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

const TestIndexName = "test-index"

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

	exists, err := client.IndexExists(TestIndexName)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *ClientTestSuite) Test_IndexExists_ReturnsFalseWhenNotFound() {
	client := newMockClient(nil, http.StatusNotFound)

	exists, err := client.IndexExists(TestIndexName)

	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *ClientTestSuite) Test_IndexDocument_SuccessfulIndexing() {
	body := `{
		"_id": "doc123",
		"result": "created"
	}`
	client := newMockClient(unmarshalJSON(body), http.StatusCreated)

	err := client.IndexDocument(context.Background(), "doc123", map[string]string{"foo": "bar"}, TestIndexName)

	assert.NoError(suite.T(), err)
}

func (suite *ClientTestSuite) Test_IndexDocument_FailureOnBadRequest() {
	body := `{
		"error": "some error"
	}`
	client := newMockClient(unmarshalJSON(body), http.StatusBadRequest)

	err := client.IndexDocument(context.Background(), "doc123", map[string]string{"foo": "bar"}, TestIndexName)

	assert.Error(suite.T(), err)
}

func (suite *ClientTestSuite) Test_Search_Successful() {
	responseJSON := `{
		"hits": {
			"total": { "value": 2 },
			"hits": [
				{ "_source": { "name": "Locate Us", "description": "Find our nearest branch" } },
				{ "_source": { "name": "Contact Us", "description": "Reach out to our support team" } }
			]
		}
	}`
	client := newMockClient(unmarshalJSON(responseJSON), http.StatusOK)

	requestJSON := `{
		"query": {
			"multi_match": {
				"query": "us",
				"fields": ["name", "description"]
			}
		},
		"from": 0,
		"size": 10
	}`
	searchBody := unmarshalJSON(requestJSON)

	hits, total, err := client.Search(context.Background(), "services", searchBody)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, total)
	assert.Len(suite.T(), hits, 2)
	assert.Equal(suite.T(), "Locate Us", hits[0]["name"])
	assert.Equal(suite.T(), "Contact Us", hits[1]["name"])
}

func (suite *ClientTestSuite) Test_Search_Error() {
	client := newMockClient(nil, http.StatusBadRequest)

	requestJSON := `{
		"query": {
			"multi_match": {
				"query": "fail",
				"fields": ["name", "description"]
			}
		}
	}`
	searchBody := unmarshalJSON(requestJSON)

	hits, total, err := client.Search(context.Background(), "services", searchBody)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), hits)
	assert.Equal(suite.T(), 0, total)
}

func (suite *ClientTestSuite) Test_FindDocumentByID_Success() {
	body := `{
		"found": true,
		"_source": {
			"id": "finddoc-id",
			"name": "FindDoc",
			"description": "FindDoc Desc"
		}
	}`
	client := newMockClient(unmarshalJSON(body), http.StatusOK)
	ctx := context.Background()
	found, err := client.FindDocumentByID(ctx, TestIndexName, "finddoc-id")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), found)
	assert.Equal(suite.T(), "FindDoc", found["name"])
}

func (suite *ClientTestSuite) Test_FindDocumentByID_NotFound() {
	body := `{
		"found": false
	}`
	client := newMockClient(unmarshalJSON(body), http.StatusOK)
	ctx := context.Background()
	_, err := client.FindDocumentByID(ctx, TestIndexName, "not-exist-id")
	assert.Error(suite.T(), err)
}

func (suite *ClientTestSuite) Test_DeleteDocumentByID_Success() {
	body := `{
		"result": "deleted"
	}`
	client := newMockClient(unmarshalJSON(body), http.StatusOK)
	ctx := context.Background()
	err := client.DeleteDocumentByID(ctx, TestIndexName, "delete-id")
	assert.NoError(suite.T(), err)
}

func (suite *ClientTestSuite) Test_DeleteDocumentByID_NotFound() {
	body := `{
		"error": "not found"
	}`
	client := newMockClient(unmarshalJSON(body), http.StatusNotFound)
	ctx := context.Background()
	err := client.DeleteDocumentByID(ctx, TestIndexName, "not-exist-id")
	assert.Error(suite.T(), err)
}