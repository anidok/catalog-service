package middleware

import (
	"catalog-service/internal/appcontext"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const CorrelationIDKeyHeader = "X-Correlation-ID"

func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader(CorrelationIDKeyHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}
		c.Set(string(appcontext.CorrelationIDKey), correlationID)
		ctx := context.WithValue(c.Request.Context(), appcontext.CorrelationIDKey, correlationID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(CorrelationIDKeyHeader, correlationID)
		c.Next()
	}
}
