<!--
SPDX-FileCopyrightText: 2025 SAP SE

SPDX-License-Identifier: Apache-2.0
-->

# Getting Started with Hermes

Hermes is an audit trail service for OpenStack that enables easy access to audit events on a tenant basis. With Hermes, you can view project-level audit events through an API or as an optional module in the OpenStack dashboard, Elektra.

## What is an audit event?
An audit event is a JSON record that contains the details of a given OpenStack event, such as the user who made the request, the request itself, and when it occurred. The event log contains information about actions taken within your OpenStack tenant or domain.

## 7 “W”s of audit
Hermes provides detailed information about each event, including the 7 “W”s of audit: What, When, Who, FromWhere, OnWhat, Where, ToWhere. This information is presented in the CADF format, which includes both mandatory and optional properties.

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
