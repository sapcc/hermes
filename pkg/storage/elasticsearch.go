package storage

import "github.com/sapcc/hermes/pkg/data"

type elasticSearch struct{}

func ElasticSearch() Interface {
	return elasticSearch{}
}

func (m elasticSearch) GetEvents(filter data.Filter) ([]data.Event, error) {
	return nil, nil
}

func (m elasticSearch) GetEvent(eventId string) (data.EventDetail, error) {
	return data.EventDetail{}, nil

}