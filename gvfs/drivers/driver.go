package drivers

import (
	"bytes"
	"fmt"
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
	driver      *gvfsVolumeDriver
	password    string
	mountpoint  string
	connections int
}

type gvfsVolumeDriver interface {
	id() DriverType
	isAvailable() bool
	mountpoint() (string, error)
}

//GVfsDriver the global driver responding to call
type GVfsDriver struct {
	sync.RWMutex
	root    string
	env     []string
	volumes map[string]*gvfsVolume
}

//Init start all needed deps and serve response to API call
func Init(root string, dbus string) *GVfsDriver {
	d := &GVfsDriver{
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

func (d GVfsDriver) startFuseDeamon() error {
	//TODO check needed gvfsd + gvfsd-ftp Maybe already on dbus ?
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

// start deamon in context of this gvfs drive with custome env
func (d GVfsDriver) startCmd(cmd string) error {
	log.Debugf(cmd)
	return setEnv(cmd, d.env).Start()
}

// run deamon in context of this gvfs drive with custome env
func (d GVfsDriver) runCmd(cmd string) error {
	log.Debugf(cmd)
	return setEnv(cmd, d.env).Run()
}

//Create create and init the requested volume
func (d GVfsDriver) Create(r volume.Request) volume.Response {
	log.Debugf("Entering Create: name: %s, options %v", r.Name, r.Options)
	d.Lock()
	defer d.Unlock()

	if r.Options == nil || r.Options["url"] == "" {
		return volume.Response{Err: "url option required"}
	}

	dr, m, err := getDriver(r.Options["url"])
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	v := &gvfsVolume{
		url:         r.Options["url"],
		driver:      dr,
		password:    r.Options["password"],
		mountpoint:  filepath.Join(d.root, m),
		connections: 0,
	}

	d.volumes[r.Name] = v
	log.Debugf("Volume Created: %v", v)
	return volume.Response{}
}

//Remove remove the requested volume
func (d GVfsDriver) Remove(r volume.Request) volume.Response {
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

//List volumes handled by thos driver
func (d GVfsDriver) List(r volume.Request) volume.Response {
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

//Get get info on the requested volume
func (d GVfsDriver) Get(r volume.Request) volume.Response {
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

//Path get path of the requested volume
func (d GVfsDriver) Path(r volume.Request) volume.Response {
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

//Mount mount the requested volume
func (d GVfsDriver) Mount(r volume.MountRequest) volume.Response {
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

//Unmount unmount the requested volume
//TODO Monitor for unmount to remount ?
func (d GVfsDriver) Unmount(r volume.UnmountRequest) volume.Response {
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

//Capabilities Send capabilities of the local driver
func (d GVfsDriver) Capabilities(r volume.Request) volume.Response {
	log.Debugf("Entering Capabilities: %v", r)
	return volume.Response{
		Capabilities: volume.Capability{
			Scope: "local",
		},
	}
}
