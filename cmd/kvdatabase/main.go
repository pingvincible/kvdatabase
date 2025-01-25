package main

import (
	"log/slog"
	"os"

	"github.com/pingvincible/kvdatabase/internal/compute"
	"github.com/pingvincible/kvdatabase/internal/logger"
	"github.com/pingvincible/kvdatabase/internal/storage/engine"
	"github.com/pingvincible/kvdatabase/internal/tcp"
)

func main() {
	logger.Configure(slog.LevelDebug)
	slog.Info("KV database started")

	kvEngine := engine.New()
	computer := compute.NewComputer(kvEngine)

	server, err := tcp.NewServer(":1234", computer)
	if err != nil {
		slog.Error(
			"failed to start server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	server.Start()
}
