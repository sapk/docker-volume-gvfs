package drivers

import (
	"net/url"
	"testing"
)

func TestSSHId(t *testing.T) {
	tmp := SSHVolumeDriver{}
	id := tmp.id()
	if id != SSH {
		t.Error("Expected SSH(", SSH, "), got ", id)
	}
}
func TestSSHIdName(t *testing.T) {
	tmp := SSHVolumeDriver{}
	name := tmp.id().String()
	if name != "ssh" {
		t.Error("Expected 'ssh', got ", name)
	}
}

func TestSSHSimpleIPURL(t *testing.T) {
	u, _ := url.Parse("sftp://10.0.0.1")
	tmp := SSHVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "sftp:host=10.0.0.1" {
		t.Error("Expected sftp:host=10.0.0.1, got ", m, err)
	}
}
func TestSSHSimpleHostURL(t *testing.T) {
	u, _ := url.Parse("sftp://hostname/")
	tmp := SSHVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "sftp:host=hostname/" {
		t.Error("Expected sftp:host=hostname/, got ", m, err)
	}
}
func TestSSHSimpleHostSSHURL(t *testing.T) {
	u, _ := url.Parse("ssh://hostname/")
	tmp := SSHVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "sftp:host=hostname/" {
		t.Error("Expected sftp:host=hostname/, got ", m, err)
	}
}
func TestSSHCustomPortURL(t *testing.T) {
	u, _ := url.Parse("sftp://hostname:42/")
	tmp := SSHVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "sftp:host=hostname,port=42/" {
		t.Error("Expected sftp:host=hostname,port=42/, got ", m, err)
	}
}
func TestSSHDefaultPortURL(t *testing.T) {
	u, _ := url.Parse("sftp://hostname:22/")
	tmp := SSHVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "sftp:host=hostname/" {
		t.Error("Expected sftp:host=hostname/, got ", m, err)
	}
}
func TestSSHSetUserURL(t *testing.T) {
	u, _ := url.Parse("sftp://sapk@hostname:42/")
	tmp := SSHVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "sftp:host=hostname,port=42,user=sapk/" {
		t.Error("Expected sftp:host=hostname,port=42,user=sapk/, got ", m, err)
	}
}
func TestSSHSimpleSomePathURL(t *testing.T) {
	u, _ := url.Parse("sftp://hostname/some/path")
	tmp := SSHVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "sftp:host=hostname/some/path" {
		t.Error("Expected sftp:host=hostname/some/path, got ", m, err)
	}
}
func TestSSHSimpleSomePathURL2(t *testing.T) {
	u, _ := url.Parse("ssh://hostname/some/path")
	tmp := SSHVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "sftp:host=hostname/some/path" {
		t.Error("Expected sftp:host=hostname/some/path, got ", m, err)
	}
}
func TestSSHComplexSomePathURL(t *testing.T) {
	u, _ := url.Parse("sftp://sapk@hostname:42/some/path")
	tmp := SSHVolumeDriver{url: u}
	m, err := tmp.mountpoint()
	if err != nil || m != "sftp:host=hostname,port=42,user=sapk/some/path" {
		t.Error("Expected sftp:host=hostname,port=42,user=sapk/some/path, got ", m, err)
	}
}
