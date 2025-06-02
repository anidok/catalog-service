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

type ServiceAPIDetailIntegrationSuite struct {
	suite.Suite
	server *httptest.Server
	client *opensearch.ClientImpl
	repo   repository.ServiceRepositoryImpl
}

func TestServiceAPIDetailIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServiceAPIDetailIntegrationSuite))
}

func (s *ServiceAPIDetailIntegrationSuite) SetupSuite() {
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

func (s *ServiceAPIDetailIntegrationSuite) TearDownSuite() {
	utils.CleanupTestData(s.client, testconstants.ServiceIndexName, s.T())
	if s.server != nil {
		s.server.Close()
	}
}

func (suite *ServiceAPIDetailIntegrationSuite) Test_GetServiceByID_Success() {
	listResp := suite.doGet("/api/services", nil)
	defer listResp.Body.Close()
	var listResult struct {
		Data struct {
			Services []struct {
				ID string `json:"id"`
			} `json:"services"`
		} `json:"data"`
	}
	suite.decodeResponse(listResp.Body, &listResult)
	svcID := listResult.Data.Services[0].ID

	detailResp := suite.doGet("/api/services/"+svcID, nil)
	defer detailResp.Body.Close()
	assert.Equal(suite.T(), http.StatusOK, detailResp.StatusCode)

	var svcResult dto.ServiceDetailResponse
	suite.decodeResponse(detailResp.Body, &svcResult)
	assert.True(suite.T(), svcResult.Success)
	assert.Equal(suite.T(), svcID, svcResult.Data.ID)
}

func (suite *ServiceAPIDetailIntegrationSuite) Test_GetServiceByID_NotFound() {
	resp := suite.doGet("/api/services/not-existing-id", nil)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)

	var result dto.ServiceDetailResponse
	suite.decodeResponse(resp.Body, &result)
	assert.False(suite.T(), result.Success)
	assert.NotEmpty(suite.T(), result.Errors)
	assert.Equal(suite.T(), "service", result.Errors[0].Entity)
}

func (s *ServiceAPIDetailIntegrationSuite) doGet(path string, headers map[string]string) *http.Response {
	req, err := http.NewRequest("GET", s.server.URL+path, nil)
	s.Require().NoError(err)

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	return resp
}

func (s *ServiceAPIDetailIntegrationSuite) decodeResponse(body io.Reader, out interface{}) {
	decoder := json.NewDecoder(body)
	s.Require().NoError(decoder.Decode(out))
}
