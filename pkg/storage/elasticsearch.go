package storage

import "github.com/sapcc/hermes/pkg/data"

type elasticSearch struct{}

func ElasticSearch() Interface {
	return elasticSearch{}
}

func (m elasticSearch) GetEvents(filter data.Filter) ([]*data.Event, int, error) {
	return nil, 0, nil
}

func (m elasticSearch) GetEvent(eventId string) (data.EventDetail, error) {
	return data.EventDetail{}, nil

}