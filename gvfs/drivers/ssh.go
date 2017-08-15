package drivers

import (
	"net/url"
	"strings"
)

//SSHVolumeDriver volume driver for ssh
type SSHVolumeDriver struct {
	url *url.URL
}

func (d SSHVolumeDriver) id() DriverType {
	return SSH
}

func (d SSHVolumeDriver) isAvailable() bool {
	is, err := isFile("/usr/lib/gvfs/gvfsd-sftp")
	if err == nil {
		return is
	}
	return false
}

func (d SSHVolumeDriver) mountpoint() (string, error) {
	mount := "sftp" + ":host=" + d.url.Host
	if strings.Contains(d.url.Host, ":") {
		el := strings.Split(d.url.Host, ":")
		mount = "sftp" + ":host=" + el[0] //Default don't show port
		if el[1] != "22" {
			mount += ",port=" + el[1] //add port if not default
		}
	}
	if d.url.User != nil {
		mount += ",user=" + d.url.User.Username()
	}

	if d.url.Path != "" { //Add relative path
		mount += d.url.Path
	}

	return mount, nil
}
