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
		buildErrorListResponse(c, httpCode, errs)
		return
	}

	services, total, err := h.usecase.Search(ctx, query, page, limit)
	if err != nil {
		log.Errorf(err, "failed to search services")
		buildErrorListResponse(c, http.StatusInternalServerError, []dto.ErrorObj{
			{
				Code:   constants.Error_GENERIC_SERVICE_ERROR,
				Entity: "service",
				Cause:  "search failed",
			},
		})
		return
	}

	buildSuccessListResponse(c, services, total, query, page, limit)
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	log := logger.NewContextLogger(ctx, "ServiceHandler/GetByID")
	log.Infof("fetching service by id='%s'", id)

	if errs, httpCode := validator.ValidateID(id); len(errs) > 0 {
		c.JSON(httpCode, dto.ServiceDetailResponse{
			Success: false,
			Errors:  errs,
		})
		return
	}

	service, err := h.usecase.FindByID(ctx, id)
	if err != nil {
		buildErrorDetailResponse(c, http.StatusNotFound, []dto.ErrorObj{
			{
				Code:   constants.Error_SERVICE_NOT_FOUND,
				Entity: "service",
				Cause:  "service not found",
			},
		})
		return
	}

	buildSuccessDetailResponse(c, service)
}

func (h *ServiceHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.NewContextLogger(ctx, "ServiceHandler/Create")

	var req dto.ServiceDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf(err, "invalid request body")
		c.JSON(http.StatusBadRequest, dto.ServiceDetailResponse{
			Success: false,
			Errors: []dto.ErrorObj{{
				Code:   constants.Error_MALFORMED_DATA,
				Entity: "service",
				Cause:  "invalid request body",
			}},
		})
		return
	}

	if errs, httpCode := validator.ValidateCreateRequest(&req); len(errs) > 0 {
		c.JSON(httpCode, dto.ServiceDetailResponse{
			Success: false,
			Errors:  errs,
		})
		return
	}

	service, err := h.usecase.Create(ctx, &req)
	if err != nil {
		log.Errorf(err, "failed to create service")
		c.JSON(http.StatusInternalServerError, dto.ServiceDetailResponse{
			Success: false,
			Errors: []dto.ErrorObj{{
				Code:   "900",
				Entity: "service",
				Cause:  "failed to create service",
			}},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ServiceDetailResponse{
		Success: true,
		Data:    service,
	})
}

func buildSuccessListResponse(c *gin.Context, services []*dto.ServiceDTO, total int, query string, page, limit int) {
	c.JSON(http.StatusOK, dto.ServiceListResponse{
		Success: true,
		Data: &dto.ServiceListData{
			Count:    total,
			Services: services,
			Next:     buildNextURL(c, query, page, limit, total),
		},
	})
}

func buildSuccessDetailResponse(c *gin.Context, service *dto.ServiceDTO) {
	c.JSON(http.StatusOK, dto.ServiceDetailResponse{
		Success: true,
		Data:    service,
	})
}

func buildErrorListResponse(c *gin.Context, httpCode int, errs []dto.ErrorObj) {
	c.JSON(httpCode, dto.ServiceListResponse{
		Errors: errs,
	})
}

func buildErrorDetailResponse(c *gin.Context, httpCode int, errs []dto.ErrorObj) {
	c.JSON(httpCode, dto.ServiceDetailResponse{
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
