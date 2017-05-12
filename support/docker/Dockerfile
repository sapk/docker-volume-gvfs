FROM alpine:latest
MAINTAINER Antoine GIRARD <antoine.girard@sapk.fr>

ENV GOPATH="/usr/share/docker-volume-gvfs"

RUN apk --no-cache --no-progress --force add \
          build-tools make gcc go musl-dev git ca-certificates dbus dbus-x11 gvfs gvfs-fuse gvfs-dav gvfs-smb\
 && git clone https://github.com/sapk/docker-volume-gvfs.git /usr/share/docker-volume-gvfs \
 && cd /usr/share/docker-volume-gvfs && make \
 && mv ./docker-volume-gvfs /usr/bin/docker-volume-gvfs \
 && cd && rm -r /usr/share/docker-volume-gvfs \
 && apk del --purge build-tools make gcc go musl-dev git \
 && mkdir /var/lib/docker-volumes/gvfs \
 && rm -rf $GOPATH /var/cache/apk/*

ENTRYPOINT [ "/usr/bin/docker-volume-gvfs" ]
CMD [ "daemon" ]
