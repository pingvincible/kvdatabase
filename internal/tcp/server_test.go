package tcp_test

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/pingvincible/kvdatabase/internal/compute"
	"github.com/pingvincible/kvdatabase/internal/config"
	"github.com/pingvincible/kvdatabase/internal/logger"
	"github.com/pingvincible/kvdatabase/internal/storage/engine"
	"github.com/pingvincible/kvdatabase/internal/tcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func executeClient(tcpAddr *net.TCPAddr, text string) (*net.TCPConn, error) {
	client, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	_, err = client.Write([]byte(text))
	if err != nil {
		return nil, fmt.Errorf("failed to write: %w", err)
	}

	message := make([]byte, 1)

	_, err = client.Read(message)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	return client, nil
}

func TestTcpServerMaxConnectionsSequentially(t *testing.T) {
	t.Parallel()

	const maxConnections = 50

	eng := engine.New()
	computer := compute.NewComputer(eng)

	server, err := tcp.NewServer(config.NetworkConfig{
		Address:        "",
		MaxConnections: maxConnections,
		MaxMessageSize: "1KB",
		IdleTimeout:    2 * time.Minute,
	}, computer, logger.NewDiscardLogger())
	require.NoError(t, err)

	addr, err := server.Addr()
	require.NoError(t, err)

	wgServer := sync.WaitGroup{}
	wgServer.Add(1)

	go func() {
		defer wgServer.Done()
		server.Start()
	}()

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	require.NoError(t, err)

	clients := make([]*net.TCPConn, maxConnections)
	for i := range maxConnections {
		clients[i], err = executeClient(tcpAddr, "unknown\n")

		require.NoError(t, err)
	}

	_, err = executeClient(tcpAddr, "unknown\n")
	require.Error(t, err)

	assert.Equal(t, maxConnections, server.ClientsHandled)
	assert.Equal(t, 1, server.ClientsDiscarded)
	assert.Equal(t, maxConnections, server.GetClients())

	for _, client := range clients {
		if client != nil {
			_ = client.Close()
		}
	}

	_ = server.Stop()

	wgServer.Wait()

	assert.Equal(t, 0, server.GetClients())
}

func TestTcpServerMaxConnectionsConcurrently(t *testing.T) {
	t.Parallel()

	// FIXME if maxConnections > 100, when all tests run together, they fail
	const maxConnections = 50

	const overConnections = 20

	eng := engine.New()
	computer := compute.NewComputer(eng)
	server, err := tcp.NewServer(config.NetworkConfig{
		Address:        "",
		MaxConnections: maxConnections,
		MaxMessageSize: "1KB",
		IdleTimeout:    2 * time.Minute,
	}, computer, logger.NewDiscardLogger())
	require.NoError(t, err)

	addr, err := server.Addr()
	require.NoError(t, err)

	wgServer := sync.WaitGroup{}
	wgServer.Add(1)

	go func() {
		defer wgServer.Done()
		server.Start()
	}()

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	require.NoError(t, err)

	clients := make([]*net.TCPConn, maxConnections+overConnections)

	clientWg := sync.WaitGroup{}
	for index := range maxConnections + overConnections {
		clientWg.Add(1)

		go func(wg *sync.WaitGroup, index int) {
			defer wg.Done()

			clients[index], err = executeClient(tcpAddr, "GET a\n")
		}(&clientWg, index)
	}

	clientWg.Wait()
	assert.Equal(t, maxConnections, server.GetClients())

	for _, client := range clients {
		if client != nil {
			_ = client.Close()
		}
	}

	_ = server.Stop()

	wgServer.Wait()

	assert.Equal(t, maxConnections, server.ClientsHandled)
	assert.Equal(t, overConnections, server.ClientsDiscarded)
	assert.Equal(t, 0, server.GetClients())
}
