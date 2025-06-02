package repository_test

import (
	"context"
	"testing"

	"catalog-service/internal/config"
	"catalog-service/internal/logger"
	"catalog-service/internal/models"
	"catalog-service/internal/opensearch"
	"catalog-service/internal/repository"
	testconstants "catalog-service/test/constants"
	"catalog-service/test/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceRepoSearchIntegrationSuite struct {
	suite.Suite
	client *opensearch.ClientImpl
	repo   repository.ServiceRepositoryImpl
}

func TestServiceRepoSearchIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServiceRepoSearchIntegrationSuite))
}

func (suite *ServiceRepoSearchIntegrationSuite) SetupSuite() {
	config.Load()
	logger.Setup("INFO", "json")
	client, err := opensearch.NewClient(config.OpenSearch().Host())
	suite.Require().NoError(err)
	suite.client = client
	suite.repo = repository.ServiceRepositoryImpl{Client: client}
	utils.CleanupTestData(suite.client, testconstants.ServiceIndexName, suite.T())
	utils.LoadTestData(suite.repo, suite.T())
}

func (suite *ServiceRepoSearchIntegrationSuite) TearDownSuite() {
	utils.CleanupTestData(suite.client, testconstants.ServiceIndexName, suite.T())
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
