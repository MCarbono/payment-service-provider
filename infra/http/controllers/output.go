package controllers

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type ControllerOutput struct {
	TraceID string `json:"trace_id"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
}

func newControllerOutput(ctx context.Context, data any, err error) ControllerOutput {
	output := ControllerOutput{TraceID: trace.SpanFromContext(ctx).SpanContext().TraceID().String()}
	if err != nil {
		output.Error = err.Error()
	}
	output.Data = struct{}{}
	if data != nil {
		output.Data = data
	}
	return output
}
