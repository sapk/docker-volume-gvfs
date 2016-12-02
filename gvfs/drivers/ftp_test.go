package drivers

import (
	"testing"
)

func TestSimpleIPURL(t *testing.T) {
	tmp := FTPVolumeDriver{url: "ftp://10.0.0.1"}
	if tmp.mountpoint() != "ftp:host=10.0.0.1" {
		t.Error("Expected ftp:host=10.0.0.1, got ", v)
	}
}
func TestHostURL(t *testing.T) {
	tmp := FTPVolumeDriver{url: "ftp://host"}
	if tmp.mountpoint() != "ftp:host=host" {
		t.Error("Expected ftp:host=host, got ", v)
	}
}
func TestUserInURL(t *testing.T) {
	tmp := FTPVolumeDriver{url: "ftp://user@host"}
	if tmp.mountpoint() != "ftp:host=host,user=user" {
		t.Error("Expected ftp:host=host,user=user, got ", v)
	}
}
func TestPortInURL(t *testing.T) {
	tmp := FTPVolumeDriver{url: "ftp://host:42"}
	if tmp.mountpoint() != "ftp:host=host,port=42" {
		t.Error("Expected ftp:host=host,port=42, got ", v)
	}
}
func TestDefaultPortInURL(t *testing.T) {
	tmp := FTPVolumeDriver{url: "ftp://host:21"}
	if tmp.mountpoint() != "ftp:host=host" {
		t.Error("Expected ftp:host=host, got ", v)
	}
}
func TestFullURL(t *testing.T) {
	tmp := FTPVolumeDriver{url: "ftp://user@host:42"}
	if tmp.mountpoint() != "ftp:host=host,port=42,user=sapk" {
		t.Error("Expected ftp:host=host,port=42,user=user, got ", v)
	}
}
func TestFullDefaultPortURL(t *testing.T) {
	tmp := FTPVolumeDriver{url: "ftp://user@host:21"}
	if tmp.mountpoint() != "ftp:host=host,user=sapk" {
		t.Error("Expected ftp:host=host,user=user, got ", v)
	}
}
