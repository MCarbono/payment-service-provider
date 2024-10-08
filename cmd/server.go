package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"payment-service-provider/application/usecase"
	"payment-service-provider/config"
	"payment-service-provider/infra/db"
	"payment-service-provider/infra/http/controllers"
	"payment-service-provider/infra/http/router"
	"payment-service-provider/infra/tracing"
	"payment-service-provider/infra/uow"

	sqlc "payment-service-provider/infra/db/sqlc"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
)

var (
	env = flag.String("env", "local", "used to know what environment the project is running")
)

func StartServer() {
	flag.Parse()
	cfg, err := config.LoadEnvConfig(*env)
	if err != nil {
		panic(err)
	}
	conn, err := db.Open(cfg.DatabaseConfig)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = fmt.Sprintf("file:///%s/infra/db/migration", dir)
	err = db.MigrateUp(conn, cfg.DatabaseConfig.Name, dir)
	if err != nil {
		panic(err)
	}
	uow := uow.NewUow(conn)
	querier := sqlc.New(conn)
	processTransaction := usecase.NewProcessTransaction(uow)
	listTransaction := usecase.NewListTransactions(querier)
	balance := usecase.NewClientBalance(querier)
	controller := controllers.NewTransaction(processTransaction, listTransaction, balance)
	r := router.New(controller)
	tracing, err := tracing.New(context.Background(), cfg.TracingExporterURL)
	if err != nil {
		panic(err)
	}
	defer tracing.Provider.Shutdown(context.Background())
	otel.SetTracerProvider(tracing.Provider)
	otel.SetTextMapPropagator(tracing.Propagator)
	fmt.Printf("Starting server on port %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.ServerPort), r))
}
