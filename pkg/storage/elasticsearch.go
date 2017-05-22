package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sapcc/hermes/pkg/util"
	"github.com/spf13/viper"
	elastic "gopkg.in/olivere/elastic.v5"
)

type elasticSearch struct {
	client *elastic.Client
}

var es elasticSearch

func ElasticSearch() Interface {
	if es.client == nil {
		es.init()
	}
	return es
}

func (es *elasticSearch) init() {
	util.LogDebug("Initiliasing ElasticSearch()")

	// Create a client
	var err error
	var url = viper.GetString("elasticsearch.url")
	util.LogDebug("Using ElasticSearch URL: %s", url)
	es.client, err = elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		panic(err)
	}
}

func (es elasticSearch) GetEvents(filter *Filter, tenant_id string) ([]*EventDetail, int, error) {
	index := indexName(tenant_id)
	util.LogDebug("Looking for events in index %s", index)

	query := elastic.NewBoolQuery()
	if filter.Source != "" {
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("publisher_id", filter.Source))
	}
	if filter.ResourceType != "" {
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("payload.target.typeURI", filter.ResourceType))
	}
	if filter.ResourceId != "" {
		query = query.Filter(elastic.NewTermQuery("payload.target.id.raw", filter.ResourceId))
	}
	if filter.UserId != "" {
		query = query.Filter(elastic.NewTermQuery("payload.initiator.user_id.raw", filter.UserId))
	}
	if filter.EventType != "" {
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("event_type", filter.EventType))
	}
	if filter.Time != "" {
		// TODO: it's complicated
	}
	if filter.Sort != "" {
		// TODO: it's complicated
	}

	search := es.client.Search().
		Index(index).
		Query(query).
		Sort("@timestamp", false).
		From(int(filter.Offset)).Size(int(filter.Limit))

	searchResult, err := search.Do(context.Background()) // execute
	if err != nil {
		return nil, 0, err
	}

	util.LogDebug("Got %d hits", searchResult.TotalHits())

	//Construct EventDetail array from search results
	var events []*EventDetail
	for _, hit := range searchResult.Hits.Hits {
		var de EventDetail
		err := json.Unmarshal(*hit.Source, &de)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, &de)
	}
	total := searchResult.TotalHits()

	return events, int(total), nil
}

func (es elasticSearch) GetEvent(eventId string, tenant_id string) (*EventDetail, error) {
	index := indexName(tenant_id)
	util.LogDebug("Looking for event %s in index %s", eventId, index)

	query := elastic.NewTermQuery("payload.id.raw", eventId)
	search := es.client.Search().
		Index(index).
		Query(query)

	searchResult, err := search.Do(context.Background())
	if err != nil {
		return nil, err
	}
	total := searchResult.TotalHits()

	if total > 0 {
		hit := searchResult.Hits.Hits[0]
		var de EventDetail
		err := json.Unmarshal(*hit.Source, &de)
		return &de, err
	} else {
		return nil, nil
	}

}

func indexName(tenant_id string) string {
	index := "audit-*"
	if tenant_id != "" {
		index = fmt.Sprintf("audit-%s-*", tenant_id)
	}
	return index
}
