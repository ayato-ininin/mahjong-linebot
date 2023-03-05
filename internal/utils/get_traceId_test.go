package utils

import (
	"context"
	"testing"
)

func TestGetTraceID(t *testing.T) {
	// Create a context with traceId key and value
	ctx := context.WithValue(context.Background(), "traceId", "12345") //スタブ(traceIdがセットされたcontextが来ることはこの関数のテストの対象外)

	// Test the function with the created context
	traceId, err := GetTraceID(ctx)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if traceId != "12345" {
		t.Errorf("Unexpected traceId value: %s", traceId)
	}
}
