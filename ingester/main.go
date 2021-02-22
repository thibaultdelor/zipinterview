package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
)

// ClickEvent is the json we received from the API
type ClickEvent struct {
	UserID int `json:"user_id" binding:"required"`
}

// ElasticIndexName is the name of the Index where we store the events received
const ElasticIndexName = "click_events"

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

var elasticClient = InitElastic()

// IngestEvent take an event and put it in elastic
func IngestEvent(c *gin.Context) {
	// Validate input
	var input ClickEvent
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(fmt.Sprintf("Invalid Json, error is \"%s\"", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint("Error decoding Json")})
		return
	}

	// Add to elastic
	event, _ := json.Marshal(map[string]interface{}{
		"@timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		"user_id":    input.UserID,
	})
	addIndexResp, err := elasticClient.Index(ElasticIndexName, bytes.NewReader(event))
	if err != nil {
		log.Printf("Error trying to add event: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save event to Elastic :O"})
	} else {
		defer addIndexResp.Body.Close()
		if addIndexResp.IsError() {
			log.Printf("Invalid Response trying to add event: %s", addIndexResp)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save event to Elastic :O"})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": "event added!"})
		}
	}
}

func main() {
	r := gin.Default()
	r.PUT("/", IngestEvent)
	r.Run()
}
