package data

// Event contains high-level data about an event, intended as a list item
type Event struct {
	ID   string `json:"event_id"`
	Type string `json:"event_type"`
	Time string `json:"event_time"`
}

// EventDetail contains the CADF payload, enhanced with names for IDs
// TODO - add lots of fields
type EventDetail struct {
	ID   string `json:"event_id"`
	Type string `json:"event_type"`
	Time string `json:"eventTime"`
}

type Filter struct {
	source       string
	resourceType string
	resourceName string
	userName     string
	eventType    string
	time         string
	offset       uint32
	limit        uint8
	sort         string
}
