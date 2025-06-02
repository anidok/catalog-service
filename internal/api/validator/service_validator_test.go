package validator

import (
	"catalog-service/internal/constants"
	"net/http"
	"testing"

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
