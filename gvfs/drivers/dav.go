package drivers

import (
	"net/url"
	"strings"
)

//DavVolumeDriver volume driver for webdav
type DavVolumeDriver struct {
	url *url.URL
}

func (d DavVolumeDriver) id() DriverType {
	return DAV
}

func (d DavVolumeDriver) isAvailable() bool {
	is, err := isFile("/usr/lib/gvfs/gvfsd-dav")
	if err == nil {
		return is
	}
	return false
}

func (d DavVolumeDriver) mountpoint() (string, error) {
	//TODO test a lot
	mount := "dav" + ":host=" + d.url.Host
	if d.url.Scheme == "davs" { //HTTPS
		mount += ",ssl=true"
	}
	if strings.Contains(d.url.Host, ":") { //Contain custom port
		el := strings.Split(d.url.Host, ":")
		mount = "dav" + ":host=" + el[0] //Default don't show port
		if d.url.Scheme == "davs" {      //HTTPS
			mount += ",ssl=true"
			if el[1] != "443" {
				mount += ",port=" + el[1] //add port if not default
			}
		} else { //HTTP
			if el[1] != "80" {
				mount += ",port=" + el[1] //add port if not default
			}
		}
	}
	if d.url.User != nil {
		mount += ",user=" + d.url.User.Username()
	}
	if d.url.Path == "" {
		mount += ",prefix=" + d.url.EscapedPath() //TODO test
	}
	return mount, nil
}
