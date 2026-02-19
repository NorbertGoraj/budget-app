package infrastructure

import (
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
)

var DB *pg.DB

func Connect() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://budget:budget@localhost:5432/budget?sslmode=disable"
	}

	opts, err := pg.ParseURL(dsn)
	if err != nil {
		return fmt.Errorf("invalid DATABASE_URL: %w", err)
	}

	DB = pg.Connect(opts)

	if _, err := DB.Exec("SELECT 1"); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
