package storage

import "github.com/sapcc/hermes/pkg/data"
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	elastic "gopkg.in/olivere/elastic.v5"
	"log"
	"strings"
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
	fmt.Println("Initiliasing ElasticSearch()")

	// Create a client
	var err error
	var url = viper.GetString("elasticsearch.url")
	log.Printf("Using ElasticSearch URL: %s", url)
	es.client, err = elastic.NewClient()
	if err != nil {
		panic(err)
	}
}

func (es elasticSearch) GetEvents(filter data.Filter, tenant_id string) ([]*data.Event, int, error) {
	index := fmt.Sprintf("audit-%s-*", tenant_id)

	// Search with a term query
	query := elastic.NewMatchAllQuery()
	search := es.client.Search().
		Index(index).
		Query(query).
		Sort("@timestamp", true).
		From(int(filter.Offset)).Size(int(filter.Limit)).
		Pretty(true) // pretty print request and response JSON

	ctx := context.Background()
	searchResult, err := search.Do(ctx) // execute

	if err != nil {
		// Handle error
		panic(err)
	}
	//Construct data.Event array from search results
	var events []*data.Event
	for _, hit := range searchResult.Hits.Hits {
		var de data.EventDetail
		err := json.Unmarshal(*hit.Source, &de)
		p := de.Payload
		ev := data.Event{
			Source:       strings.SplitN(de.EventType, ".", 2)[0],
			ID:           p.ID,
			Type:         de.EventType,
			Time:         p.EventTime,
			ResourceId:   de.Payload.Target.ID,
			ResourceType: de.Payload.Target.TypeURI,
		}
		err = copier.Copy(&ev.Initiator, &de.Payload.Initiator)
		if err != nil {
			return nil, 0, nil
		}
		events = append(events, &ev)
	}
	total := searchResult.TotalHits()

	return events, int(total), nil
}

func (es elasticSearch) GetEvent(eventId string, tenant_id string) (data.EventDetail, error) {
	return data.EventDetail{}, nil

}
