package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var Client *elasticsearch.Client

func Init() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	var err error
	Client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := Client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	log.Println("Connected to Elasticsearch")
}

func IndexResource(resourceType, id string, content interface{}) {
	body, err := json.Marshal(content)
	if err != nil {
		log.Printf("Error marshaling document: %s", err)
		return
	}

	req := esapi.IndexRequest{
		Index:      strings.ToLower(resourceType),
		DocumentID: id,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), Client)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), id)
	} else {
		log.Printf("[%s] Indexed document ID=%s", res.Status(), id)
	}
}

func SearchResources(resourceType string, query map[string]string) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	
	// Basic query construction (match all if empty, or match specific fields)
	queryMap := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	if len(query) > 0 {
		matches := []map[string]interface{}{}
		for k, v := range query {
			matches = append(matches, map[string]interface{}{
				"match": map[string]interface{}{
					k: v,
				},
			})
		}
		queryMap = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": matches,
				},
			},
		}
	}

	if err := json.NewEncoder(&buf).Encode(queryMap); err != nil {
		return nil, err
	}

	res, err := Client.Search(
		Client.Search.WithContext(context.Background()),
		Client.Search.WithIndex(strings.ToLower(resourceType)),
		Client.Search.WithBody(&buf),
		Client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.Status())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		results = append(results, source)
	}

	return results, nil
}
