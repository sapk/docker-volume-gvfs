package drivers

import (
	"fmt"
	"net/url"
)

//SSHVolumeDriver volume driver for ssh
type SSHVolumeDriver struct {
	url *url.URL
}

func (d SSHVolumeDriver) id() DriverType {
	return SSH
}

func (d SSHVolumeDriver) isAvailable() bool {
	return true //TODO check for gvfsd-ssh
}

func (d SSHVolumeDriver) mountpoint() (string, error) {
	return "", fmt.Errorf("Driver SSH not implemented yet")
}
