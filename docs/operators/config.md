# Configuration Guide

Hermes is configured using a TOML config file that is by default located in `etc/hermes/hermes.conf`.
An example configuration file is located in etc/ which can help you get started.

#### Main Hermes config

\[hermes\]
* PolicyFilePath - Location of [OpenStack policy file](https://docs.OpenStack.org/security-guide/identity/policies.html) - policy.json file for which roles are required to access audit events. 
Example located in `etc/policy.json`

#### ElasticSearch configuration
Any data served by Hermes requires an underlying ElasticSearch installation to act as the Datastore.

\[ElasticSearch\]
* url - Url for ElasticSearch


#### Environment Variables

To configure secure access to Elasticsearch, set the following environment variables:

- `HERMES_ES_USERNAME`: The username for connecting to Elasticsearch.
- `HERMES_ES_PASSWORD`: The password for connecting to Elasticsearch.

These environment variables can be set in the deployment environment, or you may include them in your Kubernetes configuration if you are deploying Hermes there.

#### Example usage:

```bash
export HERMES_ES_USERNAME="your_username_here"
export HERMES_ES_PASSWORD="your_password_here"
```


#### Integration for OpenStack Keystone
\[keystone\] 
* auth_url - Location of v3 keystone identity - ex. https://keystone.example.com/v3
* username - OpenStack *service* user to authenticate and authorize clients.
* password 
* user_domain_name 
* project_name
* token_cache_time - In order to improve responsiveness and protect Keystone from too much load, Hermes will
re-check authorizations for users by default every 15 minutes (900 seconds).

