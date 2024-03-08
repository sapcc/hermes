# Configuration Guide

Hermes is configured using environment variables. You can set the following environment variables to configure Hermes:

#### Main Hermes config

- `HERMES_DEBUG`: Enable debug logging. Set to `true` to enable debug mode.
- `HERMES_KEYSTONE_DRIVER`: Specifies the keystone driver to use. Default is `keystone`.
- `HERMES_STORAGE_DRIVER`: Specifies the storage driver to use. Default is `elasticsearch`.
- `HERMES_POLICY_FILE_PATH`: Location of [OpenStack policy file](https://docs.OpenStack.org/security-guide/identity/policies.html) - policy.json file for which roles are required to access audit events. Default is `etc/policy.json`.

#### ElasticSearch configuration

Any data served by Hermes requires an underlying ElasticSearch installation to act as the Datastore.

- `HERMES_ES_URL`: URL for ElasticSearch. Default is `http://localhost:9200`.
- `HERMES_ES_USERNAME`: The username for connecting to Elasticsearch.
- `HERMES_ES_PASSWORD`: The password for connecting to Elasticsearch.
- `HERMES_ES_MAX_RESULT_WINDOW`: The maximum result window for Elasticsearch queries. Default is `20000`.

#### Integration for OpenStack Keystone

- `HERMES_OS_AUTH_URL`: Location of v3 keystone identity - ex. `https://keystone.example.com/v3`.
- `HERMES_OS_USERNAME`: OpenStack *service* user to authenticate and authorize clients.
- `HERMES_OS_PASSWORD`: Password for the OpenStack service user.
- `HERMES_OS_USER_DOMAIN_NAME`: User domain name for the OpenStack service user.
- `HERMES_OS_PROJECT_NAME`: Project name for the OpenStack service user.
- `HERMES_OS_PROJECT_DOMAIN_NAME`: Project domain name for the OpenStack service user.
- `HERMES_OS_TOKEN_CACHE_TIME`: In order to improve responsiveness and protect Keystone from too much load, Hermes will re-check authorizations for users by default every 15 minutes (900 seconds). You can adjust this value as needed.
- `HERMES_OS_MEMCACHED_SERVERS`: Comma-separated list of memcached servers to use for token caching (optional).

#### API Configuration

- `HERMES_API_LISTEN_ADDRESS`: The address and port for the Hermes API server to listen on. Default is `0.0.0.0:8788`.

#### Example usage:

```bash
export HERMES_DEBUG=true
export HERMES_KEYSTONE_DRIVER=keystone
export HERMES_STORAGE_DRIVER=elasticsearch
export HERMES_POLICY_FILE_PATH=/path/to/policy.json
export HERMES_ES_URL=http://elasticsearch:9200
export HERMES_ES_USERNAME=your_username_here
export HERMES_ES_PASSWORD=your_password_here
export HERMES_ES_MAX_RESULT_WINDOW=20000
export HERMES_OS_AUTH_URL=https://keystone.example.com/v3
export HERMES_OS_USERNAME=hermes_service_user
export HERMES_OS_PASSWORD=your_password_here
export HERMES_OS_USER_DOMAIN_NAME=Default
export HERMES_OS_PROJECT_NAME=service
export HERMES_OS_PROJECT_DOMAIN_NAME=Default
export HERMES_OS_TOKEN_CACHE_TIME=900
export HERMES_API_LISTEN_ADDRESS=0.0.0.0:8788