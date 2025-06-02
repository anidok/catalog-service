package opensearch

import (
	"bytes"
	"catalog-service/internal/config"
	"catalog-service/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type Client interface {
	IndexExists(indexName string) (bool, error)
	IndexDocument(ctx context.Context, id string, document interface{}, indexName string) error
	Search(ctx context.Context, indexName string, searchBody map[string]interface{}) ([]map[string]interface{}, int, error)
	FindDocumentByID(ctx context.Context, indexName, id string) (map[string]interface{}, error)
}

type ClientImpl struct {
	*opensearch.Client
}

func NewClient(addresses []string) (*ClientImpl, error) {
	osCfg := config.OpenSearch()
	transport := &http.Transport{
		MaxIdleConns:        osCfg.MaxIdleConns(),
		MaxIdleConnsPerHost: osCfg.MaxIdleConnsPerHost(),
		IdleConnTimeout:     osCfg.IdleConnTimeout(),
		DialContext: (&net.Dialer{
			Timeout:   osCfg.DialTimeout(),
			KeepAlive: osCfg.KeepAlive(),
		}).DialContext,
		TLSHandshakeTimeout: osCfg.TLSHandshakeTimeout(),
	}
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: osCfg.Host(),
		Transport: transport,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenSearch client: %w", err)
	}
	return &ClientImpl{Client: client}, nil
}

func (c *ClientImpl) IndexExists(indexName string) (bool, error) {
	req := opensearchapi.IndicesExistsRequest{
		Index: []string{indexName},
	}
	res, err := req.Do(context.Background(), c.Client)
	if err != nil {
		return false, fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	return res.StatusCode == http.StatusOK, nil
}

func (c *ClientImpl) IndexDocument(ctx context.Context, id string, document interface{}, indexName string) error {
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

	res, err := req.Do(ctx, c.Client)
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

func (c *ClientImpl) Search(ctx context.Context, indexName string, searchBody map[string]interface{}) ([]map[string]interface{}, int, error) {
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

	res, err := req.Do(ctx, c.Client)
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

func (c *ClientImpl) FindDocumentByID(ctx context.Context, indexName, id string) (map[string]interface{}, error) {
	log := logger.NewContextLogger(ctx, "Client/FindDocumentByID")
	req := opensearchapi.GetRequest{
		Index:      indexName,
		DocumentID: id,
	}
	log.Debugf("getting document by id: %s from index: %s", id, indexName)
	res, err := req.Do(ctx, c.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to get document by id: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting document by id: %s", res.String())
	}

	var getResp struct {
		Found  bool                   `json:"found"`
		Source map[string]interface{} `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&getResp); err != nil {
		return nil, fmt.Errorf("failed to decode get response: %w", err)
	}
	if !getResp.Found {
		return nil, fmt.Errorf("document not found")
	}
	return getResp.Source, nil
}
