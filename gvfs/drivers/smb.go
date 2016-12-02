package drivers

import (
	"fmt"
	"net/url"
)

//SMBVolumeDriver volume driver for smb
type SMBVolumeDriver struct {
	url *url.URL
}

func (d SMBVolumeDriver) id() DriverType {
	return SMB
}

func (d SMBVolumeDriver) isAvailable() bool {
	return true //TODO check for gvfsd-smb
}

func (d SMBVolumeDriver) mountpoint() (string, error) {
	return "", fmt.Errorf("Driver SMB not implemented yet")
}
