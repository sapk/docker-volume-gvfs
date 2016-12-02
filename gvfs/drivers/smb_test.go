package drivers

import "testing"

func TestSMBId(t *testing.T) {
	tmp := SMBVolumeDriver{}
	id := tmp.id()
	if id != SMB {
		t.Error("Expected SMB(", SMB, "), got ", id)
	}
}
func TestSMBIdName(t *testing.T) {
	tmp := SMBVolumeDriver{}
	name := tmp.id().String()
	if name != "smb" {
		t.Error("Expected 'smb', got ", name)
	}
}
