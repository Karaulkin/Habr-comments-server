package main

import (
	"Habr-comments-server/internal/config"
	"Habr-comments-server/internal/storage/pg"
	"fmt"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg) // TODO: delete in prod

	log := setupLogger(cfg.Env)

	log.Info("starting server", slog.String("env", cfg.Env))
	// log.Debug("debug messages are enabled")

	// подключение к бд работает
	// Удалить в будущем был как тест
	_, err := pg.New(cfg.Storage)

	if err != nil {
		log.Error("database", err)
		return
	}

	log.Info("starting database")

	// TODO: init storage: pgx

	// TODO: init router: chi, "chi render" пока не то

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
