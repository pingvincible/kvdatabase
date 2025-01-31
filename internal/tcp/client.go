package tcp

import (
	"bufio"
	"fmt"
	"net"

	"github.com/pingvincible/kvdatabase/internal/kvio"
)

type Client struct {
	conn       net.Conn
	ReadWriter *kvio.ReadWriter
}

func NewClient(addr string) (*Client, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tcp address: %w", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tcp address: %w", err)
	}

	return &Client{
		conn:       conn,
		ReadWriter: kvio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
	}, nil
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close tcp connection: %w", err)
	}

	return nil
}
