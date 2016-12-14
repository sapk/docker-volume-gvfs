FROM alpine:latest
MAINTAINER Antoine GIRARD <antoine.girard@sapk.fr>

ENV GOPATH="/usr/share/docker-volume-gvfs"

RUN apk -U --force --no-progress add \
          build-tools make gcc go musl-dev git ca-certificates dbus dbus-x11 gvfs gvfs-fuse gvfs-dav gvfs-smb\
    && cd /usr/share/ \
    && git clone https://github.com/sapk/docker-volume-gvfs.git \
    && cd docker-volume-gvfs && make \
    && mv ./docker-volume-gvfs /usr/bin/docker-volume-gvfs \
    && apk del --purge build-tools make gcc go musl-dev git \
    && rm -rf $GOPATH /var/cache/apk/*

ENTRYPOINT [ "/usr/bin/docker-volume-gvfs" ]
CMD [ "daemon" ]
