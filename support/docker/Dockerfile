FROM golang:alpine AS build-env

COPY . /go/src/app
WORKDIR /go/src/app

RUN apk --no-cache --no-progress --force add \
    build-tools make go musl-dev git \
&&  make build

FROM alpine:latest
LABEL maintainer="Antoine GIRARD <antoine.girard@sapk.fr>"

RUN apk --no-cache --no-progress --force add \
            ca-certificates dbus dbus-x11 gvfs gvfs-fuse gvfs-dav gvfs-smb openssh-client

RUN mkdir -p /var/lib/docker-volumes/gvfs /root/.ssh

COPY --from=build-env /go/src/app/docker-volume-gvfs /usr/bin/docker-volume-gvfs

RUN /usr/bin/docker-volume-gvfs version

ENTRYPOINT [ "/usr/bin/docker-volume-gvfs" ]
CMD [ "daemon" ]
