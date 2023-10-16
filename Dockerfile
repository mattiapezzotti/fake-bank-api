FROM golang:1.15-alpine AS builder

COPY release.tar.gz /srv/app/

ENTRYPOINT []
