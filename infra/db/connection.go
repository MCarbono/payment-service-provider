package db

import (
	"database/sql"
	"fmt"
	"payment-service-provider/config"
	"time"
)

const maxRetries = 5
const retryInterval = 2000 * time.Millisecond

func Open(cfg config.DatabaseConfig) (DB *sql.DB, err error) {
	DB, err = sql.Open("postgres", fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name))
	if err != nil {
		return
	}
	for i := 0; i < maxRetries; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		if err != nil {
			fmt.Printf("Connection failed (Attempt %d): %v\n", i+1, err)
			time.Sleep(retryInterval)
		}
	}
	fmt.Println("Connected to the database!")
	return
}
