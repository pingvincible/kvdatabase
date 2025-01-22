package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/pingvincible/kvdatabase/internal/compute"
	"github.com/pingvincible/kvdatabase/internal/logger"
	"github.com/pingvincible/kvdatabase/internal/storage/engine"
)

func main() {
	logger.Configure(slog.LevelDebug)
	slog.Info("KV database started")

	kvEngine := engine.New()
	computer := compute.NewComputer(kvEngine)

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
