package main

import (
	"log"

	"github.com/biboyqg/late_days/internal/db"
	"github.com/biboyqg/late_days/internal/env"
	"github.com/biboyqg/late_days/internal/store"
	_ "github.com/lib/pq"
)

//	@title			Late Days API
//	@description	API for managing late days
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/v1
func main() {
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5433/late_days?sslmode=disable"),
			maxOpenConns: 25,
			maxIdleConns: 25,
			maxIdleTime: "15m",
		},
		env: env.GetString("ENV", "dev"),
		version: env.GetString("VERSION", "0.0.1"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store: store,
	}

	mux := app.mount()

	err = app.run(mux)
	if err != nil {
		log.Fatal(err)
	}
}
