package drivers

import "testing"

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
