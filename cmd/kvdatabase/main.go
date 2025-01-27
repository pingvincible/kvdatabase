package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/pingvincible/kvdatabase/internal/compute"
	"github.com/pingvincible/kvdatabase/internal/config"
	"github.com/pingvincible/kvdatabase/internal/logger"
	"github.com/pingvincible/kvdatabase/internal/storage/engine"
	"github.com/pingvincible/kvdatabase/internal/tcp"
)

func main() {
	cfg, err := config.Load("./kvdatabase.yaml")
	if err != nil {
		log.Fatal(err)
	}

	logger.Configure(cfg.Logging.Level)
	slog.Info("KV database started")

	kvEngine := engine.New()
	computer := compute.NewComputer(kvEngine)

	server, err := tcp.NewServer(cfg.Network.Address, computer)
	if err != nil {
		slog.Error(
			"failed to start server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	server.Start()
}
