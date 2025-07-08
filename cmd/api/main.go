package main

import (
	"social-api/internal/env"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	app := &Application{
		Config: Config{
			Addr: ":" + env.GetString("ADDR", "3000"),
		},
	}
	mux := app.mount()

	if err := app.run(mux); err != nil {
		panic(err)
	}
}
