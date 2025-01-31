package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/pingvincible/kvdatabase/internal/kvio"
	"github.com/pingvincible/kvdatabase/internal/logger"
	"github.com/pingvincible/kvdatabase/internal/tcp"
)

func main() {
	kvLogger := logger.Configure("debug")
	kvLogger.Info("KV CLI started")

	var addr string

	flag.StringVar(&addr, "hostname", "localhost:3223", "address to connect to")
	flag.Parse()

	client, err := tcp.NewClient(addr)
	if err != nil {
		kvLogger.Error(
			"failed to create tcp client",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	consoleReadWriter := kvio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))

	err = Run(consoleReadWriter, client.ReadWriter)
	if err != nil {
		kvLogger.Error(
			"failed to run",
			slog.String("error", err.Error()),
		)
	}

	err = client.Close()
	if err != nil {
		kvLogger.Error(
			"failed to close tcp client",
			slog.String("error", err.Error()),
		)
	}
}

func Run(consoleReadWriter *kvio.ReadWriter, clientReadWriter *kvio.ReadWriter) error {
	for {
		err := consoleReadWriter.Write(">>")
		if err != nil {
			return fmt.Errorf("failed to write to console: %w", err)
		}

		request, err := consoleReadWriter.ReadLine()
		if err != nil {
			return fmt.Errorf("failed to read from console: %w", err)
		}

		if request == "\n" {
			continue
		}

		if request == "Q\n" {
			return nil
		}

		err = clientReadWriter.WriteLine(request)
		if err != nil {
			return fmt.Errorf("failed to write to server: %w", err)
		}

		response, err := clientReadWriter.ReadLine()
		if err != nil {
			return fmt.Errorf("failed to read from server: %w", err)
		}

		err = consoleReadWriter.WriteLine("->: " + response)
		if err != nil {
			return fmt.Errorf("failed to write to console: %w", err)
		}
	}
}
