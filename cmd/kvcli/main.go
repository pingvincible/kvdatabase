package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/pingvincible/kvdatabase/internal/logger"
)

func main() {
	logger.Configure("debug")
	slog.Info("KV CLI started")

	var addr string

	flag.StringVar(&addr, "hostname", "localhost:3223", "address to connect to")
	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		slog.Error(
			"failed to resolve tcp address",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		slog.Error(
			"failed to connect to tcp address",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			slog.Error(
				"failed to close tcp connection",
				slog.String("error", err.Error()),
			)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")

		text, err := reader.ReadString('\n')
		if err != nil {
			slog.Error(
				"failed to read line",
				slog.String("error", err.Error()),
			)

			break
		}

		_, err = fmt.Fprintf(conn, "%s\n", text)
		if err != nil {
			slog.Error(
				"failed to send line to tcp server",
				slog.String("error", err.Error()),
			)

			break
		}

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			slog.Error(
				"failed to receive response from tcp server",
				slog.String("error", err.Error()),
			)

			break
		}

		fmt.Print("->: " + message)
	}
}
