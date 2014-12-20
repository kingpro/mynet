package mynet

import (
	"net"
	"time"
)

type Client struct {
	conn     net.Conn
	protocol Protocol
}

func DialTimeout(network string, address string, timeout time.Duration, p Protocol) (*Client, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}
	c := &Client{
		conn:     conn,
		protocol: p,
	}
	return c, nil
}

func (c *Client) ReadLoop(handler func(b []byte)) {
	defer Recover()

	buf := make([]byte, 2048)
	for {
		n, err := c.protocol.Read(c.conn, buf)
		if err != nil {
			break
		}
		handler(buf[:n])
	}
}

func (c *Client) Send(b []byte, buf []byte) error {
	return c.protocol.Write(c.conn, b, buf)
}
