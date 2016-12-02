package drivers

import "testing"

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
