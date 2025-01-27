package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/pingvincible/kvdatabase/internal/compute"
	"github.com/pingvincible/kvdatabase/internal/config"
)

var ErrServerIsNotListening = errors.New("server is not listening")

type Server struct {
	cfg      config.NetworkConfig
	listen   net.Listener
	computer *compute.Computer
}

func NewServer(cfg config.NetworkConfig, computer *compute.Computer) (*Server, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tcp addr: %s: %w", cfg.Address, err)
	}

	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to start listening on addr: %s: %w", cfg.Address, err)
	}

	return &Server{
		cfg:      cfg,
		computer: computer,
		listen:   listen,
	}, nil
}

func (s *Server) Addr() (string, error) {
	if s.listen == nil {
		return "", ErrServerIsNotListening
	}

	return s.listen.Addr().String(), nil
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
		slog.Info("test", result, err)
		if err != nil {
			slog.Error(
				"failed to process client query",
				slog.String("error", err.Error()),
			)

			// TODO sometimes nothing is sent to client
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
