package main

import (
	"github.com/sapk/docker-volume-gvfs/driver"
)

var (
	Version string
	Commit  string
)

func main() {
	driver.Version = Version
	driver.Commit = Commit
	driver.Start()
}
