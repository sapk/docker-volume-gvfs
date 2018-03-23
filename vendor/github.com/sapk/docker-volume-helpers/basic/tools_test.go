package basic

import (
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
)

func TestMountName(t *testing.T) {
	name := GetMountName(&Driver{
		Config: DriverConfig{
			MountUniqName: false,
		},
	}, &volume.CreateRequest{
		Name: "test",
		Options: map[string]string{
			"voluri": "gluster-node:volname",
		},
	})

	if name != "test" {
		t.Error("Expected to be test, got ", name)
	}

	nameuniq := GetMountName(&Driver{
		Config: DriverConfig{
			MountUniqName: true,
		},
	}, &volume.CreateRequest{
		Name: "test",
		Options: map[string]string{
			"voluri": "gluster-node:volname",
		},
	})

	if nameuniq != "gluster-node:volname" {
		t.Error("Expected to be gluster-node:volname, got ", name)
	}
}
