package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"payment-service-provider/application/usecase"
	"payment-service-provider/config"
	"payment-service-provider/infra/db"
	"payment-service-provider/infra/http/controllers"
	"payment-service-provider/infra/http/router"
	"payment-service-provider/infra/uow"

	sqlc "payment-service-provider/infra/db/sqlc"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func StartServer() {
	cfg, err := config.LoadEnvConfig()
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
	fmt.Printf("Starting server on port %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.ServerPort), r))
}
