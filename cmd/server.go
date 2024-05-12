package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"payment-service-provider/config"
	"payment-service-provider/infra/db"

	"github.com/golang-migrate/migrate/v4"
	migratePg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func StartServer() {
	cfg, err := config.LoadEnvConfig()
	if err != nil {
		panic(err)
	}
	db, err := db.Open(cfg.DatabaseConfig)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = fmt.Sprintf("file:///%s/infra/db/migration", dir)
	err = MigrateUp(db, cfg.DatabaseConfig.Name, dir)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting server on port 3000")
	panic(http.ListenAndServe(":3000", nil))
}

func MigrateUp(conn *sql.DB, dbName string, migrationsPath string) error {
	driver, err := migratePg.WithInstance(conn, &migratePg.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %w", err)
	}
	migrateClient, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate client: %w", err)
	}
	err = migrateClient.Up()
	if err != nil {
		return fmt.Errorf("failed to perform migration: %w", err)
	}
	fmt.Println("Migration completed")
	return nil
}
