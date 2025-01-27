package tcp

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"

	"github.com/pingvincible/kvdatabase/internal/compute"
)

type Server struct {
	addr     string
	listen   net.Listener
	computer *compute.Computer
}

func NewServer(addr string, computer *compute.Computer) (*Server, error) {
	server := Server{
		addr:     addr,
		computer: computer,
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tcp addr: %s: %w", addr, err)
	}

	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to start listening on addr: %s: %w", addr, err)
	}

	server.listen = listen

	return &server, nil
}

func (s *Server) Start() {
	for {
		conn, err := s.listen.Accept()
		if err != nil {
			slog.Error(
				"failed to accept connection",
				slog.String("error", err.Error()),
			)

			continue
		}

		go s.handleClient(conn)
	}
}

func (s *Server) Stop() error {
	// TODO stop accepting connections
	err := s.listen.Close()
	if err != nil {
		return fmt.Errorf("failed to close tcp server: %w", err)
	}

	return nil
}

func (s *Server) handleClient(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			slog.Error(
				"failed to close client connection",
				slog.String("error", err.Error()),
			)
		}
	}()

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			slog.Error(
				"failed to read data from tcp client",
				slog.String("error", err.Error()),
			)

			return
		}

		slog.Info(
			"received data from client",
			slog.String("data", netData),
		)

		result, err := s.computer.Process(netData)
		if err != nil {
			slog.Error(
				"failed to process client query",
				slog.String("error", err.Error()),
			)

			_, err = conn.Write([]byte(fmt.Sprintf("error: %s\n", err)))
			if err != nil {
				slog.Error(
					"failed to send data to client",
					slog.String("error", err.Error()),
				)

				return
			}
		}

		_, err = conn.Write([]byte(result + "\n"))
		if err != nil {
			slog.Error(
				"failed to send data to client",
				slog.String("error", err.Error()),
			)

			return
		}
	}
}
