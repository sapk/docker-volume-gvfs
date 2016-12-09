package drivers

import "testing"

func TestNFSId(t *testing.T) {
	tmp := NFSVolumeDriver{}
	id := tmp.id()
	if id != NFS {
		t.Error("Expected NFS(", NFS, "), got ", id)
	}
}
func TestNFSIdName(t *testing.T) {
	tmp := NFSVolumeDriver{}
	name := tmp.id().String()
	if name != "nfs" {
		t.Error("Expected 'nfs', got ", name)
	}
}
