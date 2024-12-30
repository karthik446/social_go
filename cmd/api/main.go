package main

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	"github.com/karthik446/social/internal/db"
	"github.com/karthik446/social/internal/env"
	"github.com/karthik446/social/internal/store"
)

func main() {
	// Your code here
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cfg := config{
		addr: env.GetString("SERVER_ADDR", ":8081"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://admin:password123@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 20),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:     env.GetString("ENV", "development"),
		version: env.GetString("VERSION", "1.0.0"),
	}

	database, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Fatalf("error closing db: %v", err)
		}
	}(database)
	log.Println("Connected to db")
	postgresStorage := store.NewPostgresStorage(database)

	log.Println("Starting server on", cfg.addr)
	app := &application{
		config: cfg,
		store:  postgresStorage,
	}
	mux := app.mount()

	log.Fatal(app.run(mux))
}
