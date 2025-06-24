package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/aybangueco/golang-be-template/internal/database"
	"github.com/aybangueco/golang-be-template/internal/env"
	"github.com/aybangueco/golang-be-template/internal/smtp"
	"github.com/aybangueco/golang-be-template/internal/version"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL     string
	port        int
	tokenSecret string
	db          struct {
		dsn      string
		username string
		password string
		host     string
		port     int
		name     string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		from     string
	}
}

type application struct {
	config config
	db     *database.Queries
	dbpool *pgxpool.Pool
	logger *slog.Logger
	mailer *smtp.Mailer
	wg     sync.WaitGroup
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:3000")
	cfg.port = env.GetInt("PORT", 4000)
	cfg.tokenSecret = env.GetString("TOKEN_SECRET", "secret")

	cfg.db.username = env.GetString("DB_USERNAME", "postgres")
	cfg.db.password = env.GetString("DB_PASSWORD", "postgres")
	cfg.db.host = env.GetString("DB_HOST", "localhost")
	cfg.db.port = env.GetInt("DB_PORT", 5432)
	cfg.db.name = env.GetString("DB_NAME", "default")

	cfg.smtp.username = env.GetString("SMTP_USERNAME", "example_username")
	cfg.smtp.password = env.GetString("SMTP_PASSWORD", "pa55word")
	cfg.smtp.host = env.GetString("SMTP_HOST", "example.smtp.host")
	cfg.smtp.port = env.GetInt("SMTP_PORT", 25)
	cfg.smtp.from = env.GetString("SMTP_FROM", "Example Name <no_reply@example.org>")

	showVersion := flag.Bool("version", false, "Display version and exit")
	flag.StringVar(&cfg.db.dsn, "database-dsn", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.db.username, cfg.db.password, cfg.db.host, cfg.db.port, cfg.db.name), "Database dsn of server")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		return err
	}
	defer dbpool.Close()

	mailer, err := smtp.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.from)
	if err != nil {
		return err
	}

	app := &application{
		config: cfg,
		db:     database.New(dbpool),
		dbpool: dbpool,
		logger: logger,
		mailer: mailer,
	}

	return app.serveHTTP()
}
