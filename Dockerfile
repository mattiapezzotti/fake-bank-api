FROM golang:1.15-alpine AS builder

WORKDIR /srv/app/

COPY release.tar.gz /srv/app/

ENTRYPOINT []
