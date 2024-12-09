{
  "index_patterns": ["export_events*"],  // Pattern to match the names of indices that should use this template
  "settings": {
    "number_of_shards": 1,  // Adjust based on expected load and data volume
    "number_of_replicas": 1,  // Can be adjusted for high availability
    "refresh_interval": "5s"  // Adjust based on write and query performance needs
  },
  "mappings": {
    "properties": {
      "project_id": {
        "type": "keyword"  // Using keyword type for exact value matching and aggregations
      },
      "enabled": {
        "type": "boolean"  // Boolean type for enabled/disabled status
      },
      "bucket_name": {
        "type": "keyword"  // Using keyword to facilitate exact matches
      },
      "last_run_time": {
        "type": "date",  // Date type for timestamping the last run
        "format": "epoch_millisecond"  // Storing timestamp in milliseconds since the epoch
      },
      "filters": {
        "type": "object",  // Object type to nest filter configurations
        "properties": {
          "event_types": {
            "type": "keyword"  // Array of keywords for event types
          },
          "resource_types": {
            "type": "keyword"  // Array of keywords for resource types
          },
          "exclude_event_types": {
            "type": "keyword"  // Array of keywords for excluded event types
          },
          "exclude_resource_types": {
            "type": "keyword"  // Array of keywords for excluded resource types
          }
        }
      }
    }
  }
}