package middleware

import (
	"net/http"

	"catalog-service/internal/logger"
	"catalog-service/internal/constants"

	"github.com/gin-gonic/gin"
)

func PanicRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log := logger.NewContextLogger(c.Request.Context(), "PanicRecoveryMiddleware")
				log.Error("panic recovered: %v", nil)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"errors": []map[string]string{
						{
							"code":   constants.Error_GENERIC_SERVICE_ERROR,
							"entity": "internal",
							"cause":  "internal server error",
						},
					},
				})
			}
		}()
		c.Next()
	}
}
