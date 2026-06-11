package main

import (
	"log/slog"
	"os"

	"github.com/iammikek/go-101/internal/app"
)

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	application, err := app.New(app.Config{
		DatabaseURL: envOrDefault("DATABASE_URL", "app.db"),
		APIKey:      envOrDefault("API_KEY", "dev-key-123"),
		LogLevel:    envOrDefault("LOG_LEVEL", "info"),
	})
	if err != nil {
		slog.Error("application bootstrap failed", "error", err)
		os.Exit(1)
	}

	slog.SetDefault(application.Logger)

	port := envOrDefault("PORT", "8000")
	application.Logger.Info("starting server", "port", port)
	if err := application.Router.Run(":" + port); err != nil {
		application.Logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
