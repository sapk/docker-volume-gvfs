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

func TestDAVURLSSL(t *testing.T) {
	expected := "dav:host=host,ssl=true,prefix=%2Fsome%2Fpath"
	u, _ := url.Parse("davs://host/some/path")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != expected {
		t.Error("Expected ", expected, ", got ", m, err)
	}
}

func TestDAVURLSSLTraillingSlash(t *testing.T) {
	expected := "dav:host=host,ssl=true,prefix=%2Fsome%2Fpath"
	u, _ := url.Parse("davs://host/some/path/")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != expected {
		t.Error("Expected ", expected, ", got ", m, err)
	}
}

func TestDAVURLSSLUser(t *testing.T) {
	expected := "dav:host=host,ssl=true,user=user,prefix=%2Fsome%2Fpath"
	u, _ := url.Parse("davs://user@host/some/path/")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != expected {
		t.Error("Expected ", expected, ", got ", m, err)
	}
}

func TestDAVURLSSLDefaultPort(t *testing.T) {
	expected := "dav:host=host,ssl=true,prefix=%2Fsome%2Fpath"
	u, _ := url.Parse("davs://host:443/some/path")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != expected {
		t.Error("Expected ", expected, ", got ", m, err)
	}
}

func TestDAVURLSSLCustomPort(t *testing.T) {
	expected := "dav:host=host,port=4443,ssl=true,prefix=%2Fsome%2Fpath"
	u, _ := url.Parse("davs://host:4443/some/path")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != expected {
		t.Error("Expected ", expected, ", got ", m, err)
	}
}

func TestDAVURLDefault(t *testing.T) {
	expected := "dav:host=host,ssl=false,prefix=%2Fsome%2Fpath"
	u, _ := url.Parse("dav://host/some/path")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != expected {
		t.Error("Expected ", expected, ", got ", m, err)
	}
}

func TestDAVURLUser(t *testing.T) {
	expected := "dav:host=host,ssl=false,user=user,prefix=%2Fsome%2Fpath"
	u, _ := url.Parse("dav://user@host/some/path")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != expected {
		t.Error("Expected ", expected, ", got ", m, err)
	}
}

func TestDAVURLDefaultPort(t *testing.T) {
	expected := "dav:host=host,ssl=false,prefix=%2Fsome%2Fpath"
	u, _ := url.Parse("dav://host:80/some/path")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != expected {
		t.Error("Expected ", expected, ", got ", m, err)
	}
}

func TestDAVURLIP(t *testing.T) {
	expected := "dav:host=10.0.0.1,ssl=false,prefix=%2Fsome%2Fpath"
	u, _ := url.Parse("dav://10.0.0.1/some/path")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != expected {
		t.Error("Expected ", expected, ", got ", m, err)
	}
}
func TestDAVURLCustomPort(t *testing.T) {
	expected := "dav:host=host,port=8080,ssl=false,prefix=%2Fsome%2Fpath"
	u, _ := url.Parse("dav://host:8080/some/path")
	tmp := DavVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != expected {
		t.Error("Expected ", expected, ", got ", m, err)
	}
}
