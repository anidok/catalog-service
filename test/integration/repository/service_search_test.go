package repository_test

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"catalog-service/internal/config"
	"catalog-service/internal/logger"
	"catalog-service/internal/models"
	"catalog-service/internal/opensearch"
	"catalog-service/internal/repository"
	"catalog-service/test/utils"

	"runtime"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

const ServiceIndexName = "services"

type ServiceRepoSearchIntegrationSuite struct {
	suite.Suite
	client *opensearch.ClientImpl
	repo   repository.ServiceRepositoryImpl
}

func TestServiceRepoSearchIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServiceRepoSearchIntegrationSuite))
}

func (suite *ServiceRepoSearchIntegrationSuite) SetupSuite() {
	suite.T().Parallel()
	config.Load()
	logger.Setup("INFO", "json")
	client, err := opensearch.NewClient(config.OpenSearch().Host())
	suite.Require().NoError(err)
	suite.client = client
	suite.cleanupTestData()
	suite.repo = repository.ServiceRepositoryImpl{Client: client}
	suite.loadTestData()
}

func (suite *ServiceRepoSearchIntegrationSuite) TearDownSuite() {
	suite.cleanupTestData()
}

func (suite *ServiceRepoSearchIntegrationSuite) Test_Search_Pagination() {
	tests := []struct {
		name     string
		query    string
		page     int
		limit    int
		wantLen  int
		minTotal int
		expected []models.Service
	}{
		{
			name:     "Given_EmptyQuery_When_FirstPage_Then_Returns5Items",
			page:     1,
			limit:    5,
			wantLen:  5,
			minTotal: 40,
			expected: []models.Service{
				suite.buildService("Loan Calculator", "Calculate your loan EMIs", "1.0"),
				suite.buildService("Investment Tracker", "Track your investments and get insights", "1.0"),
				suite.buildService("Tax Filing", "File your taxes online", "1.0"),
				suite.buildService("Travel Insurance", "Buy travel insurance instantly", "1.0"),
				suite.buildService("Account Statement", "Download your account statements", "1.0"),
			},
		},
		{
			name:     "Given_EmptyQuery_When_Page3_Then_Returns4Items",
			page:     3,
			limit:    4,
			wantLen:  4,
			minTotal: 40,
			expected: []models.Service{
				suite.buildService("Smart Alerts", "Receive smart alerts for account activity", "1.0"),
				suite.buildService("AI Insights", "Get AI-driven insights for your finances", "1.0"),
				suite.buildService("ESG Investing", "Invest in companies with strong ESG practices", "1.0"),
				suite.buildService("Crypto Services", "Buy, sell, and store cryptocurrencies", "1.0"),
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			ctx := context.Background()
			services, total, err := suite.repo.Search(ctx, tt.query, tt.page, tt.limit)
			suite.assertSearchResults(services, total, err, tt.wantLen, tt.minTotal)

			for i, svc := range services {
				suite.assertService(*svc, tt.expected[i])
			}
		})
	}
}

func (suite *ServiceRepoSearchIntegrationSuite) Test_Search_Queries() {
	tests := []struct {
		name     string
		query    string
		wantLen  int
		minTotal int
		expected []models.Service
	}{
		{
			name:     "Given_BankingQuery_Then_ReturnsMatchingServices",
			query:    "banking",
			wantLen:  9,
			minTotal: 9,
			expected: []models.Service{
				suite.buildService("Voice Banking", "Bank using voice commands", "1.0"),
				suite.buildService("Premium Banking", "Get exclusive offers and rewards", "1.0"),
				suite.buildService("Digital Banking", "Open your account online", "1.0"),
			},
		},
		{
			name:     "Given_ExactPhraseQuery_Then_ReturnsExactMatch",
			query:    "special rates",
			wantLen:  1,
			minTotal: 1,
			expected: []models.Service{
				suite.buildServiceWithVersion("Forex Card", "Forex card for students with special rates", "2.0", "Student offer"),
			},
		},
		{
			name:     "Given_NonExistentQuery_Then_ReturnsEmpty",
			query:    "thisshouldnotmatchanything",
			wantLen:  0,
			minTotal: 0,
			expected: nil,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			ctx := context.Background()
			services, total, err := suite.repo.Search(ctx, tt.query, 1, 10)
			suite.assertSearchResults(services, total, err, tt.wantLen, tt.minTotal)
		})
	}
}

func (suite *ServiceRepoSearchIntegrationSuite) Test_Search_Error_InvalidPageLimit() {
	ctx := context.Background()
	services, total, err := suite.repo.Search(ctx, "", -1, -10)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), services)
	assert.Equal(suite.T(), 0, total)
}

func (suite *ServiceRepoSearchIntegrationSuite) loadTestData() {
	ctx := context.Background()
	_, filename, _, _ := runtime.Caller(0)
	testdataPath := filepath.Join(filepath.Dir(filename), "..", "testdata", "services.json")

	services, err := utils.UnmarshalServiceList(testdataPath)
	suite.Require().NoError(err)

	for _, svc := range services {
		err := suite.repo.Create(ctx, svc)
		suite.Require().NoError(err)
	}
}

func (suite *ServiceRepoSearchIntegrationSuite) cleanupTestData() {
	ctx := context.Background()
	deleteByQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	body, _ := json.Marshal(deleteByQuery)
	req := opensearchapi.DeleteByQueryRequest{
		Index: []string{ServiceIndexName},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, suite.client)
	if err != nil {
		suite.T().Logf("Cleanup delete by query failed: %v", err)
		return
	}
	defer res.Body.Close()
}

func (suite *ServiceRepoSearchIntegrationSuite) buildService(name, desc string, version string) models.Service {
	return models.Service{
		Name:        name,
		Description: desc,
		Versions:    []models.Version{{VersionNumber: version, Details: "Initial release"}},
	}
}

func (suite *ServiceRepoSearchIntegrationSuite) buildServiceWithVersion(name, desc, version, details string) models.Service {
	return models.Service{
		Name:        name,
		Description: desc,
		Versions:    []models.Version{{VersionNumber: version, Details: details}},
	}
}

func (suite *ServiceRepoSearchIntegrationSuite) assertSearchResults(services []*models.Service, total int, err error, wantLen int, minTotal int) {
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), wantLen, len(services))
	assert.GreaterOrEqual(suite.T(), total, minTotal)
}

func (suite *ServiceRepoSearchIntegrationSuite) assertService(actual models.Service, expected models.Service) {
	assert.Equal(suite.T(), expected.Name, actual.Name)
	assert.Equal(suite.T(), expected.Description, actual.Description)

	assert.Len(suite.T(), actual.Versions, len(expected.Versions))
	for i, v := range actual.Versions {
		assert.Equal(suite.T(), expected.Versions[i].VersionNumber, v.VersionNumber)
		assert.Equal(suite.T(), expected.Versions[i].Details, v.Details)
	}
}
