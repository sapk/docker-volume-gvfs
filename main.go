package main

import (
	"github.com/sapk/docker-volume-gvfs/driver"
)

var VERSION string = ""
var BUILD_DATE string = ""

func main() {
	driver.Version = VERSION
	driver.BuildDate = BUILD_DATE
	driver.Start()
}
