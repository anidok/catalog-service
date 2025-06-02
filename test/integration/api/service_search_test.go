package api_test

import (
	"encoding/json"
	"io"
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

type ServiceAPISearchIntegrationSuite struct {
	suite.Suite
	server *httptest.Server
	client *opensearch.ClientImpl
	repo   repository.ServiceRepositoryImpl
}

func TestServiceAPISearchIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServiceAPISearchIntegrationSuite))
}

func (s *ServiceAPISearchIntegrationSuite) SetupSuite() {
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

func (s *ServiceAPISearchIntegrationSuite) TearDownSuite() {
	utils.CleanupTestData(s.client, testconstants.ServiceIndexName, s.T())
	if s.server != nil {
		s.server.Close()
	}
}

func (s *ServiceAPISearchIntegrationSuite) Test_SearchAPI_EmptyQuery_DefaultPagination() {
	resp := s.doGet("/api/services", nil)
	defer resp.Body.Close()
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var result dto.ServiceListResponse
	s.decodeResponse(resp.Body, &result)

	assert.GreaterOrEqual(s.T(), result.Data.Count, 5)
	assert.Equal(s.T(), 10, len(result.Data.Services))
	assert.Nil(s.T(), result.Errors)
	assert.NotNil(s.T(), result.Data.Next)
}

func (s *ServiceAPISearchIntegrationSuite) Test_SearchAPI_QueryPhrase() {
	resp := s.doGet("/api/services?q=special%20rates", nil)
	defer resp.Body.Close()
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var result dto.ServiceListResponse
	s.decodeResponse(resp.Body, &result)

	assert.Equal(s.T(), 1, result.Data.Count)
	assert.Len(s.T(), result.Data.Services, 1)
	assert.Equal(s.T(), "Forex Card", result.Data.Services[0].Name)
	assert.Contains(s.T(), result.Data.Services[0].Description, "special rates")
	assert.Nil(s.T(), result.Errors)
	assert.Nil(s.T(), result.Data.Next)
}

func (s *ServiceAPISearchIntegrationSuite) Test_SearchAPI_InvalidPage() {
	resp := s.doGet("/api/services?page=0", nil)
	defer resp.Body.Close()
	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	var result dto.ServiceListResponse
	s.decodeResponse(resp.Body, &result)

	assert.Nil(s.T(), result.Data)
	assert.NotEmpty(s.T(), result.Errors)
	assert.Equal(s.T(), "page", result.Errors[0].Entity)
	assert.Equal(s.T(), "invalid page", result.Errors[0].Cause)
	assert.Equal(s.T(), "101", result.Errors[0].Code)
}

func (s *ServiceAPISearchIntegrationSuite) doGet(path string, headers map[string]string) *http.Response {
	req, err := http.NewRequest("GET", s.server.URL+path, nil)
	s.Require().NoError(err)

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	return resp
}

func (s *ServiceAPISearchIntegrationSuite) decodeResponse(body io.Reader, out interface{}) {
	decoder := json.NewDecoder(body)
	s.Require().NoError(decoder.Decode(out))
}
