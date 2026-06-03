package utils

import (
	"context"
	"fmt"
	"time"
	"transaction_api/config"

	"github.com/olivere/elastic/v7"
)

var logQueue = make(chan logEntry, 100)

type logEntry struct {
	esClient    *elastic.Client
	indexPrefix string
	data        map[string]interface{}
}

func init() {
	go processLogQueue()
}

func processLogQueue() {
	for entry := range logQueue {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := entry.esClient.Index().Index(entry.indexPrefix).BodyJson(entry.data).Do(ctx)
		if err != nil {
			fmt.Printf("Error logging to ES: %v\n", err)
		}
		cancel()
	}
}

func LogToES(indexPrefix, operation, message string, metadata map[string]interface{}) {
	esClient := config.GetEsClient()
	if esClient == nil {
		return
	}
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["@timestamp"] = time.Now()
	metadata["operation"] = operation
	metadata["message"] = message

	logQueue <- logEntry{
		esClient:    esClient,
		indexPrefix: "transaction_api-"+indexPrefix,
		data:        metadata,
	}
}
