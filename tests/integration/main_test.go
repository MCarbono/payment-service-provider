package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"payment-service-provider/config"
	"payment-service-provider/infra/db"
	"payment-service-provider/infra/uow"
	"payment-service-provider/tests/postgrescontainer"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	migratePg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	DbConn *sql.DB
	Uow    *uow.UowImpl
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := postgrescontainer.StartPostgresContainer(ctx)
	if err != nil {
		panic(err)
	}
	cfg := config.DatabaseConfig{
		Host:     container.Host,
		Port:     container.Port,
		User:     container.User,
		Password: container.Password,
		Name:     container.Name,
	}
	conn, err := db.Open(cfg)
	if err != nil {
		panic(err)
	}
	DbConn = conn
	Uow = uow.NewUow(DbConn)
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir = strings.Replace(dir, "/tests/integration", "", 1)
	dir = fmt.Sprintf("file:///%s/infra/db/migration", dir)
	err = migrateUp(conn, cfg.Name, dir)
	if err != nil {
		panic(err)
	}
	exitCode := m.Run()
	conn.Close()
	if err := container.Terminate(ctx); err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

func migrateUp(conn *sql.DB, dbName string, migrationsPath string) error {
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

func cleanDatabaseTables(conn *sql.DB) {
	_, err := conn.Exec("DELETE FROM payables;")
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec("DELETE from transactions;")
	if err != nil {
		panic(err)
	}
}
