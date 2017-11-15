# Public API specification

## GET /v1/events

Lists a project’s or domain's audit events. The project or domain comes from the 
scope of the authentication token.

The list of events can be filtered by the parameters passed in via the URL.

Only basic event data will be listed here (event id, event name, resource name,
resource type, user name). Clients must make a separate call to retrieve the full 
CADF payload data for each individual event.

The website for CADF is [here](http://www.dmtf.org/standards/cadf).
More details on the CADF format, with examples for each OpenStack service, can
be found in the PDF 
[here](http://www.dmtf.org/sites/default/files/standards/documents/DSP2038_1.1.0.pdf).

**Parameters**

| **Name** | **Type** | **Description** |
| --- | --- | --- |
| source | string | Selects all events with source similar to this value. |
| resource\_type | string | Selects all events with resource type similar to this value. |
| resource\_name | string | Selects all events with resource name similar to this value. |
| user\_name | integer | Selects all events with user name equal to this value. |
| event\_type | string | Selects all events with event\_type equal to this value. |
| time | string | Date filter to select all events with _event_time_ matching the specified criteria. See Date Filters below for more detail. |
| offset | integer | The starting index within the total list of the events that you would like to retrieve. |
| limit | integer | The maximum number of records to return (up to 100). The default limit is 10. |
| sort | string | Determines the sorted order of the returned list. See Sorting below for more detail. |

**Date Filters:**

The value for the `time` parameter is a comma-separated list of time stamps in ISO 
8601 format. The time stamps can be prefixed with any of these comparison operators:
`gt:` (greater-than), `gte:` (greater-than-or-equal), `lt:` (less-than), `lte:` 
(less-than-or-equal).

For example, to get a list of events that will expire in January of 2020:
```
GET /v1/events?time=gte:2020-01-01T00:00:00,lt:2020-02-01T00:00:00
```

**Sorting:**

The value of the sort parameter is a comma-separated list of sort keys. Supported 
sort keys include `time`, `source`, `resource_type`, `resource_name`, and `event_type`.

Each sort key may also include a direction. Supported directions are `:asc` for 
ascending and `:desc` for descending. The service will use `:asc` for every key 
that does not include a direction.

For example, to sort the list from most recently created to oldest:

```
GET /v1/events?sort=time:desc
```

**Request:**

```
GET /v1/events?offset=1&limit=2&sort=time

Headers:
    Accept: application/json
    X-Auth-Token: {keystone_token}
```

**Response:**

This example shows the audit events for creating & deleting a project.

```json
{
  "next": "http://{hermes_host}:8788/v1/events?limit=2&offset=3",
  "previous": "http://{hermes_host}:8788/v1/events?limit=2&offset=0",
  "events": [
    {
      "event_id": "3824e534-6cd4-53b2-93d4-33dc4ab50b8c",
      "event_name": "identity.project.created",
      "event_time": "2017-04-20T11:27:15.834562+0000",
      "resource_name": "temp_project",
      "resource_id": "3a7e3d2421384f56a8fb6cf082a8efab",
      "resource_type": "data/security/project",
      "initiator": {
        "domain_id": "39a253e16e4a4a3686edca72c8e101bc",
        "domain_name": "monsoon3",
        "typeURI": "service/security/account/user",
        "user_id": "275e9a16294b3805c8dd2ab77123531af6aacd92182ddcd491933e5c09864a1d",
        "user_name": "I056593",
        "host": {
           "agent": "python-keystoneclient",
           "address": "100.66.0.24"
        }
      }
    },
    {
      "event_id": "1ff4703a-d8c3-50f8-94d1-8ab382941e80",
      "event_name": "identity.project.deleted",
      "event_time": "2017-04-20T11:28:32.521298+0000",
      "resource_name": "temp_project",
      "resource_id": "3a7e3d2421384f56a8fb6cf082a8efab",
      "resource_type": "data/security/project",
      "initiator": {
        "domain_id": "39a253e16e4a4a3686edca72c8e101bc",
        "domain_name": "monsoon3",
        "typeURI": "service/security/account/user",
        "user_id": "275e9a16294b3805c8dd2ab77123531af6aacd92182ddcd491933e5c09864a1d",
        "user_name": "I056593",
        "host": {
           "agent": "python-keystoneclient",
           "address": "100.66.0.24"
        }
      }    }
  ],
  "total": 5
}
```

**Response Attributes**

| **Name** | **Type** | **Description** |
| --- | --- | --- |
| events | list | Contains a list of events. The attributes in the event objects are the same as for an individual event. |
| total | integer | The total number of events available to the user. |
| next | string | A HATEOAS URL to retrieve the next set of events based on the offset and limit parameters. This attribute is only available when the total number of events is greater than offset and limit parameter combined. |
| previous | string | A HATEOAS URL to retrieve the previous set of events based on the offset and limit parameters. This attribute is only available when the request offset is greater than 0. |

**HTTP Status Codes**

| **Code** | **Description** |
| --- | --- |
| 200 | Successful Request |
| 401 | Invalid X-Auth-Token or the token doesn&#39;t have permissions to this resource |

## GET /v1/events/:event_id

Returns the full CADF payload for an individual
event, e.g.:

```json
{
   "publisher_id": "identity.keystone-2031324599-cgpyi",
   "event_type": "identity.project.deleted",
   "payload": {
      "observer": {
         "typeURI": "service/security",
         "id": "3824e534-6cd4-53b2-93d4-33dc4ab50b8c"
      },
      "resource_info": "d2eec974d849446da1715923e60d0b3b",
      "typeURI": "http://schemas.dmtf.org/cloud/audit/1.0/event",
      "initiator": {
         "domain_id": "39a253e16e4a4a3686edca72c8e101bc",
         "domain_name": "monsoon3",
         "typeURI": "service/security/account/user",
         "user_id": "275e9a16294b3805c8dd2ab77123531af6aacd92182ddcd491933e5c09864a1d",
         "user_name": "I056593",
         "host": {
            "agent": "python-keystoneclient",
            "address": "100.66.0.24"
         },
         "id": "493b9a5284675cbb9f3f6439bd222eb6"
      },
      "eventTime": "2017-04-20T11:28:32.521298+0000",
      "action": "deleted.project",
      "eventType": "activity",
      "id": "1ff4703a-d8c3-50f8-94d1-8ab382941e80",
      "outcome": "success",
      "target": {
         "typeURI": "data/security/project",
         "id": "d2eec974d849446da1715923e60d0b3b",
         "name": "temp_project"
      }
   },
   "message_id": "d4f88c45-5fea-4013-80ec-2d357eab37c3",
   "priority": "info",
   "timestamp": "2017-04-20 11:28:32.521769"
}
```