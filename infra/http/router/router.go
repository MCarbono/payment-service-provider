package router

import (
	"net/http"
	"payment-service-provider/infra/http/controllers"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func New(controller *controllers.Transaction) http.Handler {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(JSONContetTyeResponse)
	r.Post("/transactions", controller.ProcessTransaction)
	r.Get("/transactions/{id}", controller.ListTransactions)
	r.Get("/balance/{id}", controller.Balance)
	return r
}

func JSONContetTyeResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
