package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"catalog-service/internal/api"
	"catalog-service/internal/config"
	"catalog-service/internal/logger"
	"catalog-service/internal/opensearch"
	"catalog-service/internal/repository"
	testconstants "catalog-service/test/constants"
	"catalog-service/test/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceAPIDeleteIntegrationSuite struct {
	suite.Suite
	server *httptest.Server
	client *opensearch.ClientImpl
	repo   repository.ServiceRepositoryImpl
}

func TestServiceAPIDeleteIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServiceAPIDeleteIntegrationSuite))
}

func (s *ServiceAPIDeleteIntegrationSuite) SetupSuite() {
	config.Load()
	logger.Setup("INFO", "json")

	client, err := opensearch.NewClient(config.OpenSearch().Host())
	s.Require().NoError(err)
	s.client = client
	s.repo = repository.ServiceRepositoryImpl{Client: client}

	utils.CleanupTestData(s.client, testconstants.ServiceIndexName, s.T())
	utils.LoadTestData(s.repo, s.T())

	s.server = httptest.NewServer(api.NewRouter(&s.repo))
}

func (s *ServiceAPIDeleteIntegrationSuite) TearDownSuite() {
	utils.CleanupTestData(s.client, testconstants.ServiceIndexName, s.T())
	if s.server != nil {
		s.server.Close()
	}
}

func (suite *ServiceAPIDeleteIntegrationSuite) Test_DeleteServiceByID_Success() {
	payload := map[string]interface{}{
		"name":        "Delete Test Service",
		"description": "To be deleted",
		"versions": []map[string]interface{}{
			{"version_number": "1.0", "Details": "Initial"},
		},
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(suite.server.URL+"/api/services", "application/json", bytes.NewReader(body))
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	suite.decodeResponse(resp.Body, &result)
	assert.True(suite.T(), result.Success)
	id := result.Data.ID

	// Now delete
	req, _ := http.NewRequest("DELETE", suite.server.URL+"/api/services/"+id, nil)
	delResp, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)
	defer delResp.Body.Close()
	assert.Equal(suite.T(), http.StatusOK, delResp.StatusCode)

	var delResult struct {
		Success bool `json:"success"`
	}
	suite.decodeResponse(delResp.Body, &delResult)
	assert.True(suite.T(), delResult.Success)

	// Confirm not found
	getResp := suite.doGet("/api/services/"+id, nil)
	defer getResp.Body.Close()
	assert.Equal(suite.T(), http.StatusNotFound, getResp.StatusCode)
}

func (s *ServiceAPIDeleteIntegrationSuite) doGet(path string, headers map[string]string) *http.Response {
	req, err := http.NewRequest("GET", s.server.URL+path, nil)
	s.Require().NoError(err)

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	return resp
}

func (s *ServiceAPIDeleteIntegrationSuite) decodeResponse(body io.Reader, out interface{}) {
	decoder := json.NewDecoder(body)
	s.Require().NoError(decoder.Decode(out))
}
