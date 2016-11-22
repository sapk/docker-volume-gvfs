package driver

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
)

type gvfsVolume struct {
	url         string
	mountpoint  string
	connections int
}

type gvfsDriver struct {
	sync.RWMutex
	root    string
	volumes map[string]*gvfsVolume
}

func newGVfsDriver(root string) *gvfsDriver {
	d := &gvfsDriver{
		root:    root,
		volumes: make(map[string]*gvfsVolume),
	}
	d.startFuseDeamon()
	return d
}

func (d gvfsDriver) startFuseDeamon() error {
	//TODO check if not allready started by other ? -> Normaly gvfsd-fuse block such so this like crash
	//TODO check if folder need to be created like in Mount
	cmd := fmt.Sprintf("/usr/lib/gvfs/gvfsd-fuse %s -f -o big_writes", d.root)
	log.Debugf(cmd)
	return exec.Command("sh", "-c", cmd).Start() //We execut in background
}

func urlToMountPoint(root string, urlString string) (string, error) {
	//Done ftp://sapk@10.8.0.7 -> ftp:host=10.8.0.7,user=sapk
	//TODO ftp://sapk@10.8.0.7:42 -> ftp:host=10.8.0.7,user=sapk,port=4242 ???
	//TODO ftp://10.8.0.7 -> ftp:host=10.8.0.7 ???
	//TODO ftp://sapk.fr -> ftp:host=sapk.fr ???
	//TODO other sheme
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	name := u.Scheme + ":host=" + u.Host
	if u.User != nil { //TODO test it
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
		/** //Maybe a little to much to remove all ?
			if err := os.RemoveAll(v.mountpoint); err != nil {
				return volume.Response{Err: err.Error()}
			}
		/**/
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
	return volume.Response{Mountpoint: v.mountpoint}
}

func (d gvfsDriver) Mount(r volume.MountRequest) volume.Response {
	//TODO manage allready mountpoint allready exist before ? maybe init to +1 ?
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

	fi, err := os.Lstat(v.mountpoint)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(v.mountpoint, 0755); err != nil {
			return volume.Response{Err: err.Error()}
		}
	} else if err != nil {
		return volume.Response{Err: err.Error()}
	}

	if fi != nil && !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v already exist and it's not a directory", v.mountpoint)}
	}

	//Example gvfs-mount ftp://sapk@10.8.0.7
	cmd := fmt.Sprintf("gvfs-mount %s", v.url)
	if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{Mountpoint: v.mountpoint}
}

func (d gvfsDriver) Unmount(r volume.UnmountRequest) volume.Response {
	//Example gvfs-mount -u ftp://sapk@10.8.0.7
	d.Lock()
	defer d.Unlock()
	v, ok := d.volumes[r.Name]
	if !ok {
		return volume.Response{Err: fmt.Sprintf("volume %s not found", r.Name)}
	}
	if v.connections <= 1 {
		cmd := fmt.Sprintf("gvfs-mount -u %s", v.url)
		if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
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
