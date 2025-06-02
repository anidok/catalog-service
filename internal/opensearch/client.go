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
	Search(ctx context.Context, indexName string, searchBody map[string]interface{}) ([]map[string]interface{}, int, error)
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

func (c *Client) Search(ctx context.Context, indexName string, searchBody map[string]interface{}) ([]map[string]interface{}, int, error) {
	log := logger.NewContextLogger(ctx, "Client/Search")
	searchBodyBytes, err := json.Marshal(searchBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal search query: %w", err)
	}

	log.Debugf("search body: %s", searchBodyBytes)
	req := opensearchapi.SearchRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(searchBodyBytes),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("error executing search query: %s", res.String())
	}

	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, 0, fmt.Errorf("failed to decode search response: %w", err)
	}

	hitsArr, ok := searchResponse["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("unexpected hits format in response")
	}
	total := 0
	if totalVal, ok := searchResponse["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64); ok {
		total = int(totalVal)
	}

	hits := make([]map[string]interface{}, 0, len(hitsArr))
	for _, h := range hitsArr {
		hitMap, ok := h.(map[string]interface{})
		if !ok {
			continue
		}
		source, ok := hitMap["_source"].(map[string]interface{})
		if ok {
			hits = append(hits, source)
		}
	}

	return hits, total, nil
}
