package main

import (
	"github.com/sapk/docker-volume-gvfs/driver"
)

var (
	version string
	branch  string
	commit  string
)

func main() {
	driver.Version = version
	driver.Commit = commit
	driver.Branch = branch
	driver.Start()
}
