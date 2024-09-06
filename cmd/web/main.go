package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	dsnString := fmt.Sprintf("postgres://%s?sslmode=disable", os.Getenv("DB_DETAILS"))
	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", dsnString, "PostgreSQL data source name")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = false // TODO change to true when TSL sorted

	app := &application{
		logger:         logger,
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	logger.Info("starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping(context.Background())
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
