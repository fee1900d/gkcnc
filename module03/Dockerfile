FROM golang:1.18-alpine AS build

LABEL author=yz
MAINTAINER "yz 791018061@qq.com"

RUN apk add git && \
    git clone https://github.com/fee1900d/gkcnc.git

WORKDIR gkcnc/module02 

RUN go build main/2.go

ENTRYPOINT ./2

EXPOSE 80

HEALTHCHECK --interval=5s --timeout=3s \
    CMD curl -fs http://localhost/healthz || exit 1
