package mynet

import (
	"errors"
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	server *Server

	id       uint64
	conn     net.Conn
	protocol Protocol

	sendChan  chan []byte
	closeChan chan bool
	closeFlag int32
}

var (
	readBufSize  int = 2048
	writeBufSize int = 8192
)

func NewConn(server *Server, id uint64, conn net.Conn, p Protocol, sendChSize int) *Conn {
	c := &Conn{
		server:    server,
		id:        id,
		conn:      conn,
		protocol:  p,
		sendChan:  make(chan []byte, sendChSize),
		closeChan: make(chan bool),
	}
	c.server.handler.Open(c)
	go c.sendLoop()
	return c
}

func (c *Conn) Id() uint64 {
	return c.id
}

func (c *Conn) read(buf []byte) (n int, err error) {
	return c.protocol.Read(c.conn, buf)
}

func (c *Conn) readLoop() {
	defer Recover()

	Log.Info("new conn id=%d", c.id)

	buf := make([]byte, readBufSize)
	for {
		n, err := c.read(buf)
		if err != nil {
			c.Close(err)
			break
		}
		c.server.handler.Data(c, buf[:n])
	}
}

func (c *Conn) Send(b []byte, timeout time.Duration) error {
	if c.IsClosed() {
		return errors.New("conn close")
	}

	select {
	case c.sendChan <- b:
	case <-time.After(timeout):
		return errors.New("send timeout")
	}
	return nil
}

func (c *Conn) send(msg []byte, buf []byte) error {
	return c.protocol.Write(c.conn, msg, buf)
}

func (c *Conn) sendLoop() {
	buf := make([]byte, writeBufSize)
	for {
		select {
		case msg := <-c.sendChan:
			err := c.send(msg, buf)
			if err != nil {
				c.Close(err)
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}

func (c *Conn) Close(err error) {
	if atomic.CompareAndSwapInt32(&c.closeFlag, 0, 1) {
		c.server.handler.Close(c)
		c.server.delConn(c)

		c.conn.Close()
		close(c.closeChan)
		Log.Info("conn close id=%d err=%v", c.id, err)
	}
}

func (c *Conn) SetDeadline(d time.Duration) {
	c.conn.SetDeadline(time.Now().Add(d))
}
