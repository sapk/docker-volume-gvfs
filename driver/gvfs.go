package driver

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
)

const (
	//MountTimeout timeout before killing a mount try in seconds
	MountTimeout = 30
)

type gvfsVolume struct {
	url         string
	password    string
	mountpoint  string
	connections int
}

type gvfsDriver struct {
	sync.RWMutex
	root    string
	env     []string
	volumes map[string]*gvfsVolume
}

func newGVfsDriver(root string, dbus string) *gvfsDriver {
	d := &gvfsDriver{
		root:    root,
		env:     make([]string, 1),
		volumes: make(map[string]*gvfsVolume),
	}
	if dbus == "" {
		// start needed dbus like (eval `dbus-launch --sh-syntax`) and get env variable
		result, err := exec.Command("dbus-launch", "--sh-syntax").CombinedOutput() //DBUS_SESSION_BUS_ADDRESS='unix:abstract=/tmp/dbus-JHGXLpeJ6A,guid=25ab632502ebccd43cd403bc58388fab';\n ...
		if err != nil {
			panic(err)
		}
		env := string(result)
		log.Debugf("dbus-launch --sh-syntax -> \n%s", env)
		reDBus := regexp.MustCompile("DBUS_SESSION_BUS_ADDRESS='(.*?)';")
		//rePID := regexp.MustCompile("DBUS_SESSION_BUS_PID=(.*?);")
		matchDBuse := reDBus.FindStringSubmatch(env)
		//matchPID := rePID.FindStringSubmatch(env)
		dbus = matchDBuse[1]
		//TODO plan to kill this add closing ?
	}
	d.env[0] = fmt.Sprintf("DBUS_SESSION_BUS_ADDRESS=%s", dbus)
	err := d.startFuseDeamon()
	if err != nil {
		panic(err) //Something went wrong
	}
	return d
}

func setEnv(cmd string, env []string) *exec.Cmd {
	c := exec.Command("sh", "-c", cmd)
	//c := exec.Command(strings.Split(cmd, " ")) //TODO better
	c.Env = env
	return c
}

// start deamon in context of this gvfs drive with custome env
func (d gvfsDriver) startCmd(cmd string) error {
	log.Debugf(cmd)
	return setEnv(cmd, d.env).Start()
}

// run deamon in context of this gvfs drive with custome env
func (d gvfsDriver) runCmd(cmd string) error {
	log.Debugf(cmd)
	return setEnv(cmd, d.env).Run()
}

func (d gvfsDriver) startFuseDeamon() error {
	//TODO check needed gvfsd + gvfsd-ftp Maybe allready on dbus ?
	// Normaly gvfsd-fuse block such so this like crash but global ?

	fi, err := os.Lstat(d.root)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(d.root, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	if fi != nil && !fi.IsDir() {
		return fmt.Errorf("%v already exist and it's not a directory", d.root)
	}

	err = d.startCmd("/usr/lib/gvfs/gvfsd --no-fuse") //Start global deamon
	if err != nil {
		return err
	}

	err = d.startCmd(fmt.Sprintf("/usr/lib/gvfs/gvfsd-fuse %s -f -o big_writes,use_ino,allow_other,auto_cache,umask=0022", d.root)) //Start ftp handler
	if err != nil {
		return err
	}
	return nil
}

func urlToMountPoint(root string, urlString string) (string, error) {
	//Done ftp://sapk@10.8.0.7 -> ftp:host=10.8.0.7,user=sapk
	//Done ftp://10.8.0.7 -> ftp:host=10.8.0.7
	//Done ftp://sapk.fr -> ftp:host=sapk.fr
	//Done ftp://sapk@10.8.0.7:2121 -> ftp:host=10.8.0.7,port=2121,user=sapk
	//Done ftp://sapk@10.8.0.7:21 -> ftp:host=10.8.0.7,user=sapk
	//TODO other sheme
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	name := u.Scheme + ":host=" + u.Host
	if strings.Contains(u.Host, ":") {
		el := strings.Split(u.Host, ":")
		name = u.Scheme + ":host=" + el[0] //Default don't show port
		if u.Scheme == "ftp" && el[1] != "21" {
			name = u.Scheme + ":host=" + el[0] + ",port=" + el[1] //add port if not default
		}
	}
	if u.User != nil {
		name += ",user=" + u.User.Username()
	}
	return filepath.Join(root, name), nil
}

func (d gvfsDriver) Create(r volume.Request) volume.Response {
	log.Debugf("Entering Create: name: %s, options %v", r.Name, r.Options)
	d.Lock()
	defer d.Unlock()

	if r.Options == nil || r.Options["url"] == "" {
		return volume.Response{Err: "url option required"}
	}

	m, err := urlToMountPoint(d.root, r.Options["url"])
	if err != nil {
		return volume.Response{Err: err.Error()}
	}
	v := &gvfsVolume{
		url:         r.Options["url"],
		password:    r.Options["password"],
		mountpoint:  m,
		connections: 0,
	}

	d.volumes[r.Name] = v
	log.Debugf("Volume Created: %v", v)
	return volume.Response{}
}

func (d gvfsDriver) Remove(r volume.Request) volume.Response {
	log.Debugf("Entering Remove: name: %s, options %v", r.Name, r.Options)
	d.Lock()
	defer d.Unlock()
	v, ok := d.volumes[r.Name]

	if !ok {
		return volume.Response{Err: fmt.Sprintf("volume %s not found", r.Name)}
	}
	if v.connections == 0 {
		delete(d.volumes, r.Name)
		return volume.Response{}
	}
	return volume.Response{Err: fmt.Sprintf("volume %s is currently used by a container", r.Name)}
}

func (d gvfsDriver) List(r volume.Request) volume.Response {
	log.Debugf("Entering List: name: %s, options %v", r.Name, r.Options)

	d.Lock()
	defer d.Unlock()

	var vols []*volume.Volume
	for name, v := range d.volumes {
		vols = append(vols, &volume.Volume{Name: name, Mountpoint: v.mountpoint})
		log.Debugf("Volume found: %s", v)
	}
	return volume.Response{Volumes: vols}
}

func (d gvfsDriver) Get(r volume.Request) volume.Response {
	log.Debugf("Entering Get: name: %s", r.Name)
	d.Lock()
	defer d.Unlock()

	v, ok := d.volumes[r.Name]
	if !ok {
		return volume.Response{Err: fmt.Sprintf("volume %s not found", r.Name)}
	}

	log.Debugf("Volume found: %s", v)
	return volume.Response{Volume: &volume.Volume{Name: r.Name, Mountpoint: v.mountpoint}}
}

func (d gvfsDriver) Path(r volume.Request) volume.Response {
	log.Debugf("Entering Path: name: %s, options %v", r.Name)

	d.RLock()
	defer d.RUnlock()
	v, ok := d.volumes[r.Name]
	if !ok {
		return volume.Response{Err: fmt.Sprintf("volume %s not found", r.Name)}
	}
	log.Debugf("Volume found: %s", v)
	return volume.Response{Mountpoint: v.mountpoint}
}

func (d gvfsDriver) Mount(r volume.MountRequest) volume.Response {
	log.Debugf("Entering Mount: %v", r)
	d.Lock()
	defer d.Unlock()

	v, ok := d.volumes[r.Name]
	if !ok {
		return volume.Response{Err: fmt.Sprintf("volume %s not found", r.Name)}
	}

	if v.connections > 0 {
		v.connections++
		return volume.Response{Mountpoint: v.mountpoint}
	}

	cmd := fmt.Sprintf("gvfs-mount %s", v.url)
	if v.password != "" {
		p := setEnv(cmd, d.env)
		inStd, err := p.StdinPipe()
		if err != nil { //Get a input buffer
			return volume.Response{Err: err.Error()}
		}
		var outStd bytes.Buffer
		p.Stdout = &outStd
		var errStd bytes.Buffer
		p.Stderr = &errStd

		if err := p.Start(); err != nil {
			return volume.Response{Err: err.Error()}
		}
		inStd.Write([]byte(v.password + "\n")) //Send password to process + Send return line

		// wait or timeout
		donec := make(chan error, 1)
		go func() {
			donec <- p.Wait() //Process finish
		}()
		select {
		case <-time.After(MountTimeout * time.Second):
			sOut := outStd.String()
			sErr := errStd.String()
			p.Process.Kill()
			log.Debugf("out : %s", sOut)
			log.Debugf("outErr : %s", sErr)
			return volume.Response{Err: fmt.Sprintf("The command %s timeout", cmd)}
		case <-donec:
			sOut := outStd.String()
			sErr := errStd.String()
			log.Debugf("Password send and command %s return", cmd)
			log.Debugf("out : %s", sOut)
			log.Debugf("outErr : %s", sErr)
			// handle erros like : "Error mounting location: Location is already mounted" or Error mounting location: Could not connect to 10.8.0.7: No route to host
			if strings.Contains(sErr, "Error mounting location") {
				return volume.Response{Err: fmt.Sprintf("Error mounting location : %s", sErr)}
			}
			v.connections++
			break
		}
	} else {
		if err := d.runCmd(cmd); err != nil {
			return volume.Response{Err: err.Error()}
		}
	}

	return volume.Response{Mountpoint: v.mountpoint}
}

//TODO Monitor for unmount to remount ?
func (d gvfsDriver) Unmount(r volume.UnmountRequest) volume.Response {
	//Execute gvfs-mount -u $params
	log.Debugf("Entering Unmount: %v", r)

	d.Lock()
	defer d.Unlock()
	v, ok := d.volumes[r.Name]
	if !ok {
		return volume.Response{Err: fmt.Sprintf("volume %s not found", r.Name)}
	}
	if v.connections <= 1 {
		cmd := fmt.Sprintf("gvfs-mount -u %s", v.url)
		if err := d.runCmd(cmd); err != nil {
			return volume.Response{Err: err.Error()}
		}
		v.connections = 0
	} else {
		v.connections--
	}

	return volume.Response{}
}

func (d gvfsDriver) Capabilities(r volume.Request) volume.Response {
	log.Debugf("Entering Capabilities: %v", r)
	return volume.Response{
		Capabilities: volume.Capability{
			Scope: "local",
		},
	}
}
