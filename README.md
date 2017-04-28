# Hermes

[![Build Status](https://travis-ci.org/sapcc/hermes.svg?branch=master)](https://travis-ci.org/sapcc/hermes)
[![Coverage Status](https://coveralls.io/repos/github/sapcc/hermes/badge.svg?branch=master)](https://coveralls.io/github/sapcc/hermes?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sapcc/hermes)](https://goreportcard.com/report/github.com/sapcc/hermes)
[![GoDoc](https://godoc.org/github.com/sapcc/hermes?status.svg)](https://godoc.org/github.com/sapcc/hermes)

Hermes is an OpenStack-compatible audit data retrieval service, originally designed for SAP's internal cloud.

It is named after the Futurama character, not the Greek god.

# Background

*TO DO*

# Installation

There's a Makefile, so do:

* `make` to just compile and run the binaries from the `build/` directory
* `make && make install` to install to `/usr`
* `make && make install PREFIX=/some/path` to install to `/some/path`
* `make docker` to build the Docker image (set image name and tag with the `DOCKER_IMAGE` and `DOCKER_TAG` variables)

## Usage
 *TO DO

1. Write a configuration file for your environment, by following the [example configuration][ex-conf].

[ex-conf]:  ./etc/hermes.conf
