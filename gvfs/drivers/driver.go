package drivers

import (
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sapk/docker-volume-helpers/basic"
	"github.com/sapk/docker-volume-helpers/driver"
	"github.com/sirupsen/logrus"
)

const (
	//MountTimeout timeout before killing a mount try in seconds
	MountTimeout = 30
	//CfgVersion current config version compat
	CfgVersion = 1
	//CfgFolder config folder
	CfgFolder = "/etc/docker-volumes/gvfs/"
)

//GVfsDriver docker volume plugin driver extension of basic plugin
type GVfsDriver = basic.Driver

//Init start all needed deps and serve response to API call
func Init(root string, dbus string, fuseOpts string) *GVfsDriver {
	logrus.Debugf("Init gluster driver at %s, UniqName: %v", root, mountUniqName)
	config := basic.DriverConfig{
		Version:       CfgVersion,
		Root:          root,
		Folder:        CfgFolder,
		MountUniqName: mountUniqName,
	}
	eventHandler := basic.DriverEventHandler{
		IsValidURI:    isValidURI,
		OnMountVolume: mountVolume,
		OnInit:        initDriver,
	}
	return basic.Init(&config, &eventHandler)
}

func mountVolume(d *basic.Driver, v driver.Volume, m driver.Mount, r *volume.MountRequest) (*volume.MountResponse, error) {
	//TODO
	return nil, nil
}

func initDriver(d *basic.Driver) error {
	//TODO
	return
}
