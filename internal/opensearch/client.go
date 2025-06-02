package opensearch

import (
	"bytes"
	"catalog-service/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type ClientInterface interface {
	IndexExists(indexName string) (bool, error)
	IndexDocument(ctx context.Context, id string, document interface{}, indexName string) error
}

type Client struct {
	client *opensearch.Client
}

func NewClient(addresses []string) (*Client, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: addresses,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenSearch client: %w", err)
	}
	return &Client{client: client}, nil
}

func (c *Client) IndexExists(indexName string) (bool, error) {
	req := opensearchapi.IndicesExistsRequest{
		Index: []string{indexName},
	}
	res, err := req.Do(context.Background(), c.client)
	if err != nil {
		return false, fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	return res.StatusCode == http.StatusOK, nil
}

func (c *Client) IndexDocument(ctx context.Context, id string, document interface{}, indexName string) error {
	log := logger.NewContextLogger(ctx, "Client/IndexDocument")

	docJSON, err := json.Marshal(document)
	if err != nil {
		log.Errorf(err, "failed to marshal document: %v", err)
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	req := opensearchapi.IndexRequest{
		Index:      indexName,
		DocumentID: id,
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if res.IsError() {
		return fmt.Errorf("error indexing document: %v", response)
	}
	docID := response["_id"].(string)
	log.Infof("document indexed successfully: %s", docID)
	return nil
}
