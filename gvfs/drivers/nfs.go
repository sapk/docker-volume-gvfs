package drivers

import (
	"net/url"
	"strings"
)

//NFSVolumeDriver volume driver for ssh
type NFSVolumeDriver struct {
	url *url.URL
}

func (d NFSVolumeDriver) id() DriverType {
	return NFS
}

func (d NFSVolumeDriver) isAvailable() bool {
	is, err := isFile("/usr/lib/gvfs/gvfsd-nfs")
	if err == nil {
		return is
	}
	return false
}

//TODO test
func (d NFSVolumeDriver) mountpoint() (string, error) {
	mount := "nfs" + ":host=" + d.url.Host
	if strings.Contains(d.url.Host, ":") {
		el := strings.Split(d.url.Host, ":")
		mount = "nfs" + ":host=" + el[0] //Default don't show port
		if el[1] != "111" {
			mount += ",port=" + el[1] //add port if not default
		}
	}
	if d.url.User != nil {
		mount += ",user=" + d.url.User.Username()
	}
	return mount, nil
}
