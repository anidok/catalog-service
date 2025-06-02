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

	page, pageErr, pageCode := ValidatePageWithError(pageStr)
	if pageErr != nil {
		errors = append(errors, *pageErr)
		httpCode = pageCode
	}
	limit, limitErr, limitCode := ValidateLimitWithError(limitStr)
	if limitErr != nil {
		errors = append(errors, *limitErr)
		if httpCode == http.StatusOK {
			httpCode = limitCode
		}
	}
	return page, limit, errors, httpCode
}

func ValidatePageWithError(pageStr string) (int, *dto.ErrorObj, int) {
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

func ValidateLimitWithError(limitStr string) (int, *dto.ErrorObj, int) {
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
