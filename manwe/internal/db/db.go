package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open connects to the configured Postgres database. DATABASE_DSN, if
// set, overrides the DSN assembled from POSTGRES_USER/POSTGRES_PASSWORD/
// POSTGRES_DB (the vars docker-compose's db service is seeded from).
func Open() (*sqlx.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@db:5432/%s?sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
		)
	}

	sqlDB, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening postgres connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	return sqlDB, nil
}
