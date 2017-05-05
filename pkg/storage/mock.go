package storage

import "github.com/sapcc/hermes/pkg/data"

type mock struct{}

func Mock() Interface {
	return mock{}
}

func (m mock) GetEvents(filter data.Filter) ([]data.Event, error) {
	return nil, nil
}

func (m mock) GetEvent(eventId string) (data.EventDetail, error) {
	return data.EventDetail{}, nil

}
