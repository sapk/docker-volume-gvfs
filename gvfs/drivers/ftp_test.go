package drivers

import (
	"net/url"
	"testing"
)

func TestFTPId(t *testing.T) {
	tmp := FTPVolumeDriver{}
	id := tmp.id()
	if id != FTP {
		t.Error("Expected FTP(", FTP, "), got ", id)
	}
}
func TestFTPIdName(t *testing.T) {
	tmp := FTPVolumeDriver{}
	name := tmp.id().String()
	if name != "ftp" {
		t.Error("Expected 'ftp', got ", name)
	}
}
func TestFTPSimpleIPURL(t *testing.T) {
	u, _ := url.Parse("ftp://10.0.0.1")
	tmp := FTPVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "ftp:host=10.0.0.1" {
		t.Error("Expected ftp:host=10.0.0.1, got ", m, err)
	}
}
func TestFTPHostURL(t *testing.T) {
	u, _ := url.Parse("ftp://host")
	tmp := FTPVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "ftp:host=host" {
		t.Error("Expected ftp:host=host, got ", m, err)
	}
}
func TestFTPUserInURL(t *testing.T) {
	u, _ := url.Parse("ftp://user@host")
	tmp := FTPVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "ftp:host=host,user=user" {
		t.Error("Expected ftp:host=host,user=user, got ", m, err)
	}
}
func TestFTPPortInURL(t *testing.T) {
	u, _ := url.Parse("ftp://host:42")
	tmp := FTPVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "ftp:host=host,port=42" {
		t.Error("Expected ftp:host=host,port=42, got ", m, err)
	}
}
func TestFTPDefaultPortInURL(t *testing.T) {
	u, _ := url.Parse("ftp://host:21")
	tmp := FTPVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "ftp:host=host" {
		t.Error("Expected ftp:host=host, got ", m, err)
	}
}
func TestFTPFullURL(t *testing.T) {
	u, _ := url.Parse("ftp://user@host:42")
	tmp := FTPVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "ftp:host=host,port=42,user=user" {
		t.Error("Expected ftp:host=host,port=42,user=user, got ", m, err)
	}
}
func TestFTPFullDefaultPortURL(t *testing.T) {
	u, _ := url.Parse("ftp://user@host:21")
	tmp := FTPVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "ftp:host=host,user=user" {
		t.Error("Expected ftp:host=host,user=user, got ", m, err)
	}
}
