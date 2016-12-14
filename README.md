# docker-volume-gvfs [![License](https://img.shields.io/badge/license-MIT-red.svg)](https://github.com/sapk/docker-volume-gvfs/blob/master/LICENSE) ![Project Status](http://img.shields.io/badge/status-alpha-red.svg)
[![GitHub release](https://img.shields.io/github/release/sapk/docker-volume-gvfs.svg)](https://github.com/sapk/docker-volume-gvfs/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/sapk/docker-volume-gvfs)](https://goreportcard.com/report/github.com/sapk/docker-volume-gvfs)
[![codecov](https://codecov.io/gh/sapk/docker-volume-gvfs/branch/master/graph/badge.svg)](https://codecov.io/gh/sapk/docker-volume-gvfs)
 master : [![Travis master](https://api.travis-ci.org/sapk/docker-volume-gvfs.svg?branch=master)](https://travis-ci.org/sapk/docker-volume-gvfs) develop : [![Travis develop](https://api.travis-ci.org/sapk/docker-volume-gvfs.svg?branch=develop)](https://travis-ci.org/sapk/docker-volume-gvfs)

Use GVfs as a backend for docker volume

Status : **proof of concept (working)**


By using [GVfs](https://wiki.gnome.org/Projects/gvfs) this plugins allow the use of various backend as storage.

Dedends on gvfs and gvfsd-fuse (so dbus indirectly)

Packages : [![Archlinux](https://img.shields.io/badge/Archlinux-AUR-blue.svg)](https://aur.archlinux.org/packages/docker-volume-gvfs-git/)

Working: SSH/SFTP/DAV/DAVS/FTP  
In Progress: FTPS/NFS/SMB/CIFS/...
## Build
```
make
```

## Start daemon
```
./docker-volume-gvfs daemon
OR in a docker container
docker run -d --device=/dev/fuse:/dev/fuse --cap-add=SYS_ADMIN --cap-add=MKNOD  -v /run/docker/plugins:/run/docker/plugins -v /var/lib/docker-volumes/gvfs:/var/lib/docker-volumes/gvfs:shared sapk/docker-volume-gvfs
```

For more advance params : ```./docker-volume-gvfs --help OR ./docker-volume-gvfs daemon --help```
```
Run listening volume drive deamon to listen for mount request

Usage:
  docker-volume-gvfs daemon [flags]

Flags:
  -d, --dbus string        DBus address to use for gvfs.  Can also set default environment DBUS_SESSION_BUS_ADDRESS
  -o, --fuse-opts string   Fuse options to use for gvfs moint point (default "big_writes,allow_other,auto_cache")

Global Flags:
  -b, --basedir string   Mounted volume base directory (default "/var/lib/docker-volumes/gvfs")
  -v, --verbose          Turns on verbose logging
```

## Create and Mount volume
```
docker volume create --driver gvfs --opt url=ftp://user@url --opt password=pass --name test
docker run -v test:/mnt --rm -ti ubuntu
```
NB : For mounting SSH/SFTP it is usefull to set a [ssh_config](https://linux.die.net/man/5/ssh_config) file for the running user in order to use a ssh key as authentification.

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
