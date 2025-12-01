package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeDatabase() *pgxpool.Pool {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		log.Fatal("PG_DSN environment variable is not set")
	}

	db, err := pgxpool.New(context.TODO(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
