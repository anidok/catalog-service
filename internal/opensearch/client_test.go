package opensearch

import (
	"context"
	"net/http"
	"testing"

	"catalog-service/internal/config"
	"catalog-service/internal/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
}

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (suite *ClientTestSuite) SetupTest() {
	config.Load()
	logger.Setup("INFO", "json")
}

func (suite *ClientTestSuite) Test_IndexExists_ReturnsTrueWhenStatusOK() {
	client := newMockClient(nil, http.StatusOK)

	exists, err := client.IndexExists("test-index")

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *ClientTestSuite) Test_IndexExists_ReturnsFalseWhenNotFound() {
	client := newMockClient(nil, http.StatusNotFound)

	exists, err := client.IndexExists("test-index")

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *ClientTestSuite) Test_IndexDocument_SuccessfulIndexing() {
	respBody := map[string]interface{}{
		"_id":    "doc123",
		"result": "created",
	}
	client := newMockClient(respBody, http.StatusCreated)

	err := client.IndexDocument(context.Background(), "doc123", map[string]string{"foo": "bar"}, "test-index")

	assert.NoError(suite.T(), err)
}

func (suite *ClientTestSuite) Test_IndexDocument_FailureOnBadRequest() {
	respBody := map[string]interface{}{
		"error": "some error",
	}
	client := newMockClient(respBody, http.StatusBadRequest)

	err := client.IndexDocument(context.Background(), "doc123", map[string]string{"foo": "bar"}, "test-index")

	assert.Error(suite.T(), err)
}
