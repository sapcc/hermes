# Hermes

[![Build Status](https://travis-ci.org/sapcc/hermes.svg?branch=master)](https://travis-ci.org/sapcc/hermes)
[![Go Report Card](https://goreportcard.com/badge/github.com/sapcc/hermes)](https://goreportcard.com/report/github.com/sapcc/hermes)
[![GoDoc](https://godoc.org/github.com/sapcc/hermes?status.svg)](https://godoc.org/github.com/sapcc/hermes)

----

**Hermes** is an audit trail service for OpenStack, originally designed for SAP's internal Openstack Cloud. 

It is named after the Futurama character, not the Greek god.

# The idea: Audit trail for OpenStack

OpenStack has an audit log, but no way for users to view these audit events. Hermes enables easy access 
to audit events on a tenant basis, relying on the ELK stack for storage. Now users can view their project level
audit events through an API, or as a module in [Elektra](https://github.com/sapcc/elektra) an OpenStack Dashboard.

It enables easy access for audit relevant events that occur from OpenStack in an Open Standards CADF Format.
* [CADF Format](http://www.dmtf.org/sites/default/files/standards/documents/DSP0262_1.0.0.pdf)
* [CADF Standards](http://www.dmtf.org/standards/cadf)

Related projects:
* [OpenStack Audit Middleware](https://github.com/sapcc/openstack-audit-middleware)
* [Keystone Event Notifications](https://docs.openstack.org/keystone/pike/advanced-topics/event_notifications.html)

----

## Features 

* A managed service for Auditing in OpenStack
* OpenStack Identity v3 authentication and authorization
* Project and domain-level access control (scoping)
* Compatible with other cloud based audit APIs 
* Exposed Prometheus metrics

## Supported Services
* Nova
* Neutron
* Designate
* Cinder
* Manila
* [Limes](https://github.com/sapcc/limes)

# Documentation

## For users

* [Hermes Users Guide](./docs/users/index.md)
* [Hermes API Reference](./docs/users/hermes-v1-reference.md)

## For operators

* [Hermes Operators Guide](./docs/operators/operators-guide.md)
