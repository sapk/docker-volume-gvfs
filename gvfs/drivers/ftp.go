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
	return true //TODO check for gvfsd-ftp
}

func (d FTPVolumeDriver) mountpoint() (string, error) {
	//Done ftp://sapk@10.8.0.7 -> ftp:host=10.8.0.7,user=sapk
	//Done ftp://10.8.0.7 -> ftp:host=10.8.0.7
	//Done ftp://sapk.fr -> ftp:host=sapk.fr
	//Done ftp://sapk@10.8.0.7:2121 -> ftp:host=10.8.0.7,port=2121,user=sapk
	//Done ftp://sapk@10.8.0.7:21 -> ftp:host=10.8.0.7,user=sapk
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