package router

import (
	"fmt"
	"net/http"
	"payment-service-provider/infra/http/controllers"
	"payment-service-provider/infra/tracing"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func New(controller *controllers.Transaction) http.Handler {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(JSONContetTyeResponse)
	r.Use(otelchi.Middleware("payment-service-provider", otelchi.WithChiRoutes(r)))
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracing.Tracer.Start(r.Context(), "ping init")
		defer span.End()
		span.AddEvent("middle of the ping", trace.WithAttributes(attribute.String("middle", "123456"), attribute.Bool("teste-bool", true)))
		span.SetAttributes(attribute.String("test-123", "123"))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	})
	r.Post("/transactions", controller.ProcessTransaction)
	r.Get("/transactions/{client_id}", controller.ListTransactions)
	r.Get("/balance/{client_id}", controller.Balance)
	return r
}

func JSONContetTyeResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
