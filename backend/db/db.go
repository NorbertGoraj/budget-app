package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://budget:budget@localhost:5432/budget?sslmode=disable"
	}

	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	return Pool.Ping(context.Background())
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
