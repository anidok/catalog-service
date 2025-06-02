package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"

	"catalog-service/internal/opensearch"
	"catalog-service/internal/repository"

	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/stretchr/testify/require"
)

func LoadTestData(repo repository.ServiceRepositoryImpl, t *testing.T) {
	ctx := context.Background()
	_, filename, _, _ := runtime.Caller(0)
	testdataPath := filepath.Join(filepath.Dir(filename), "..", "integration", "testdata", "services.json")

	services, err := UnmarshalServiceList(testdataPath)
	require.NoError(t, err)

	for _, svc := range services {
		err := repo.Create(ctx, svc)
		require.NoError(t, err)
	}
}

func CleanupTestData(client *opensearch.ClientImpl, indexName string, t *testing.T) {
	if client == nil {
		return
	}
	ctx := context.Background()
	deleteByQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	body, _ := json.Marshal(deleteByQuery)
	req := opensearchapi.DeleteByQueryRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		t.Logf("Cleanup delete by query failed: %v", err)
		return
	}
	defer res.Body.Close()
}
