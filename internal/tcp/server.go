package tcp

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
)

type Server struct {
	addr   string
	listen net.Listener
}

func NewServer(addr string) (*Server, error) {
	server := Server{
		addr: addr,
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
			// TODO think what to do: stop server or try to accept another connection?
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
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			slog.Error(
				"failed to read data from tcp client",
				slog.String("error", err.Error()),
			)

			return
		}

		fmt.Print("-> ", netData)

		_, err = conn.Write([]byte(netData))
		if err != nil {
			slog.Error(
				"failed to send data to tcp client",
				slog.String("error", err.Error()),
			)

			return
		}
	}
}
