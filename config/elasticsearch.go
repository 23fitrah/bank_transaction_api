package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
)

type ESConfig struct {
	URL      string
	Username string
	Password string
	Timeout  time.Duration
}

var (
	esClient     *elastic.Client
	esClientLock sync.RWMutex
)

func SetEsClient(client *elastic.Client) {
	esClientLock.Lock()
	defer esClientLock.Unlock()
	esClient = client
}

func GetEsClient() *elastic.Client {
	esClientLock.RLock()
	defer esClientLock.RUnlock()
	return esClient
}

func NewElasticsearchConnection(config ESConfig) (*elastic.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	options := []elastic.ClientOptionFunc{
		elastic.SetURL(config.URL),
		elastic.SetSniff(false),
	}

	if config.Username != "" && config.Password != "" {
		options = append(options, elastic.SetBasicAuth(config.Username, config.Password))
	}

	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	info, code, err := client.Ping(config.URL).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping Elasticsearch: %w", err)
	}

	fmt.Printf("Elasticsearch Connected! Version: %s | Code: %d\n", info.Version.Number, code)

	SetEsClient(client)
	return client, nil
}
