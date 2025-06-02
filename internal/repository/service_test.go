package repository

import (
	"context"
	"testing"

	"catalog-service/internal/config"
	"catalog-service/internal/logger"
	opensearchmock "catalog-service/test/mocks/opensearch"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceRepoTestSuite struct {
	suite.Suite
}

func TestClient(t *testing.T) {
	suite.Run(t, new(ServiceRepoTestSuite))
}

func (suite *ServiceRepoTestSuite) SetupTest() {
	config.Load()
	logger.Setup("INFO", "json")
}

func (suite *ServiceRepoTestSuite) Test_Search_Success() {
	mockClient := new(opensearchmock.Client)
	mockClient.On("Search", mock.Anything, "services", mock.Anything).Return(
		[]map[string]interface{}{
			{
				"name":        "Locate Us",
				"description": "Find our nearest branch",
			},
			{
				"name":        "Contact Us",
				"description": "Reach out to our support team",
			},
		}, 2, nil,
	)

	repo := &ServiceRepositoryImpl{Client: mockClient}

	services, total, err := repo.Search(context.Background(), "us", 1, 10)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, total)
	assert.Len(suite.T(), services, 2)
	assert.Equal(suite.T(), "Locate Us", services[0].Name)
	assert.Equal(suite.T(), "Contact Us", services[1].Name)
}

func (suite *ServiceRepoTestSuite) Test_Search_Error() {
	mockClient := new(opensearchmock.Client)
	mockClient.On("Search", mock.Anything, "services", mock.Anything).Return(
		nil, 0, assert.AnError,
	)

	repo := &ServiceRepositoryImpl{Client: mockClient}

	services, total, err := repo.Search(context.Background(), "fail", 1, 10)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), services)
	assert.Equal(suite.T(), 0, total)
}
