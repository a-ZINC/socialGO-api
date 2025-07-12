package main

import (
	"log"
	"net/http"
	"social-api/cmd/middlewares"
	"social-api/cmd/utils"
	"social-api/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	Config     Config
	Store      *store.Store
	Err        *utils.ErrorHandler
	Middleware *middlewares.Middleware
}

type Config struct {
	Addr string
	Db   DbConfig
}

var Version string = "1.0.0"

type DbConfig struct {
	Addr            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime string
	ConnMaxIdleTime string
}

func (app *Application) mount() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.CreatePosthandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.Middleware.PostContext)
				r.Get("/", app.GetPostByIDHandler)
				r.Delete("/", app.DeletePostHandler)
			})
		})
		r.Route("/comments", func(r chi.Router) {
			r.Route("/post", func(r chi.Router) {
				r.Post("/{postId}", app.CreateCommentHandler)
			})
		})
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
