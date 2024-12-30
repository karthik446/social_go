package main

import (
	"database/sql"
	"log"

	"github.com/karthik446/social/internal/db"
	"github.com/karthik446/social/internal/env"
	"github.com/karthik446/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgresql://admin:password123@localhost:5432/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v", err)
		}
	}(conn)

	postgresStorage := store.NewPostgresStorage(conn)
	db.Seed(postgresStorage)
}
