package main

import (
	"log"
	"social-api/cmd/middlewares"
	"social-api/cmd/utils"
	_ "social-api/docs"
	"social-api/internal/db"
	"social-api/internal/env"
	"social-api/internal/store"

	"github.com/joho/godotenv"
)

//	@title			Social Media API
//	@description	This is a simple social media API built with Go.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@schemes	http

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	cfg := Config{
		Addr: ":" + env.GetString("ADDR", "3000"),
		Db: DbConfig{
			Addr:            env.GetString("DB_ADDR", "postgres://postgres:password@localhost:5432/social_api?sslmode=disable"),
			MaxOpenConns:    env.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    env.GetInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: env.GetString("DB_CONN_MAX_LIFETIME", "5m"),
			ConnMaxIdleTime: env.GetString("DB_CONN_MAX_IDLE_TIME", "5m"),
		},
		ApiUrl: env.GetString("API_URL", "localhost:3000"),
	}
	errorConfig := &utils.ErrorHandler{}
	middlewareCfg := &middlewares.Middleware{
		Err: errorConfig,
	}

	sql := db.New(cfg.Db.Addr, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleConns, cfg.Db.ConnMaxLifetime, cfg.Db.ConnMaxIdleTime)
	log.Printf("Connected to database at %s", cfg.Db.Addr)

	store := store.NewStorage(sql)

	app := &Application{
		Config:     cfg,
		Store:      store,
		Err:        errorConfig,
		Middleware: middlewareCfg,
	}
	mux := app.mount()

	if err := app.run(mux); err != nil {
		panic(err)
	}
}
