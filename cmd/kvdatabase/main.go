package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pingvincible/kvdatabase/internal/compute"
	"github.com/pingvincible/kvdatabase/internal/config"
	"github.com/pingvincible/kvdatabase/internal/logger"
	"github.com/pingvincible/kvdatabase/internal/storage/engine"
	"github.com/pingvincible/kvdatabase/internal/tcp"
)

func main() {
	cfg, err := config.Load("./config/kvdatabase.yaml")
	if err != nil {
		log.Fatal(err)
	}

	flags := handleFlags(cfg)
	cfg.UpdateWithFlags(flags)
	log.Printf("%+v", cfg)

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

	addr, err := server.Addr()
	if err != nil {
		slog.Error(
			"failed to get server address",
			slog.String("error", err.Error()),
		)
	}

	slog.Info("tcp server started", slog.String("address", addr))

	server.Start()
}

func handleFlags(cfg *config.Config) config.Flags {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagSet.Usage = cleanenv.FUsage(flagSet.Output(), cfg, nil, flagSet.Usage)

	configFlags := config.Flags{
		EngineType:     flagSet.String("engineType", cfg.Engine.Type, "database engine type"),
		Address:        flagSet.String("address", cfg.Network.Address, "address to listen"),
		MaxConnections: flagSet.Int("maxConnections", cfg.Network.MaxConnections, "max client connections"),
		MaxMessageSize: flagSet.String("maxMessageSize", cfg.Network.MaxMessageSize, "max message size"),
		IdleTimeout:    flagSet.Duration("idleTimeout", cfg.Network.IdleTimeout, "idle timeout"),
		LoggingLevel:   flagSet.String("logLevel", cfg.Logging.Level, "log level"),
		LoggingOutput:  flagSet.String("logOutput", cfg.Logging.Output, "log output filename"),
	}

	_ = flagSet.Parse(os.Args[1:])

	return configFlags
}
