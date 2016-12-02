# docker-volume-gvfs [![License](https://img.shields.io/badge/license-MIT-red.svg)](https://github.com/sapk/docker-volume-gvfs/blob/master/LICENSE) ![Project Status](http://img.shields.io/badge/status-alpha-red.svg)
[![GitHub release](https://img.shields.io/github/release/sapk/docker-volume-gvfs.svg)](https://github.com/sapk/docker-volume-gvfs/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/sapk/docker-volume-gvfs)](https://goreportcard.com/report/github.com/sapk/docker-volume-gvfs) [![Travis](https://api.travis-ci.org/sapk/docker-volume-gvfs.svg?branch=master)](https://travis-ci.org/sapk/docker-volume-gvfs)

Use GVfs as a backend for docker volume

Status : **proof of concept (working)**

By using [GVfs](https://wiki.gnome.org/Projects/gvfs) this plugins allow the use of various backend as storage.

Dedends on gvfsd-fuse

## Build
```
make
```

## Start daemon
```
./docker-volume-gvfs daemon
For more advance params : ./docker-volume-gvfs --help
```

## Create and Mount volume
```
docker volume create --driver gvfs --opt url=ftp://user@url --opt password=pass --name test
docker run -v test:/mnt --rm -ti ubuntu
```

## Known bug:
If when trying to start a container you get this error :

`docker: Error response from daemon: chown /var/lib/docker-volumes/gvfs/ftp:host=10.8.0.7,user=sapk: function not implemented.`

Try to start the container with the :nocopy attribute on the volume like that :

```
docker run -v test:/mnt:nocopy --rm -ti ubuntu
```


## Inspired from :
 - https://github.com/ContainX/docker-volume-netshare/
 - https://github.com/vieux/docker-volume-sshfs/

## TODO :
 - Add test for differents backends (ftp, ssh, smb, dav, ...)
