package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"catalog-service/internal/api"
	"catalog-service/internal/config"
	"catalog-service/internal/dto"
	"catalog-service/internal/logger"
	"catalog-service/internal/opensearch"
	"catalog-service/internal/repository"
	testconstants "catalog-service/test/constants"
	"catalog-service/test/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceAPICreateIntegrationSuite struct {
	suite.Suite
	server *httptest.Server
	client *opensearch.ClientImpl
	repo   repository.ServiceRepositoryImpl
	url    string
}

func TestServiceAPICreateIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServiceAPICreateIntegrationSuite))
}

func (s *ServiceAPICreateIntegrationSuite) SetupSuite() {
	config.Load()
	logger.Setup("INFO", "json")

	client, err := opensearch.NewClient(config.OpenSearch().Host())
	s.Require().NoError(err)
	s.client = client
	s.repo = repository.ServiceRepositoryImpl{Client: client}

	utils.CleanupTestData(s.client, testconstants.ServiceIndexName, s.T())
	utils.LoadTestData(s.repo, s.T())

	s.server = httptest.NewServer(api.NewRouter(&s.repo))
	s.url = s.server.URL
}

func (s *ServiceAPICreateIntegrationSuite) TearDownSuite() {
	utils.CleanupTestData(s.client, testconstants.ServiceIndexName, s.T())
	if s.server != nil {
		s.server.Close()
	}
}

func (suite *ServiceAPICreateIntegrationSuite) Test_CreateService_Success() {
	payload := map[string]interface{}{
		"name":        "Integration Test Service",
		"description": "Created via integration test",
		"versions": []map[string]interface{}{
			{"version_number": "1.0", "Details": "Initial"},
		},
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(suite.url+"/api/services", "application/json", bytes.NewReader(body))
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var result dto.ServiceDetailResponse
	decoder := json.NewDecoder(resp.Body)
	suite.Require().NoError(decoder.Decode(&result))
	assert.True(suite.T(), result.Success)
	assert.Equal(suite.T(), "Integration Test Service", result.Data.Name)
	assert.Equal(suite.T(), "Created via integration test", result.Data.Description)
}

func (suite *ServiceAPICreateIntegrationSuite) Test_CreateService_InvalidBody() {
	resp, err := http.Post(suite.url+"/api/services", "application/json", bytes.NewReader([]byte(`{invalid json}`)))
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	var result dto.ServiceDetailResponse
	decoder := json.NewDecoder(resp.Body)
	suite.Require().NoError(decoder.Decode(&result))
	assert.False(suite.T(), result.Success)
	assert.NotEmpty(suite.T(), result.Errors)
	assert.Equal(suite.T(), "service", result.Errors[0].Entity)
}
