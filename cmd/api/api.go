package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/biboyqg/late_days/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
