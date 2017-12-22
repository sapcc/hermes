# Auditing Log Reference

A Hermes event log is a JSON record containing the details of a given event. The event log contains
information about actions taken within your Converged Cloud account, such as who made the request, 
what the request was, and when the request occurred.

## Event Detail

**Parameters**

| **Name** | **Description** |
| --- | --- |
| source | OpenStack Service Origin (ex: identity) |
| event_id | Unique ID for event within Auditing system |
| event_type | Description of event (ex: identity.project.created) |
| event_time | Timestamp event occured (ex: 2017-04-20T11:27:15.834562+00:00) |
| resource_id | Unique ID for targeted resource  |
| resource_type | Targeted service description (ex: data/security/project)  |
| initiator.domain_id | Unique ID for the domain  |
| initiator.typeURI | Indicates that the initiator is a user  |
| initiator.user_id | UUID for User that initiated the action  |
| initiator.host.agent | Agent where the OpenStack compute service request came from |
| initiator.host.address | Address information where the OpenStack compute service request came from |
| initiator.id | OpenStack initiator Unique Id |



```json
{
  "events": [
    {
      "source": "identity",
      "event_id": "3824e534-6cd4-53b2-93d4-33dc4ab50b8c",
      "event_type": "identity.project.created",
      "event_time": "2017-04-20T11:27:15.834562+00:00",
      "resource_id": "3a7e3d2421384f56a8fb6cf082a8efab",
      "resource_type": "data/security/project",
      "initiator": {
        "domain_id": "39a253e16e4a4a3686edca72c8e101bc",
        "typeURI": "service/security/account/user",
        "user_id": "275e9a16294b3805c8dd2ab77123531af6aacd92182ddcd491933e5c09864a1d",
        "host": {
           "agent": "python-keystoneclient",
           "address": "100.66.0.24"
        },
        "id": "493b9a5284675cbb9f3f6439bd222eb6"
      }
    },
  ]
}
```

