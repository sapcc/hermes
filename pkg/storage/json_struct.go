package storage

// Thanks to the tool at https://mholt.github.io/json-to-go/

type eventListWithTotal struct {
	Total  int     `json:"total"`
	Events []Event `json:"events"`
}

type Event struct {
	PublisherID string `json:"publisher_id"`
	EventType   string `json:"event_type"`
	Payload     struct {
		Observer struct {
			TypeURI string `json:"typeURI"`
			ID      string `json:"id"`
		} `json:"observer"`
		ResourceInfo string `json:"resource_info"`
		TypeURI      string `json:"typeURI"`
		Initiator    struct {
			TypeURI   string `json:"typeURI"`
			DomainID string `json:"domain_id"`
			ProjectID string `json:"project_id"`
			UserID    string `json:"user_id"`
			Host      struct {
				Agent   string `json:"agent"`
				Address string `json:"address"`
			} `json:"host"`
			ID string `json:"id"`
		} `json:"initiator"`
		EventTime string `json:"eventTime"`
		Action    string `json:"action"`
		EventType string `json:"eventType"`
		ID        string `json:"id"`
		Outcome   string `json:"outcome"`
		Target    struct {
			TypeURI string `json:"typeURI"`
			ID      string `json:"id"`
		} `json:"target"`
	} `json:"payload"`
	MessageID string `json:"message_id"`
	Priority  string `json:"priority"`
	Timestamp string `json:"timestamp"`
}
