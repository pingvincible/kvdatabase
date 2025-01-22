package logger

import (
	"log/slog"
	"os"
)

func Configure(level slog.Level) {
	logLevel := &slog.LevelVar{}
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
	logLevel.Set(level)
}
