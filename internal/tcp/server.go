package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/pingvincible/kvdatabase/internal/compute"
	"github.com/pingvincible/kvdatabase/internal/config"
)

var ErrServerIsNotListening = errors.New("server is not listening")

type Server struct {
	cfg              config.NetworkConfig
	listen           net.Listener
	computer         *compute.Computer
	ClientsHandled   int
	ClientsDiscarded int
	logger           *slog.Logger

	isListening bool

	mutex   sync.Mutex
	clients int
}

func NewServer(
	cfg config.NetworkConfig,
	computer *compute.Computer,
	logger *slog.Logger,
) (*Server, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tcp addr: %s: %w", cfg.Address, err)
	}

	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to start listening on addr: %s: %w", cfg.Address, err)
	}

	server := &Server{
		cfg:      cfg,
		logger:   logger,
		computer: computer,
		listen:   listen,
		clients:  0,
	}

	server.mutex.Lock()
	server.isListening = true
	server.mutex.Unlock()

	return server, nil
}

func (s *Server) Addr() (string, error) {
	if s.listen == nil {
		return "", ErrServerIsNotListening
	}

	return s.listen.Addr().String(), nil
}

func (s *Server) Start() {
	for {
		s.mutex.Lock()
		if !s.isListening {
			s.mutex.Unlock()

			break
		}
		s.mutex.Unlock()

		clients := s.getClients()

		s.logger.Info(
			"waiting for client",
			slog.Int("clients connected", clients),
			slog.Int("max connections", s.cfg.MaxConnections),
		)

		conn, err := s.listen.Accept()
		if err != nil {
			s.logger.Error(
				"failed to accept connection",
				slog.String("error", err.Error()),
			)

			const failSleep = 10

			time.Sleep(failSleep * time.Millisecond)

			continue
		}

		if ok := s.increaseClients(); ok {
			go s.handleClient(conn)

			s.ClientsHandled++
		} else {
			s.ClientsDiscarded++

			s.logger.Info("failed to handle client, too many connections")

			err = conn.Close()
			if err != nil {
				s.logger.Error(
					"failed to close client connection",
					slog.String("error", err.Error()),
				)
			}
		}
	}
}

func (s *Server) getClients() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.clients
}

func (s *Server) handleClient(conn net.Conn) {
	defer func() {
		s.mutex.Lock()
		if s.clients > 0 {
			s.clients--
		}
		s.mutex.Unlock()

		err := conn.Close()
		if err != nil {
			s.logger.Error(
				"failed to close client connection",
				slog.String("error", err.Error()),
			)
		}
	}()

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			s.logger.Error(
				"failed to read data from tcp client",
				slog.String("error", err.Error()),
			)

			return
		}

		s.logger.Info(
			"received data from client",
			slog.String("data", netData),
		)

		response, err := s.computer.Process(netData)
		if err != nil {
			s.logger.Error(
				"failed to process client query",
				slog.String("error", err.Error()),
			)

			response = fmt.Sprintf("error: %s\n", err)
		}

		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			s.logger.Error(
				"failed to send data to client",
				slog.String("error", err.Error()),
			)

			return
		}
	}
}

func (s *Server) GetClients() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.clients
}

func (s *Server) increaseClients() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.clients < s.cfg.MaxConnections {
		s.clients++

		return true
	}

	return false
}

func (s *Server) Stop() error {
	defer func() {
		s.mutex.Lock()
		s.isListening = false
		s.mutex.Unlock()
	}()

	err := s.listen.Close()
	if err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}
