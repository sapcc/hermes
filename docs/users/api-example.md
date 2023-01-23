# Using the Hermes API

The Hermes API allows you to access audit events on a tenant basis, providing detailed information about each event, including the 7 “W”s of audit: What, When, Who, FromWhere, OnWhat, Where, ToWhere. This guide will walk you through the process of getting a token, finding the Hermes endpoint, and using the API.

If you would prefer to use a command line to access the API please use [HermesCLI](https://github.com/sapcc/hermescli)

## Getting a Token

Before you can use the API, you need to get an authentication token. You can do this by using the OpenStack client. 
Here are the steps to get a token:

1. Install the OpenStack client: Follow the instructions in the OpenStack documentation for [how to install and use it](https://docs.openstack.org/user-guide/cli.html).
2. Provide your credentials to the OpenStack client: 
3. Get a token: 

```bash
export OS_AUTH_TOKEN="$(openstack token issue -f value -c id)"
```

This command will not print any output if it is successful.

## Finding Hermes

Query the service catalog to find the Hermes endpoint. It can be identified by looking for the `resources` service type:

```bash
$ openstack catalog list
+---------------+---------------+--------------------------------------------------------------------------+
| Name          | Type          | Endpoints                                                                |
+---------------+---------------+--------------------------------------------------------------------------+
| keystone      | identity      | staging                                                                  |
|               |               |   public: https://identity.example.com:443/v3                            |
|               |               | staging                                                                  |
|               |               |   internal: http://keystone.openstack.svc.kubernetes.example.com:5000/v3 |
|               |               | staging                                                                  |
|               |               |   admin: https://identity-admin.example.com:443/v3                       |
|               |               |                                                                          |
| ...           | ...           | ...                                                                      |
|               |               |                                                                          |
| hermes        | audit-data    | staging                                                                  |
|               |               |   public: https://hermes.example.com                                     |
|               |               |                                                                          |
| ...           | ...           | ...                                                                      |
|               |               |                                                                          |
+---------------+---------------+--------------------------------------------------------------------------+
```

### Using Hermes

In this case, the endpoint URL for Hermes is `https://hermes.example.com`, so you can build a request URL by appending 
one of the paths from the [API reference][v1-api]. For example, to show quota and usage data for a project, use the
following command:

```bash
curl -H "X-Auth-Token: $OS_AUTH_TOKEN" https://hermes.example.com/v1/events
```

`$OS_AUTH_TOKEN` is the token from the first step. `$DOMAIN_ID` and `$PROJECT_ID` need to be set by you to the project
ID in question and its domain ID. If you only have a project name, you can get these IDs by calling `openstack project
show $NAME`.

[os-cli]: https://docs.openstack.org/user-guide/cli.html
[v1-api]: ./hermes-v1-reference.md
