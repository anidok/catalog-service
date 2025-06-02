package repository

import (
	"context"
	"testing"
	"time"

	"catalog-service/internal/config"
	"catalog-service/internal/logger"
	"catalog-service/internal/models"
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

func (suite *ServiceRepoTestSuite) Test_FindByID_Success() {
	mockClient := new(opensearchmock.Client)
	ctx := context.Background()
	svc := &models.Service{
		ID:          "test-find-id",
		Name:        "FindByID Service",
		Description: "Test FindByID",
		Versions:    []models.Version{{VersionNumber: "1.0", Details: "Initial"}},
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	mockClient.On("FindDocumentByID", mock.Anything, "services", "test-find-id").Return(
		map[string]interface{}{
			"id":          svc.ID,
			"name":        svc.Name,
			"description": svc.Description,
			"versions":    []map[string]interface{}{{"VersionNumber": "1.0", "Details": "Initial"}},
			"created_at":  svc.CreatedAt,
			"updated_at":  svc.UpdatedAt,
		}, nil,
	)

	repo := &ServiceRepositoryImpl{Client: mockClient}

	found, err := repo.FindByID(ctx, "test-find-id")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), found)
	assert.Equal(suite.T(), svc.ID, found.ID)
	assert.Equal(suite.T(), svc.Name, found.Name)
	assert.Equal(suite.T(), svc.Description, found.Description)
}

func (suite *ServiceRepoTestSuite) Test_FindByID_NotFound() {
	mockClient := new(opensearchmock.Client)
	ctx := context.Background()
	mockClient.On("FindDocumentByID", mock.Anything, "services", "does-not-exist-id").Return(nil, assert.AnError)

	repo := &ServiceRepositoryImpl{Client: mockClient}

	_, err := repo.FindByID(ctx, "does-not-exist-id")
	assert.Error(suite.T(), err)
}
