package basic

import (
	"io"
	"net/url"
	"os"

	"github.com/docker/go-plugins-helpers/volume"
)

//GetMountName get moint point base on request and driver config (mountUniqName)
func GetMountName(d *Driver, r *volume.CreateRequest) string {
	if d.Config.MountUniqName {
		return url.PathEscape(r.Options["voluri"])
	}
	return url.PathEscape(r.Name)
}

//FolderIsEmpty based on: http://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty
func FolderIsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
