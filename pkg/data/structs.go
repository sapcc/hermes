package data

type Filter struct {
	Source       string
	ResourceType string
	ResourceName string
	UserName     string
	EventType    string
	Time         string
	Offset       uint64
	Limit        uint64
	Sort         string
}
