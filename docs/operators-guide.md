# Hermes Operators Guide

This guide describes how to setup Openstack Auditing with Hermes.

## Introduction

Hermes depends on an existing Openstack installation that is responsible for generating
audit events from Openstack using Keystone Middleware. Each Openstack component you wish to 
use will require configuration via Keystone Middleware with the correct cadf mappings file.

This will send CADF formatted audit events to the oslo messaging bus which is held via RabbitMQ. 
From there we will take over with the Hermes infrastructure consisting of 

* Logstash for transforming events by adding metadata concerning ProjectID and DomainID
* Elasticsearch for holding the data as the main datastore
* Hermes as the Audit API for accessing these events from your Openstack Dashboard.

Optionally this can also include

* Kibana for ease of viewing data yourself. Staging in our case.

With all of this plumbing together, we build a complete Openstack Auditing setup.

## Configuration

Hermes is configured using a TOML config file that is by default located in `etc/hermes/hermes.conf`.
An example configuration file is located in etc/ which can help you get started.

#### Main Hermes config

\[hermes\]
* PolicyFilePath - Location of [OpenStack policy file](https://docs.openstack.org/security-guide/identity/policies.html) - policy.json file for which roles are required to access audit events. 
Example located in `etc/policy.json`

#### ElasticSearch configuration
Any data served by Hermes requires an underlying Elasticsearch installation to act as the Datastore.

\[elasticsearch\]
* url - Url for elasticsearch

#### Integration for Openstack Keystone
\[keystone\] 
* auth_url - Location of v3 keystone identity - ex. https://keystone.example.com/v3
* username - Openstack *service* user to authenticate and authorize clients.
* password 
* user_domain_name 
* project_name
* token_cache_time - In order to improve responsiveness and protect Keystone from too much load, Maia will
re-check authorizations for users by default every 15 minutes (900 seconds).

## Starting Hermes

Running the hermes binary will start the Server listening on `http://localhost:8788`

## Configuration of Keystone Middleware, RabbitMQ, Logstash, ElasticSearch

Documentation for [Keystone Middleware's Audit](https://docs.openstack.org/keystonemiddleware/latest/audit.html) 
describes how to enable the audit capabilities in CADF Format for
various openstack services. 

[PyCadf Audit Mappings](https://github.com/openstack/pycadf/tree/master/etc/pycadf) are used for this process.

Using the oslo.messaging bus, we configure the middleware to send audit 
events to an audit specific rabbitmq. This keeps the load on the main
oslo.messaging bus to a minimum so that auditing doesn't impact other 
core openstack services.

We then implement a Logstash instance to act as a transformation step before
loading into Elasticsearch. 

Common transforms are dropping events that don't provide value as Auditing 
events, and adding CADF mappings to events that do not currently have an 
audit map in keystone middleware due to their lack of consistent event details.
Ex: Designate Events 

From there the data is loaded into Elasticsearch where we have a rolling 
index that is created from a template to hold audit details via daily 
index.

Hermes is used as the API to query this Elasticsearch to provide API events
to the Openstack Dashboard. We use a custom Openstack Dashboard named Elektra.

## Instrumentation 

Hermes has prometheus integration located at the /metrics endpoint. Custom metrics included are

| **Name** | **Description** | 

| **hermes_request_duration_seconds** | **Duration of a Hermes request** | 
| --- | --- | 
| hermes_requests_inflight |  Number of inflight HTTP requests served by Hermes |
| hermes_response_size_bytes | Size of the Hermes response (e.g. to retrieve events) | 
| hermes_storage_errors_count | Number of technical errors occurred when accessing underlying storage | 
