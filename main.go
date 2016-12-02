package main

import (
	"github.com/sapk/docker-volume-gvfs/gvfs"
)

var (
	version string
	branch  string
	commit  string
)

func main() {
	gvfs.Version = version
	gvfs.Commit = commit
	gvfs.Branch = branch
	gvfs.Start()
}
