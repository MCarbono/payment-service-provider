package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratePg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

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
		if err == migrate.ErrNoChange {
			fmt.Println("no need to migrate. Database at last version already")
			return nil
		}
		return fmt.Errorf("failed to perform migration: %w", err)
	}
	fmt.Println("Migration completed")
	return nil
}
