package drivers

import (
	"net/url"
	"strings"
)

//SMBVolumeDriver volume driver for smb (not tested yet)
type SMBVolumeDriver struct {
	url *url.URL
}

func (d SMBVolumeDriver) id() DriverType {
	return SMB
}

func (d SMBVolumeDriver) isAvailable() bool {
	is, err := isFile("/usr/lib/gvfsd-smb")
	if err == nil {
		return is
	}
	return false
}

func (d SMBVolumeDriver) mountpoint() (string, error) {
	mount := d.url.Scheme + ":host=" + d.url.Host
	if strings.Contains(d.url.Host, ":") {
		el := strings.Split(d.url.Host, ":")
		mount = d.url.Scheme + ":host=" + el[0] //Default don't show port
		if el[1] != "445" {
			mount += ",port=" + el[1] //add port if not default
		}
	}
	if d.url.User != nil {
		mount += ",user=" + d.url.User.Username()
	}
	return mount, nil
}
