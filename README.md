# Hermes

[![Build Status](https://travis-ci.org/sapcc/hermes.svg?branch=master)](https://travis-ci.org/sapcc/hermes)
[![Go Report Card](https://goreportcard.com/badge/github.com/sapcc/hermes)](https://goreportcard.com/report/github.com/sapcc/hermes)
[![GoDoc](https://godoc.org/github.com/sapcc/hermes?status.svg)](https://godoc.org/github.com/sapcc/hermes)

----

**Hermes** is an audit trail service for OpenStack. It is named after the Futurama character, not the Greek god.

It enables easy access for audit relevant events that occur from OpenStack in an Open Standards CADF Format.
* [CADF Format](http://www.dmtf.org/sites/default/files/standards/documents/DSP0262_1.0.0.pdf)
* [CADF Standards](http://www.dmtf.org/standards/cadf)

Related projects:
* [OpenStack Audit Middleware](https://github.com/sapcc/openstack-audit-middleware)
* [Keystone Event Notifications](https://docs.openstack.org/keystone/pike/advanced-topics/event_notifications.html)

----

## Features 

* Architected as a managed service for Auditing
* OpenStack Identity v3 authentication and authorization
* Project and domain-level access control (scoping)
* Compatible with other cloud based audit APIs 

## Documentation
* [Hermes API Reference](docs/hermes-api-reference.md)
* [Hermes Operators Guide](./docs/operators-guide.md)
* [Hermes Users Guide](./docs/users-guide.md)
