package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/pingvincible/kvdatabase/internal/compute/parser"
	"github.com/pingvincible/kvdatabase/internal/logger"
)

func main() {
	// Create parser.
	// Create engine.
	// Create interactive shell.
	logger.Configure(slog.LevelDebug)
	slog.Info("Database started")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("> ")
		scanner.Scan()
		text := scanner.Text()

		command, err := parser.Parse(text)
		if err != nil {
			fmt.Printf("%s\n", err)

			continue
		}

		command.Execute()
	}
}
