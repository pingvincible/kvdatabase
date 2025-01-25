package main

import (
	"bufio"
	"fmt"
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

	server, err := tcp.NewServer(":1234")
	if err != nil {
		slog.Error(
			"failed to start server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	kvEngine := engine.New()
	computer := compute.NewComputer(kvEngine)

	server.Start()

	Run(computer)
}

func Run(computer *compute.Computer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("> ")
		scanner.Scan()

		text := scanner.Text()

		result, err := computer.Process(text)
		if err != nil {
			fmt.Println(err)

			continue
		}

		if len(result) > 0 {
			fmt.Println(result)
		}
	}
}
