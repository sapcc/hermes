<!--
SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company

SPDX-License-Identifier: Apache-2.0
-->

# ExportEvents Feature Design Document

## Overview
The ExportEvents feature is an enhancement to the Hermes audit trail service for OpenStack. It enables users to configure the export of audit events at the project level, allowing them to transfer filtered audit events to a designated S3 (Swift) bucket for long-term storage and analysis.

## Architecture
The ExportEvents feature will be implemented as an extension to the existing Hermes system. It will involve the following components:

- Elasticsearch: An "export_events" index will be created to store the export configuration for each project.
- Export Worker: A worker process will be implemented to periodically retrieve and export audit events for enabled projects.
- API Endpoints: New API endpoints will be introduced to manage the export configuration for projects.
- CLI: The Hermes CLI will be enhanced to support commands related to the ExportEvents feature.
- UI: The Elektra OpenStack Dashboard will be extended to provide a user interface for configuring the ExportEvents feature.

## Elasticsearch Index
- Index Name: "export_events"
- Fields:
  - `project_id` (string): The ID of the project.
  - `enabled` (boolean): Indicates whether the export is enabled for the project.
  - `bucket_name` (string): The name of the S3 (Swift) bucket where the exported audit events will be stored.
  - `last_run_time` (timestamp): The timestamp of the last successful export run for the project.
  - `filters` (object): An object containing the filter configuration for the project.
    - `event_types` (array of strings): An array of event types to include in the export (e.g., ["identity.user.created", "compute.instance.created"]).
    - `resource_types` (array of strings): An array of resource types to include in the export (e.g., ["user", "instance"]).
    - `exclude_event_types` (array of strings): An array of event types to exclude from the export (e.g., ["identity.role.deleted"]).
    - `exclude_resource_types` (array of strings): An array of resource types to exclude from the export (e.g., ["volume"]).

## Export Worker
- The export worker will be implemented in `storage/exportevents/worker.go`.
- It will run periodically (e.g., every 15 minutes) to process and export audit events for enabled projects.
- For each enabled project, the worker will:
  - Retrieve the last run time from the project's export configuration.
  - Fetch the audit events for the project since the last run time using the modified `GetEvents` function.
  - Filter the events based on the project's export configuration (e.g., event type, resource type).
  - Aggregate and compress the filtered events.
  - Upload the compressed events to the designated S3 (Swift) bucket.
  - Update the last run time in the project's export configuration.

## API Endpoints
- `POST /v1/project/{project_id}/export-events`: Enables or updates the export configuration for a project, including filter options.
- `GET /v1/project/{project_id}/export-events`: Retrieves the export configuration for a project, including filter options.

## CLI Commands
- `hermes export-events enable --project-id <project_id> --bucket-name <bucket_name> --retention-period <retention_period> --event-types <event_types> --resource-types <resource_types> --exclude-event-types <exclude_event_types> --exclude-resource-types <exclude_resource_types>`: Enables the ExportEvents feature for a project with filter options.
- `hermes export-events update --project-id <project_id> --event-types <event_types> --resource-types <resource_types> --exclude-event-types <exclude_event_types> --exclude-resource-types <exclude_resource_types>`: Updates the filter configuration for a project.

## UI Integration
- The Hermes module in the Elektra OpenStack Dashboard will be extended to include a section for the ExportEvents feature.
- It will provide a user-friendly interface for enabling/disabling the feature, specifying the S3 (Swift) bucket, and configuring event filters.
- The UI will integrate with the API endpoints to manage the export configuration for projects.

## Testing
- Unit tests will be written for all components related to the ExportEvents feature, including Elasticsearch functions, API endpoints, export worker, and CLI commands.
- Integration tests will be conducted to verify the end-to-end functionality of the feature, testing various scenarios and edge cases.
- Performance and scalability tests will be performed to ensure the feature can handle a large number of projects and high volume of audit events.

## Documentation
- The Hermes user guide and API reference will be updated to include detailed information about the ExportEvents feature.
- Step-by-step instructions and examples will be provided to guide users in configuring and using the feature through the API, CLI, and UI.

## Roadmap
1. Design and implement the Elasticsearch index and functions.
2. Modify the event retrieval logic to support filtering by project ID and time range.
3. Implement the export worker to process and export audit events periodically.
4. Develop the API endpoints for managing the export configuration.
5. Enhance the Hermes CLI to support the ExportEvents feature.
6. Integrate the feature with the Elektra OpenStack Dashboard.
7. Update the documentation to include information about the ExportEvents feature.
8. Conduct thorough testing, including unit tests, integration tests, and performance tests.
9. Address any identified issues and optimize the feature based on feedback and metrics.
10. Deploy the ExportEvents feature to production and monitor its usage and effectiveness.

## Conclusion
The ExportEvents feature enhances the Hermes audit trail service by providing users with the capability to export filtered audit events at the project level. By enabling the transfer of audit events to a designated S3 (Swift) bucket, this feature facilitates long-term storage, analysis, and compliance requirements.

The implementation of the ExportEvents feature involves extending various components of the Hermes system, including Elasticsearch, the export worker, API endpoints, CLI, and UI. Thorough testing, documentation, and consideration for performance, scalability, and security are essential to ensure a reliable and user-friendly feature.

By following the outlined roadmap and addressing the identified considerations, the ExportEvents feature can be successfully implemented and integrated into the Hermes audit trail service, providing enhanced functionality and value to OpenStack users.
