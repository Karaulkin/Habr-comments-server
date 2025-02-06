package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"

	"Habr-comments-server/internal/config"
	"Habr-comments-server/internal/graphql"
	"Habr-comments-server/internal/graphql/loaders"
	"Habr-comments-server/internal/service"
	"Habr-comments-server/internal/storage/pg"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.MustLoad()

	// Настраиваем логирование
	log := setupLogger(cfg.Env)
	log.Info("Starting server", slog.String("env", cfg.Env))

	// Подключаемся к БД
	db, err := pg.New(cfg.Storage)
	if err != nil {
		log.Error("Database connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	if db == nil {
		log.Error("Database connection is nil! Check config.")
		os.Exit(1)
	}

	defer func() {
		if err := db.Stop(context.Background()); err != nil {
			log.Error("Failed to close database connection", slog.Any("error", err))
		}
	}()

	log.Info("Connected to database")

	// Создаем сервисы
	svc := service.NewService(db, db)

	// Создаем DataLoader'ы
	lds := loaders.NewLoaders(svc)

	// Создаем резолвер GraphQL
	resolver := &graphql.Resolver{
		Service: svc,
		Loaders: lds,
	}

	// Запускаем GraphQL-сервер
	srv := handler.GraphQL(
		graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}),
		handler.ComplexityLimit(500), // Ограничение сложности запроса
	)

	srv = GraphQLLoggingMiddleware(log, srv)
	srv = AuthMiddleware(srv)

	mux := http.NewServeMux()

	// GraphQL Playground (UI для тестирования)
	mux.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))

	// Обработчик GraphQL API
	mux.Handle("/graphql", srv)

	// Запускаем HTTP-сервер
	serverAddr := cfg.HTTPServer.Address
	if serverAddr == "" {
		serverAddr = "localhost:8082"
	}

	log.Info("Server is running", slog.String("address", serverAddr))

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Канал для обработки завершения работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Ожидание сигнала завершения
	sig := <-quit
	log.Info("Shutting down server", slog.String("signal", sig.String()))

	// Создаем контекст с таймаутом для завершения
	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.Timeout*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Failed to gracefully shutdown server", slog.Any("error", err))
	} else {
		log.Info("Server shutdown successfully")
	}

	// Закрываем соединение с БД
	if err := db.Stop(ctx); err != nil {
		log.Error("Failed to close database connection", slog.Any("error", err))
	} else {
		log.Info("Database connection closed")
	}

	log.Info("Server exited")
}

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "userID", uint(1)) // Заглушка
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func GraphQLLoggingMiddleware(log *slog.Logger, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("Failed to read request body", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		log.Info("GraphQL Request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("body", string(body)),
			slog.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Info("GraphQL Request Processed",
			slog.String("duration", duration.String()),
			slog.String("remote_addr", r.RemoteAddr),
		)
	})
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
	default:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
