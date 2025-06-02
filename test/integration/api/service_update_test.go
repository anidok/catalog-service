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
	"catalog-service/internal/dto"
	"catalog-service/internal/logger"
	"catalog-service/internal/opensearch"
	"catalog-service/internal/repository"
	testconstants "catalog-service/test/constants"
	"catalog-service/test/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceAPIUpdateIntegrationSuite struct {
	suite.Suite
	server *httptest.Server
	client *opensearch.ClientImpl
	repo   repository.ServiceRepositoryImpl
}

func TestServiceAPIUpdateIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServiceAPIUpdateIntegrationSuite))
}

func (s *ServiceAPIUpdateIntegrationSuite) SetupSuite() {
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

func (s *ServiceAPIUpdateIntegrationSuite) TearDownSuite() {
	utils.CleanupTestData(s.client, testconstants.ServiceIndexName, s.T())
	if s.server != nil {
		s.server.Close()
	}
}

func (suite *ServiceAPIUpdateIntegrationSuite) Test_UpdateService_AppendVersionAndDescription() {
	// Get a service ID
	resp := suite.doGet("/api/services", nil)
	defer resp.Body.Close()
	var listResult struct {
		Data struct {
			Services []struct {
				ID          string `json:"id"`
				Description string `json:"description"`
				Versions    []struct {
					VersionNumber string `json:"VersionNumber"`
					Details       string `json:"Details"`
				} `json:"versions"`
			} `json:"services"`
		} `json:"data"`
	}
	suite.decodeResponse(resp.Body, &listResult)
	svc := listResult.Data.Services[0]
	svcID := svc.ID

	// Update description and append a version
	updatePayload := map[string]interface{}{
		"description": "Updated description",
		"versions": []map[string]interface{}{
			{"version_number": "2.0", "Details": "Second release"},
		},
	}
	body, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest("PUT", suite.server.URL+"/api/services/"+svcID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	updateResp, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)
	defer updateResp.Body.Close()
	assert.Equal(suite.T(), http.StatusOK, updateResp.StatusCode)

	var updateResult dto.ServiceDetailResponse
	suite.decodeResponse(updateResp.Body, &updateResult)
	assert.True(suite.T(), updateResult.Success)
	assert.Equal(suite.T(), "Updated description", updateResult.Data.Description)
	found := false
	for _, v := range updateResult.Data.Versions {
		if v.VersionNumber == "2.0" && v.Details == "Second release" {
			found = true
			break
		}
	}
	assert.True(suite.T(), found, "new version should be appended")

	// Try to update name (should fail)
	updatePayload = map[string]interface{}{
		"name": "Should Not Update Name",
	}
	body, _ = json.Marshal(updatePayload)
	req, _ = http.NewRequest("PUT", suite.server.URL+"/api/services/"+svcID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	updateResp, err = http.DefaultClient.Do(req)
	suite.Require().NoError(err)
	defer updateResp.Body.Close()
	assert.Equal(suite.T(), http.StatusBadRequest, updateResp.StatusCode)

	var failResult dto.ServiceDetailResponse
	suite.decodeResponse(updateResp.Body, &failResult)
	assert.False(suite.T(), failResult.Success)
	assert.NotEmpty(suite.T(), failResult.Errors)
	assert.Equal(suite.T(), "name", failResult.Errors[0].Entity)
}

func (s *ServiceAPIUpdateIntegrationSuite) doGet(path string, headers map[string]string) *http.Response {
	req, err := http.NewRequest("GET", s.server.URL+path, nil)
	s.Require().NoError(err)

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	return resp
}

func (s *ServiceAPIUpdateIntegrationSuite) decodeResponse(body io.Reader, out interface{}) {
	decoder := json.NewDecoder(body)
	s.Require().NoError(decoder.Decode(out))
}
