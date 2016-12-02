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
	//Done sftp://10.8.0.7 -> sftp:host=10.8.0.7
	//Done ssh://10.8.0.7 -> sftp:host=10.8.0.7
	//Done sftp://sapk@10.8.0.7 -> sftp:host=10.8.0.7,user=sapk
	//Done sftp://sapk.fr -> sftp:host=sapk.fr
	//Done sftp://sapk@10.8.0.7:2121 -> sftp:host=10.8.0.7,port=2121,user=sapk
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
	return mount, nil
}
