package main

import (
	"go_url_short/internal/config"
	"go_url_short/internal/storage/sqlite"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting loggin go_url_short", slog.String("env:", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}

	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		{
			log = slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		}
	case "dev":
		{
			log = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		}
	case "prod":
		{
			log = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			)
		}
	}

	return log
}
