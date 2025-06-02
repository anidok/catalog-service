package validator

import (
	"net/http"
	"strconv"

	"catalog-service/internal/constants"
	"catalog-service/internal/dto"
)

func ValidateSearchRequest(pageStr, limitStr string) (page int, limit int, errs []dto.ErrorObj, httpCode int) {
	var errors []dto.ErrorObj
	httpCode = http.StatusOK

	page, pageErr, pageCode := validatePageWithError(pageStr)
	if pageErr != nil {
		errors = append(errors, *pageErr)
		httpCode = pageCode
	}
	limit, limitErr, limitCode := validateLimitWithError(limitStr)
	if limitErr != nil {
		errors = append(errors, *limitErr)
		if httpCode == http.StatusOK {
			httpCode = limitCode
		}
	}
	return page, limit, errors, httpCode
}

func validatePageWithError(pageStr string) (int, *dto.ErrorObj, int) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, &dto.ErrorObj{
			Code:   constants.Error_MALFORMED_DATA,
			Entity: "page",
			Cause:  "invalid page",
		}, http.StatusBadRequest
	}
	return page, nil, http.StatusOK
}

func validateLimitWithError(limitStr string) (int, *dto.ErrorObj, int) {
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return 0, &dto.ErrorObj{
			Code:   constants.Error_MALFORMED_DATA,
			Entity: "limit",
			Cause:  "invalid limit",
		}, http.StatusBadRequest
	}
	return limit, nil, http.StatusOK
}

func ValidateID(id string) ([]dto.ErrorObj, int) {
	if id == "" {
		return []dto.ErrorObj{{
			Code:   constants.Error_MALFORMED_DATA,
			Entity: "id",
			Cause:  "missing id",
		}}, http.StatusBadRequest
	}
	return nil, http.StatusOK
}

func ValidateCreateRequest(req *dto.ServiceDTO) ([]dto.ErrorObj, int) {
	var errs []dto.ErrorObj
	if req.Name == "" {
		errs = append(errs, dto.ErrorObj{
			Code:   constants.Error_MALFORMED_DATA,
			Entity: "name",
			Cause:  "name is required",
		})
	}
	if len(req.Versions) == 0 {
		errs = append(errs, dto.ErrorObj{
			Code:   constants.Error_MALFORMED_DATA,
			Entity: "versions",
			Cause:  "at least one version is required",
		})
	}
	for i, v := range req.Versions {
		if v.VersionNumber == "" {
			errs = append(errs, dto.ErrorObj{
				Code:   constants.Error_MALFORMED_DATA,
				Entity: "versions",
				Cause:  "version_number is required for version at index " + strconv.Itoa(i),
			})
		}
	}
	if len(errs) > 0 {
		return errs, http.StatusBadRequest
	}
	return nil, http.StatusOK
}
