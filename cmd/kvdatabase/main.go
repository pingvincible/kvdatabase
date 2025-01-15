package main

import (
	"log/slog"

	"github.com/pingvincible/kvdatabase/internal/logger"
)

func main() {
	logger.Configure(slog.LevelDebug)
	slog.Info("Database started")
	// Configure logger.
	// Create parser.
	// Create engine.
	// Create interactive shell.
}
