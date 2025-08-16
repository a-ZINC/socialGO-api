package main

import (
	"net/http"
	"social-api/cmd/middlewares"
	"social-api/cmd/utils"
	"social-api/docs"
	"social-api/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Application struct {
	Config     Config
	Store      *store.Store
	Err        *utils.ErrorHandler
	Middleware *middlewares.Middleware
	Logger     *zap.SugaredLogger
}

type Config struct {
	Addr   string
	Db     DbConfig
	ApiUrl string
	Email  ConfigEmail
}

type ConfigEmail struct {
	ExpiryTime time.Duration
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
		docsUrl := "/v1/swagger/doc.json"
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.CreatePosthandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.Middleware.PostContext)
				r.Get("/", app.GetPostByIDHandler)
				r.Patch("/", app.UpdatePostHandler)
				r.Delete("/", app.DeletePostHandler)
			})
		})
		r.Route("/comments", func(r chi.Router) {
			r.Route("/post", func(r chi.Router) {
				r.Post("/{postId}", app.CreateCommentHandler)
			})
		})
		r.Route("/user", func(r chi.Router) {
			r.Put("/activate/{token}", app.ActivateUserHandler)
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.CreateUserContextMiddleware)
				r.Get("/", app.GetUserByIDHandler)
				r.Put("/follow", app.FollowUserHandler)
				r.Put("/unfollow", app.UnfollowUserHandler)
			})
			r.Get("/feed", app.GetUserFeedHandler)
		})

		// Public
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.RegisterUserHandler)
		})
	})
	return mux
}

func (app *Application) run(mux http.Handler) error {

	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Host = app.Config.ApiUrl
	docs.SwaggerInfo.Version = Version

	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Second,
	}

	app.Logger.Infof("Starting server on %s", app.Config.Addr)
	return server.ListenAndServe()
}
