package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

// ElasticIndexName is the name of the Index where we store the events received
const ElasticIndexName = "click_events"

var elasticClient = InitElastic()

// InitElastic initialise Elasticsearch and return the client
func InitElastic() *elasticsearch.Client {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Retrieve Client and Server versions
	var r map[string]interface{}
	res, err := client.Info()
	if err != nil {
		log.Fatalf("Error retrieving Elasticsearch info: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	// Create the index if not already there
	exist, err := client.Indices.Exists([]string{ElasticIndexName})
	if err != nil {
		log.Fatalf("Failed to retrieving Elasticsearch index: %s", err)
	}
	defer res.Body.Close()
	if exist.StatusCode == 404 {
		mapping, _ := json.Marshal(map[string]interface{}{
			"mappings": map[string]interface{}{
				"properties": map[string]interface{}{
					"@timestamp": map[string]interface{}{
						"type":   "date",
						"format": "epoch_millis",
					},
					"user_id": map[string]interface{}{
						"type": "integer",
					},
				},
			},
		})
		create, err := client.Indices.Create(ElasticIndexName,
			client.Indices.Create.WithBody(bytes.NewReader(mapping)),
		)
		if err != nil || create.IsError() {
			log.Fatalf("Error parsing the response body: %s - %s", err, create)
		} else {
			defer create.Body.Close()
			log.Printf("Index created: %s", create)
		}
	} else if exist.IsError() {
		log.Fatalf("Unexpected Error when checking Index: %s", exist)
	} else {
		log.Printf("Index already exist: %s", ElasticIndexName)
	}

	return client
}

func addElasticEvent(userID int) error {
	event, _ := json.Marshal(map[string]interface{}{
		"@timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		"user_id":    userID,
	})
	addIndexResp, err := elasticClient.Index(ElasticIndexName, bytes.NewReader(event))

	if err != nil {
		return fmt.Errorf("Error trying to add event: %s", err)
	}

	defer addIndexResp.Body.Close()
	if addIndexResp.IsError() {
		return fmt.Errorf("Invalid Response trying to add event: %s", addIndexResp)
	}

	return nil
}
