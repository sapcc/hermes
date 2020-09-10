FROM alpine:latest
MAINTAINER "Nathan Oyler <nathan.oyler@sap.com>"
LABEL source_repository="https://github.com/sapcc/hermes"

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ADD build/docker.tar /
ENTRYPOINT ["/usr/bin/hermes"]
