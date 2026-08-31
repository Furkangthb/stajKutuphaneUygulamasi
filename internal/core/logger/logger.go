package logger

import (
	"io"
	"log/slog"
	"os"
)

func New() *slog.Logger {
	var handler slog.Handler

	if err := os.MkdirAll("/app/logs", 0755); err != nil {
		return slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logFile, err := os.OpenFile("/app/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	var writer io.Writer = os.Stdout
	if err == nil {
		writer = io.MultiWriter(os.Stdout, logFile)
	}

	if os.Getenv("APP_ENV") == "production" {
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	return slog.New(handler)

}
