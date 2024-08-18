package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	Tracer = otel.Tracer("payment-service-provider")
)

type Tracing struct {
	Exporter   *zipkin.Exporter
	Provider   *trace.TracerProvider
	Resource   *resource.Resource
	Propagator propagation.TextMapPropagator
}

func New(ctx context.Context, exporterURL string) (*Tracing, error) {
	exporter, err := newTracerExporter(exporterURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracing exporter: %w", err)
	}
	resource, err := newResource(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracing resource: %w", err)
	}
	provider := newTracerProvider(exporter, resource)
	propagator := newPropagator()
	return &Tracing{
		Exporter:   exporter,
		Provider:   provider,
		Resource:   resource,
		Propagator: propagator,
	}, nil
}

func newTracerExporter(exporterURL string) (*zipkin.Exporter, error) {
	return zipkin.New(exporterURL)
}

func newTracerProvider(tracerExporter trace.SpanExporter, resource *resource.Resource) *trace.TracerProvider {
	return trace.NewTracerProvider(trace.WithBatcher(tracerExporter, trace.WithBatchTimeout(time.Second)), trace.WithResource(resource))
}

func newResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx, resource.WithAttributes(semconv.ServiceName("payment-service-provider")))
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
}
