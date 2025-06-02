package validator

import (
	"catalog-service/internal/constants"
	"catalog-service/internal/dto"
	"catalog-service/internal/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceValidatorSuite struct {
	suite.Suite
}

func TestServiceValidatorSuite(t *testing.T) {
	suite.Run(t, new(ServiceValidatorSuite))
}

func (suite *ServiceValidatorSuite) Test_ValidateSearchRequest_Valid() {
	page, limit, errs, code := ValidateSearchRequest("2", "10")
	suite.Equal(2, page)
	suite.Equal(10, limit)
	suite.Empty(errs)
	suite.Equal(http.StatusOK, code)
}

func (suite *ServiceValidatorSuite) Test_ValidateSearchRequest_InvalidPage() {
	page, limit, errs, code := ValidateSearchRequest("0", "10")
	suite.Equal(0, page)
	suite.Equal(10, limit)
	suite.Len(errs, 1)
	suite.Equal(http.StatusBadRequest, code)
	suite.Equal(constants.Error_MALFORMED_DATA, errs[0].Code)
	suite.Equal("page", errs[0].Entity)
}

func (suite *ServiceValidatorSuite) Test_ValidateSearchRequest_InvalidLimit() {
	page, limit, errs, code := ValidateSearchRequest("1", "0")
	suite.Equal(1, page)
	suite.Equal(0, limit)
	suite.Len(errs, 1)
	suite.Equal(http.StatusBadRequest, code)
	suite.Equal(constants.Error_MALFORMED_DATA, errs[0].Code)
	suite.Equal("limit", errs[0].Entity)
}

func (suite *ServiceValidatorSuite) Test_ValidateSearchRequest_BothInvalid() {
	page, limit, errs, code := ValidateSearchRequest("0", "0")
	suite.Equal(0, page)
	suite.Equal(0, limit)
	suite.Len(errs, 2)
	suite.Equal(http.StatusBadRequest, code)
	suite.Equal("page", errs[0].Entity)
	suite.Equal("limit", errs[1].Entity)
}

func (suite *ServiceValidatorSuite) Test_ValidateSearchRequest_NonInt() {
	page, limit, errs, code := ValidateSearchRequest("abc", "xyz")
	suite.Equal(0, page)
	suite.Equal(0, limit)
	suite.Len(errs, 2)
	suite.Equal(http.StatusBadRequest, code)
	suite.Equal("page", errs[0].Entity)
	suite.Equal("limit", errs[1].Entity)
}

func TestValidateCreateRequest_Valid(t *testing.T) {
	req := &dto.ServiceDTO{
		Name: "Test Service",
		Versions: []models.Version{
			{VersionNumber: "1.0", Details: "Initial"},
		},
	}
	errs, code := ValidateCreateRequest(req)
	assert.Empty(t, errs)
	assert.Equal(t, 200, code)
}

func TestValidateCreateRequest_MissingName(t *testing.T) {
	req := &dto.ServiceDTO{
		Name: "",
		Versions: []models.Version{
			{VersionNumber: "1.0", Details: "Initial"},
		},
	}
	errs, code := ValidateCreateRequest(req)
	assert.Len(t, errs, 1)
	assert.Equal(t, 400, code)
	assert.Equal(t, "name", errs[0].Entity)
}

func TestValidateCreateRequest_MissingVersions(t *testing.T) {
	req := &dto.ServiceDTO{
		Name:     "Test Service",
		Versions: []models.Version{},
	}
	errs, code := ValidateCreateRequest(req)
	assert.Len(t, errs, 1)
	assert.Equal(t, 400, code)
	assert.Equal(t, "versions", errs[0].Entity)
}

func TestValidateCreateRequest_MissingVersionNumber(t *testing.T) {
	req := &dto.ServiceDTO{
		Name: "Test Service",
		Versions: []models.Version{
			{VersionNumber: "", Details: "Initial"},
		},
	}
	errs, code := ValidateCreateRequest(req)
	assert.Len(t, errs, 1)
	assert.Equal(t, 400, code)
	assert.Equal(t, "versions", errs[0].Entity)
	assert.Contains(t, errs[0].Cause, "version_number is required")
}
