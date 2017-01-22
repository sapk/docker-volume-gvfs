package main

import (
	"github.com/sapk/docker-volume-gvfs/gvfs"
)

var (
	//Version version of app set by build flag
	Version string
	//Branch git branch of app set by build flag
	Branch string
	//Commit git commit of app set by build flag
	Commit string
	//BuildTime build time of app set by build flag
	BuildTime string
)

func main() {
	gvfs.Version = Version
	gvfs.Commit = Commit
	gvfs.Branch = Branch
	gvfs.BuildTime = BuildTime
	gvfs.Start()
}
