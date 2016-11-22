package driver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/volume"
)

type gvfsVolume struct {
	uri         string
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
	cmd := fmt.Sprintf("/usr/lib/gvfs/gvfsd-fuse %s -f -o big_writes", root)
	log.Debugf(cmd)
	return exec.Command("sh", "-c", cmd).Start() //We execut in background
}

func (d gvfsDriver) Create(r volume.Request) volume.Response {
	log.Debugf("Entering Create: name: %s, options %v", r.Name, r.Options)
	d.Lock()
	defer d.Unlock()

	if r.Options == nil || r.Options["uri"] == "" {
		return volume.Response{Err: "uri option required"}
	}

	v := &gvfsVolume{
		uri:         r.Options["uri"],
		mountpoint:  filepath.Join(d.root, r.Options["uri"]), //TODO test //ftp://sapk@10.8.0.7 -> ftp:host=10.8.0.7,user=sapk
		connections: 0,
	}

	d.volumes[r.Name] = v
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
	log.Debugf("Entering Mount: %v", r)
	d.Lock()
	defer d.Unlock()

	v, ok := d.volumes[r.Name]
	if !ok {
		return responseError(fmt.Sprintf("volume %s not found", r.Name))
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
	/*
		if err := d.mountVolume(v); err != nil {
			return volume.Response{Err: err.Error()}
		}
	*/
	/*
		cmd := fmt.Sprintf("sshfs -oStrictHostKeyChecking=no %s %s", v.sshcmd, v.mountpoint)
		if v.password != "" {
			cmd = fmt.Sprintf("echo %s | %s -o workaround=rename -o password_stdin", v.password, cmd)
		}
		logrus.Debug(cmd)
		return exec.Command("sh", "-c", cmd).Run() //ftp://sapk@10.8.0.7 -> ftp:host=10.8.0.7,user=sapk
	*/
	//TODO gvfs-mount ftp://sapk@10.8.0.7
	//TODO return volume.Response{Mountpoint: v.mountpoint}
}

func (d gvfsDriver) Unmount(r volume.UnmountRequest) volume.Response {
	//TODO gvfs-mount -u ftp://sapk@10.8.0.7

}

func (d gvfsDriver) Capabilities(r volume.Request) volume.Response {
	log.Debugf("Entering Capabilities: %v", r)
	return volume.Response{
		Capabilities: volume.Capability{
			Scope: "local",
		},
	}
}
