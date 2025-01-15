package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/pingvincible/kvdatabase/internal/compute/parser"
	"github.com/pingvincible/kvdatabase/internal/logger"
	"github.com/pingvincible/kvdatabase/internal/storage/engine"
)

func main() {
	logger.Configure(slog.LevelDebug)
	slog.Info("Database started")

	scanner := bufio.NewScanner(os.Stdin)
	kvEngine := engine.New()

	for {
		fmt.Printf("> ")
		scanner.Scan()
		text := scanner.Text()

		command, err := parser.Parse(text)
		if err != nil {
			fmt.Printf("%s\n", err)

			continue
		}

		Execute(kvEngine, command)
	}
}

func Execute(kvEngine *engine.Engine, command parser.Command) {
	switch command.Type {
	case parser.CommandSet:
		kvEngine.Set(command.Key, command.Value)
	case parser.CommandGet:
		fmt.Println(kvEngine.Get(command.Key))
	case parser.CommandDel:
		kvEngine.Delete(command.Key)
	}
}
