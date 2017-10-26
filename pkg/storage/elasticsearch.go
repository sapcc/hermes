package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sapcc/hermes/pkg/util"
	"github.com/spf13/viper"
	"gopkg.in/olivere/elastic.v5"
	"strings"
)

type ElasticSearch struct {
	esClient *elastic.Client
}

func (es *ElasticSearch) client() *elastic.Client {
	// Lazy initialisation - don't connect to ElasticSearch until we need to
	if es.esClient == nil {
		es.init()
	}
	return es.esClient
}

func (es *ElasticSearch) init() {
	util.LogDebug("Initializing ElasticSearch()")

	// Create a client
	var err error
	var url = viper.GetString("elasticsearch.url")
	util.LogDebug("Using ElasticSearch URL: %s", url)
	// Added disabling sniffing for Testing from Golang. This corrects a problem. Likely needs to be removed before prod deploy
	es.esClient, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	//es.esClient, err = elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		panic(err)
	}
}

// GetEvents grabs events for a given tenantID with filtering.
func (es ElasticSearch) GetEvents(filter *Filter, tenantId string) ([]*EventDetail, int, error) {
	index := indexName(tenantId)
	util.LogDebug("Looking for events in index %s", index)

	query := elastic.NewBoolQuery()
	if filter.Source != "" {
		util.LogDebug("Filtering on Source %s", filter.Source)
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("observer.typeURI", filter.Source))
	}
	if filter.ResourceType != "" {
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("target.typeURI", filter.ResourceType))
	}
	if filter.ResourceId != "" {
		query = query.Filter(elastic.NewTermQuery("target.id", filter.ResourceId))
	}
	if filter.UserId != "" {
		query = query.Filter(elastic.NewPrefixQuery("initiator.id", filter.UserId))
	}
	if filter.EventType != "" {
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("action", filter.EventType))
	}
	if filter.Time != nil && len(filter.Time) > 0 {
		for key, value := range filter.Time {
			timeField := "eventTime"
			switch key {
			case "lt":
				query = query.Filter(elastic.NewRangeQuery(timeField).Lt(value))
			case "lte":
				query = query.Filter(elastic.NewRangeQuery(timeField).Lte(value))
			case "gt":
				query = query.Filter(elastic.NewRangeQuery(timeField).Gt(value))
			case "gte":
				query = query.Filter(elastic.NewRangeQuery(timeField).Gte(value))
			}
		}
	}

	// Fields are mapped differently in ES, must maintain a mapping to json.
	esFieldMapping := map[string]string{
		"time":          "eventTime",
		"source":        "observer.typeURI",
		"resource_type": "target.typeURI",
		"resource_name": "target.id",
		"event_type":    "action",
	}

	esSearch := es.client().Search().
		Index(index).
		Query(query)

	if filter.Sort != nil {
		for _, fieldOrder := range filter.Sort {
			switch fieldOrder.Order {
			case "asc":
				esSearch = esSearch.Sort(esFieldMapping[fieldOrder.Fieldname], true)
			case "desc":
				esSearch = esSearch.Sort(esFieldMapping[fieldOrder.Fieldname], false)
			}
		}
	}

	esSearch = esSearch.
		Sort("eventTime", false).
		From(int(filter.Offset)).Size(int(filter.Limit))

	searchResult, err := esSearch.Do(context.Background()) // execute
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

// GetEvent Returns EventDetail for a single event.
func (es ElasticSearch) GetEvent(eventId string, tenantId string) (*EventDetail, error) {
	index := indexName(tenantId)
	util.LogDebug("Looking for event %s in index %s", eventId, index)

	query := elastic.NewTermQuery("id", eventId)
	esSearch := es.client().Search().
		Index(index).
		Query(query)

	searchResult, err := esSearch.Do(context.Background())
	if err != nil {
		return nil, err
	}
	total := searchResult.TotalHits()

	if total > 0 {
		hit := searchResult.Hits.Hits[0]
		var de EventDetail
		err := json.Unmarshal(*hit.Source, &de)
		return &de, err
	}
	return nil, nil
}

// GetAttributes Return all unique attributes available for filtering
// Possible queries, event_type, dns, identity, etc..
func (es ElasticSearch) GetAttributes(queryName string, tenantId string) ([]string, error) {
	index := indexName(tenantId)

	util.LogDebug("Looking for unique attributes for %s in index %s", queryName, index)

	// Mapping for attributes based on return values to API
	// Source in this case is not the cadf source, but instead the first part of event_type
	esFieldMapping := map[string]string{
		"time":          "eventTime",
		"source":        "observer.typeURI",
		"resource_type": "target.typeURI",
		"resource_name": "target.id",
		"event_type":    "action",
	}

	var esName string
	util.LogDebug("Mapped Queryname: %s", esFieldMapping[queryName])
	// Append .raw onto queryName, in Elasticsearch. Aggregations turned on for .raw
	if val, ok := esFieldMapping[queryName]; ok {
		esName = val + ".raw"
	} else {
		esName = queryName + ".raw"
	}

	queryAgg := elastic.NewTermsAggregation().Field(esName)

	esSearch := es.client().Search().Index(index).Aggregation("attributes", queryAgg)
	searchResult, err := esSearch.Do(context.Background())
	if err != nil {
		return nil, err
	}

	if searchResult.Hits == nil {
		util.LogDebug("expected Hits != nil; got: nil")
	}

	agg := searchResult.Aggregations
	if agg == nil {
		util.LogDebug("expected Aggregations, got nil")
	}

	termsAggRes, found := agg.Terms("attributes")
	if !found {
		util.LogDebug("Term %s not found in Aggregation", esName)
	}
	if termsAggRes == nil {
		util.LogDebug("termsAggRes is nil")
	}
	util.LogDebug("Number of Buckets: %d", len(termsAggRes.Buckets))

	var unique []string
	for _, bucket := range termsAggRes.Buckets {
		util.LogDebug("key: %s count: %d", bucket.Key, bucket.DocCount)
		if queryName == "source" {
			source := strings.SplitN(bucket.Key.(string), ".", 2)[0]
			unique = append(unique, source)
		} else {
			unique = append(unique, bucket.Key.(string))
		}
	}

	unique = SliceUniqMap(unique)
	return unique, nil
}

// SliceUniqMap Removes duplicates from slice
func SliceUniqMap(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

// MaxLimit grabs the configured maxlimit for results
func (es ElasticSearch) MaxLimit() uint {
	return uint(viper.GetInt("elasticsearch.max_result_window"))
}

// indexName Generates the index name for a given TenantId. If no tenantId defaults to audit-*
// Records for audit-* will not be accessible from a given Tenant.
func indexName(tenantId string) string {
	index := "audit-*"
	if tenantId != "" {
		index = fmt.Sprintf("audit-%s-*", tenantId)
	}
	return index
}
