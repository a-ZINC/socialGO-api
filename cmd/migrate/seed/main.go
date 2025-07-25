package main

import (
	"log"
	"social-api/internal/db"
	"social-api/internal/env"
	"social-api/internal/store"
)

func main() {

	sql := db.New(env.GetString("DB_ADDR", "postgres://postgres:password@localhost:5432/social_media?sslmode=disable"), env.GetInt("DB_MAX_OPEN_CONNS", 25), env.GetInt("DB_MAX_IDLE_CONNS", 25), env.GetString("DB_CONN_MAX_LIFETIME", "5m"), env.GetString("DB_CONN_MAX_IDLE_TIME", "5m"))
	log.Printf("Connected to database at %s", env.GetString("DB_ADDR", "postgres://postgres:password@localhost:5432/social_api?sslmode=disable"))

	store := store.NewStorage(sql)

	err := db.Seed(store)
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
}
