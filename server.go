package mynet

import (
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

type Config struct {
	Addr         string
	MaxConn      int
	Logger       *log.Logger
	DeadlineTime time.Duration
}

type Server struct {
	config   *Config
	mu       sync.Mutex
	listener net.Listener
	conns    map[uint64]*Conn
	handler  Handler
	protocol Protocol
}

type Handler interface {
	Open(conn *Conn)
	Data(conn *Conn, data []byte)
	Close(conn *Conn)
}

var (
	linkId uint64
	Log    *SimpleLog
)

func NewServer(conf *Config, h Handler, p Protocol) *Server {
	server := new(Server)
	server.config = conf
	server.protocol = p
	server.conns = make(map[uint64]*Conn)
	server.handler = h

	if conf.Logger == nil {
		panic("logger is not set")
	}
	Log = NewLogger(conf.Logger)

	var err error
	server.listener, err = net.Listen("tcp", server.config.Addr)
	if err != nil {
		panic(err.Error())
	}

	Log.Info("listen %s", server.config.Addr)

	return server
}

func (server *Server) Run() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			server.Close()
			break
		}

		server.newConn(conn)
	}
}

func (server *Server) newConn(conn net.Conn) {
	server.mu.Lock()
	defer server.mu.Unlock()

	if len(server.conns) > server.config.MaxConn {
		Log.Info("maximum connection limit", conn.RemoteAddr())
		return
	}

	linkId++
	c := NewConn(server, linkId, conn, server.protocol, 32)
	server.conns[c.Id()] = c
	go c.readLoop()
}

func (server *Server) delConn(c *Conn) {
	server.mu.Lock()
	defer server.mu.Unlock()

	delete(server.conns, c.Id())
}

func (server *Server) Close() {
	server.mu.Lock()
	defer server.mu.Unlock()

	server.listener.Close()

	for _, conn := range server.conns {
		conn.Close(nil)
	}
}

func Recover() {
	if err := recover(); err != nil {
		Log.Error("%v", err)
		for i := 1; i <= 20; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			Log.Error("%v %v", file, line)
		}
	}
}
