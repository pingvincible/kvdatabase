package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
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

	wgClients sync.WaitGroup
	clients   atomic.Int32
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
		cfg:       cfg,
		logger:    logger,
		computer:  computer,
		listen:    listen,
		wgClients: sync.WaitGroup{},
	}

	return server, nil
}

func (s *Server) Addr() (string, error) {
	if s.listen == nil {
		return "", ErrServerIsNotListening
	}

	return s.listen.Addr().String(), nil
}

func (s *Server) Run() {
	for {
		s.logger.Info(
			"waiting for client",
			slog.Int("clients connected", int(s.clients.Load())),
			slog.Int("max connections", s.cfg.MaxConnections),
		)

		conn, err := s.listen.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.logger.Error(
					"listener was closed",
					slog.String("error", err.Error()),
				)

				return
			}

			s.logger.Error(
				"failed to accept connection",
				slog.String("error", err.Error()),
			)

			const failSleep = 10

			time.Sleep(failSleep * time.Millisecond)

			continue
		}

		s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	if ok := incrementWithLimit(&s.clients, int32(s.cfg.MaxConnections)); ok {
		s.wgClients.Add(1)

		go s.handleClient(conn)

		s.ClientsHandled++
	} else {
		s.ClientsDiscarded++

		s.logger.Info("failed to handle client, too many connections")

		err := conn.Close()
		if err != nil {
			s.logger.Error(
				"failed to close client connection",
				slog.String("error", err.Error()),
			)
		}
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer s.wgClients.Done()

	defer func() {
		if v := recover(); v != nil {
			s.logger.Error(
				"captured panic: ",
				slog.Any("panic", v),
			)
		}

		err := conn.Close()
		if err != nil {
			s.logger.Error(
				"failed to close client connection",
				slog.String("error", err.Error()),
			)
		}
	}()

	defer func() {
		s.logger.Info("closing client connection")
		s.clients.Add(-1)

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

func (s *Server) GetClients() int32 {
	return s.clients.Load()
}

func (s *Server) Stop() error {
	err := s.listen.Close()
	if err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	s.wgClients.Wait()

	return nil
}

func incrementWithLimit(value *atomic.Int32, limit int32) bool {
	for {
		currentValue := value.Load()
		if currentValue >= limit {
			return false
		}

		nextValue := currentValue + 1

		if value.CompareAndSwap(currentValue, nextValue) {
			return true
		}
	}
}
