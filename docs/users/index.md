# Documentation for Hermes users

Hermes is an audit trail service for OpenStack. If auditing is enabled and Hermes is deployed to OpenStack, per tenant 
audit events are available as an API to users, as well as an optional Dashboard component.

A Hermes event is a JSON record containing the details of a given OpenStack event. The event log contains
information about actions taken within your OpenStack tenant or domain, such as who made the request, 
what the request was, and when the request occurred.

For any given event, the log contains the following details: What, When, Who, From Where, On What, Where, To Where of an activity<sup>*</sup>. This is also referred to as the 7 W’s of audit and compliance.

&ast;*an activity is a type of event that provides information about actions having occurred or intended to occur, and initiated by some resource or done against some resource.*

### The 7 “W”s of audit

| “W” Component | CADF Mandatory Properties  | CADF Optional Properties (where applicable) | Description |
| --- | --- | --- | --- |
| What | `event.action`<br>`event.outcome`<br>`event.eventType` | `event.reason` | “what” activity occurred; “what” was the result. |
| When | `event.eventTime` || “when” did it happen. |
| Who | `event.initiator.id`<br>`event.initiator.typeURI` | `event.initiator.name` | “who” (person or service) initiated the action. |
| FromWhere || `event.initiator.host`<br>`event.initiator.domain`<br>`event.initiator.domain_id`<br>`event.initiator.project_id` | "FromWhere" provides information describing where the action was initiated from. |
| OnWhat | `event.target.id`<br>`event.target.typeURI`  | `event.target.domain_id`<br>`event.target.project_id` | “onWhat” resource did the activity target. |
| Where | `event.observer.id`<br>`event.observer.typeURI` | `event.observer.name` | “where” did the activity get observed (reported), or modified in some way. |
| ToWhere ||| "ToWhere" provides information describing where the target resource that is affected by the action is located. |


## Available clients

* Hermes command line client [HermesCli](https://github.com/sapcc/hermescli)
* You can send requests to [the HTTP API](./hermes-v1-reference.md) directly, as shown [in this guide](./api-example.md).
* The OpenStack web dashboard [Elektra](https://github.com/sapcc/elektra) contains an optional *Audit*
  module that becomes accessible if Hermes is deployed in the target OpenStack cluster.

## Retention of audit events

Retention is configurable on a global level for all tenants. In the roadmap it is intended that retention will be
on a per tenant basis. The current basis for retention is 3 months.
