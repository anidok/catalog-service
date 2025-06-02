package migrate

import (
	"bytes"
	"catalog-service/internal/config"
	"catalog-service/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type Migrate struct {
	client *opensearch.Client
}

func New(client *opensearch.Client) *Migrate {
	return &Migrate{
		client: client,
	}
}

func (m *Migrate) Run(schemaDir string) error {
	files, err := os.ReadDir(schemaDir)
	if err != nil {
		return fmt.Errorf("failed to read schema directory %s: %w", schemaDir, err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		indexName := strings.TrimSuffix(file.Name(), ".json")

		schemaPath := filepath.Join(schemaDir, file.Name())
		schema, err := os.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", file.Name(), err)
		}

		if err := m.createOrUpdateIndex(indexName, schema); err != nil {
			return fmt.Errorf("failed to create/update index %s: %w", indexName, err)
		}

		logger.NonContext.Infof("successfully migrated index: %s\n", indexName)
	}

	return nil
}

func (m *Migrate) createOrUpdateIndex(indexName string, schema []byte) error {
	exists, err := m.indexExists(indexName)
	if err != nil {
		return fmt.Errorf("failed to check if index exists %s: %w", indexName, err)
	}

	if exists {
		logger.NonContext.Infof("index %s already exists. skipping update.", indexName)
		return nil
	}

	return m.createIndex(indexName, schema)
}

func (m *Migrate) indexExists(indexName string) (bool, error) {
	timeout := time.Duration(config.OpenSearch().DialTimeout()) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := opensearchapi.IndicesExistsRequest{
		Index: []string{indexName},
	}
	res, err := req.Do(ctx, m.client)
	if err != nil {
		return false, fmt.Errorf("failed to check index existence for %s: %w", indexName, err)
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

func (m *Migrate) createIndex(indexName string, schema []byte) error {
	var js json.RawMessage
	if err := json.Unmarshal(schema, &js); err != nil {
		return fmt.Errorf("invalid json schema for index %s: %w", indexName, err)
	}

	timeout := time.Duration(config.OpenSearch().DialTimeout()) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := opensearchapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader(schema),
	}
	res, err := req.Do(ctx, m.client)
	if err != nil {
		return fmt.Errorf("failed to create index %s: %w", indexName, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 400 && strings.Contains(res.String(), "resource_already_exists_exception") {
			logger.NonContext.Infof("index %s already exists. skipping creation.", indexName)
			return nil
		}
		return fmt.Errorf("error creating index %s: %s", indexName, res.String())
	}

	return nil
}
