FROM alpine:latest
MAINTAINER "Darren Hague <d.hague@sap.com>"

ADD build/docker.tar /
ENTRYPOINT ["/usr/bin/hermes"]
