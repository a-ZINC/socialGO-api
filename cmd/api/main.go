package main

import (
	"social-api/internal/env"
	"social-api/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	store := store.NewStorage(nil)
	app := &Application{
		Config: Config{
			Addr: ":" + env.GetString("ADDR", "3000"),
		},
		Store: store,
	}
	mux := app.mount()

	if err := app.run(mux); err != nil {
		panic(err)
	}
}
