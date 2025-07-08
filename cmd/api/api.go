package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	Config Config
}

type Config struct {
	Addr string
}

func (app *Application) mount() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthHandler)
	})
	return mux
}

func (app *Application) run(mux http.Handler) error {

	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Second,
	}

	log.Printf("Starting server on %s", app.Config.Addr)
	return server.ListenAndServe()
}
