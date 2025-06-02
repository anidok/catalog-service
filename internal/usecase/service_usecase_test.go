package usecase

import (
	"context"
	"testing"
	"time"

	"catalog-service/internal/constants"
	"catalog-service/internal/dto"
	"catalog-service/internal/models"
	mockrepo "catalog-service/test/mocks/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceUsecaseSuite struct {
	suite.Suite
}

func TestServiceUsecaseSuite(t *testing.T) {
	suite.Run(t, new(ServiceUsecaseSuite))
}

func (suite *ServiceUsecaseSuite) Test_Search_Success() {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	mockRepo := new(mockrepo.ServiceRepository)
	mockRepo.
		On("Search", mock.Anything, "", 1, 10).
		Return([]*models.Service{
			{
				ID:          "id1",
				Name:        "Service1",
				Description: "Desc1",
				Versions:    []models.Version{{VersionNumber: "1.0", Details: "Initial"}},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}, 1, nil)

	uc := NewServiceUsecase(mockRepo)
	dtos, total, err := uc.Search(context.Background(), "", 1, 10)

	suite.Require().NoError(err)
	suite.Require().Equal(1, total)
	suite.Require().Len(dtos, 1)

	want := struct {
		ID, Name, Description, VersionNumber, Details, CreatedAt, UpdatedAt string
	}{
		ID:            "id1",
		Name:          "Service1",
		Description:   "Desc1",
		VersionNumber: "1.0",
		Details:       "Initial",
		CreatedAt:     now.Format(constants.Iso8601Format),
		UpdatedAt:     now.Format(constants.Iso8601Format),
	}

	suite.assertServiceDTOEqual(dtos[0], want)
	mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceUsecaseSuite) Test_Search_Error() {
	mockRepo := new(mockrepo.ServiceRepository)
	mockRepo.
		On("Search", mock.Anything, "", 1, 10).
		Return(nil, 0, assert.AnError)

	uc := NewServiceUsecase(mockRepo)
	dtos, total, err := uc.Search(context.Background(), "", 1, 10)
	suite.Error(err)
	suite.Nil(dtos)
	suite.Equal(0, total)
	mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceUsecaseSuite) assertServiceDTOEqual(got *dto.ServiceDTO, want struct {
	ID, Name, Description, VersionNumber, Details, CreatedAt, UpdatedAt string
}) {
	suite.Equal(want.ID, got.ID)
	suite.Equal(want.Name, got.Name)
	suite.Equal(want.Description, got.Description)
	suite.Require().Len(got.Versions, 1)
	suite.Equal(want.VersionNumber, got.Versions[0].VersionNumber)
	suite.Equal(want.Details, got.Versions[0].Details)
	suite.Equal(want.CreatedAt, got.CreatedAt)
	suite.Equal(want.UpdatedAt, got.UpdatedAt)
}
