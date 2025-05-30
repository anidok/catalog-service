package appcontext

import (
	"context"
)

type key string

const (
	CorrelationIDKey = key("CorrelationID")
)

func LogFields(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})
	fields[string(CorrelationIDKey)] = Value(ctx, CorrelationIDKey)
	return fields
}

func Value(ctx context.Context, key key) string {
	if val, ok := ctx.Value(key).(string); ok {
		return val
	}
	return ""
}
