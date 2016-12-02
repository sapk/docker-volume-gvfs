package drivers

import (
	"net/url"
	"strings"
)

//FTPVolumeDriver volume driver for ftp
type FTPVolumeDriver struct {
	url *url.URL
}

func (d FTPVolumeDriver) id() DriverType {
	return FTP
}

func (d FTPVolumeDriver) isAvailable() bool {
	is, err := isFile("/usr/lib/gvfs/gvfsd-ftp")
	if err == nil {
		return is
	}
	return false
}

func (d FTPVolumeDriver) mountpoint() (string, error) {
	mount := d.url.Scheme + ":host=" + d.url.Host
	if strings.Contains(d.url.Host, ":") {
		el := strings.Split(d.url.Host, ":")
		mount = d.url.Scheme + ":host=" + el[0] //Default don't show port
		if el[1] != "21" {
			mount += ",port=" + el[1] //add port if not default
		}
	}
	if d.url.User != nil {
		mount += ",user=" + d.url.User.Username()
	}
	return mount, nil
}
