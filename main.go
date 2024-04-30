package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/zaidsasa/xbankapi/internal/api"
	"github.com/zaidsasa/xbankapi/internal/http"
	"github.com/zaidsasa/xbankapi/internal/storage"
	"github.com/zaidsasa/xbankapi/internal/validator"
)

var errMissingEnviromentVariableDatabaseURL = errors.New("missing environment variable DATABASE_URL")

const (
	defualtServiceAddr = ":3000"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal(errMissingEnviromentVariableDatabaseURL)
	}

	addr := os.Getenv("SERVCE_ADDRESS")
	if addr == "" {
		addr = defualtServiceAddr

		slog.Info("using default serivce address", "address", addr)
	}

	validator.ConfigureDefaultValidator()

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	logger := slog.Default()

	storage := storage.New(pool)

	accountService := api.NewAccountService(pool, storage, logger)

	srv := http.NewServer(
		logger,
		api.NewAccountHandler(accountService),
		api.NewPropsHandler(pool),
	)

	if err := srv.Start(addr); err != nil {
		panic(err)
	}
}
