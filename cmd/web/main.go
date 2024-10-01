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
	"github.com/berberapan/my-stuff/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	logger         *slog.Logger
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dbAddress := os.Getenv("DB_DETAILS")
	if dbAddress == "" {
		logger.Error("Env for DB address not set")
		os.Exit(1)
	}

	dsnString := fmt.Sprintf("postgres://%s?sslmode=require", os.Getenv("DB_DETAILS"))
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

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = false // TODO change to true when TSL sorted

	app := &application{
		logger:         logger,
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
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
