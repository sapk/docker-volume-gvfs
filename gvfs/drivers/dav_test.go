package drivers

import (
	"net/url"
	"testing"
)

func TestDAVId(t *testing.T) {
	tmp := DavVolumeDriver{}
	id := tmp.id()
	if id != DAV {
		t.Error("Expected DAV(", DAV, "), got ", id)
	}
}
func TestDAVIdName(t *testing.T) {
	tmp := DavVolumeDriver{}
	name := tmp.id().String()
	if name != "dav" {
		t.Error("Expected 'dav', got ", name)
	}
}

//TODO test dav:
//TODO test IP

func TestDAVURL1(t *testing.T) {
	u, _ := url.Parse("davs://host/some/path/")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "dav:host=host,ssl=true,prefix=%2Fsome%2Fpath" {
		t.Error("Expected dav:host=host,ssl=true,prefix=%2Fsome%2Fpath, got ", m, err)
	}
}
func TestDAVURL2(t *testing.T) {
	u, _ := url.Parse("davs://host/some/path")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "dav:host=host,ssl=true,prefix=%2Fsome%2Fpath" {
		t.Error("Expected dav:host=host,ssl=true,prefix=%2Fsome%2Fpath, got ", m, err)
	}
}
func TestDAVURL3(t *testing.T) {
	u, _ := url.Parse("davs://user@host/some/path/")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "dav:host=host,ssl=true,user=user,prefix=%2Fsome%2Fpath" {
		t.Error("Expected dav:host=host,ssl=true,user=user,prefix=%2Fsome%2Fpath, got ", m, err)
	}
}
