package handler

import (
	"net/http"
	"strconv"

	"catalog-service/internal/api/validator"
	"catalog-service/internal/constants"
	"catalog-service/internal/dto"
	"catalog-service/internal/logger"
	"catalog-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage  = "1"
	defaultLimit = "10"
)

type ServiceHandler struct {
	usecase usecase.ServiceUsecase
}

func NewServiceHandler(usecase usecase.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{usecase: usecase}
}

func (h *ServiceHandler) Search(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.NewContextLogger(ctx, "ServiceHandler/Search")

	query := c.Query("q")
	pageStr := c.DefaultQuery("page", defaultPage)
	limitStr := c.DefaultQuery("limit", defaultLimit)

	log.Infof("Searching: query='%s', page='%s', limit='%s'", query, pageStr, limitStr)

	page, limit, errs, httpCode := validator.ValidateSearchRequest(pageStr, limitStr)
	if len(errs) > 0 {
		buildErrorResponse(c, httpCode, errs)
		return
	}

	services, total, err := h.usecase.Search(ctx, query, page, limit)
	if err != nil {
		log.Errorf(err, "failed to search services")
		buildErrorResponse(c, http.StatusInternalServerError, []dto.ErrorObj{
			{
				Code:   constants.Error_GENERIC_SERVICE_ERROR,
				Entity: "service",
				Cause:  "search failed",
			},
		})
		return
	}

	buildSuccessResponse(c, services, total, query, page, limit)
}

func buildSuccessResponse(c *gin.Context, services []*dto.ServiceDTO, total int, query string, page, limit int) {
	c.JSON(http.StatusOK, dto.ServiceListResponse{
		Data: &dto.ServiceListData{
			Count:    total,
			Services: services,
			Next:     buildNextURL(c, query, page, limit, total),
		},
	})
}

func buildErrorResponse(c *gin.Context, httpCode int, errs []dto.ErrorObj) {
	c.JSON(httpCode, dto.ServiceListResponse{
		Errors: errs,
	})
}

func buildNextURL(c *gin.Context, query string, page, limit, total int) *string {
	if (page * limit) >= total {
		return nil
	}
	q := c.Request.URL.Query()
	q.Set("page", strconv.Itoa(page+1))
	q.Set("limit", strconv.Itoa(limit))
	if query != "" {
		q.Set("q", query)
	}
	url := c.Request.URL.Path + "?" + q.Encode()
	return &url
}