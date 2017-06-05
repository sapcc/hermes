FROM alpine:latest
MAINTAINER "Darren Hague <d.hague@sap.com>"

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ADD build/docker.tar /
ENTRYPOINT ["/usr/bin/hermes"]
