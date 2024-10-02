package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/biboyqg/late_days/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/biboyqg/late_days/docs"
)

type application struct {
	config config
	store *store.Storage
}

type config struct {
	addr string
	db dbConfig
	env string
	version string
}

type dbConfig struct {
	addr string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthcheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/late-days", func(r chi.Router) {
			r.Post("/", app.createLateDays)

			r.Route("/{studentID}", func(r chi.Router) {
				r.Get("/", app.getLateDaysByStudentID)
				r.Patch("/", app.updateLateDays)
			})
		})
	})

	return r
}

func (app *application) run(r http.Handler) error {
	docs.SwaggerInfo.Version = app.config.version
	docs.SwaggerInfo.Host = app.config.addr
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr: app.config.addr,
		Handler: r,
		WriteTimeout: 30 * time.Second,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: time.Minute,
	}

	log.Printf("Starting server on %s", app.config.addr)

	return srv.ListenAndServe()
}
