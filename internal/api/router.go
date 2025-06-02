package api

import (
	"strings"

	"github.com/gin-gonic/gin"

	"catalog-service/internal/api/handler"
	"catalog-service/internal/config"
	"catalog-service/internal/middleware"
	"catalog-service/internal/repository"
	"catalog-service/internal/usecase"
)

var allowedEnvs = []string{"dev", "test", "uat", "production"}

func NewRouter(repo repository.ServiceRepository) *gin.Engine {
	env := config.AppEnv()
	if !isAllowedEnv(env) {
		panic("invalid APP_ENV: must be one of dev, test, uat, production")
	}
	if !strings.EqualFold(env, "dev") {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CorrelationIDMiddleware())

	serviceUsecase := usecase.NewServiceUsecase(repo)
	serviceHandler := handler.NewServiceHandler(serviceUsecase)

	api := r.Group("/api")
	{
		api.GET("/services", serviceHandler.Search)
		api.GET("/services/:id", serviceHandler.GetByID)
	}

	return r
}

func isAllowedEnv(env string) bool {
	env = strings.ToLower(env)
	for _, e := range allowedEnvs {
		if env == e {
			return true
		}
	}
	return false
}
