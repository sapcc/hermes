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

#####Main Hermes config

\[hermes\]
* PolicyFilePath - Location of [OpenStack policy file](https://docs.openstack.org/security-guide/identity/policies.html) - policy.json file for which roles are required to access audit events. 
Example located in `etc/policy.json`
* enrich_keystone_events - Defaults to false, will optionally change UUIDs to real names.

#####ElasticSearch configuration
Any data served by Hermes requires an underlying Elasticsearch installation to act as the Datastore.

\[elasticsearch\]
* url - Url for elasticsearch

#####Integration for Openstack Keystone
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

TODO - Discuss configuration of each part to have a working full path Auditing System.

