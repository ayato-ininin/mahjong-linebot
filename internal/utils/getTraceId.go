package utils

import (
	"context"
	"errors"
)

// contextからtraceIdを取り出す
func GetTraceID(ctx context.Context) (string, error) {
	traceId, ok := ctx.Value("traceId").(string)
	if !ok {
		return "", errors.New("traceId not found in context")
	}
	return traceId, nil
}
