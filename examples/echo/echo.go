package main

import (
	"github.com/soki/mynet"
	"log"
	"os"
	"time"
)

func main() {
	config := &mynet.Config{
		Addr:    "0.0.0.0:8760",
		MaxConn: 100,
		Logger:  log.New(os.Stdout, "", 0),
	}
	protocol := mynet.NewSimpleProtocol(1024, 1024, 4)
	server := mynet.NewServer(config, NewService(), protocol)
	server.Run()
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Open(conn *mynet.Conn) {
	s.SendMsg(conn, []byte("hello world"))
}

func (s *Service) Data(conn *mynet.Conn, data []byte) {
	log.Printf("msg id=%d data=%s\n", conn.Id(), string(data))
	s.SendMsg(conn, data)
}

func (s *Service) Close(conn *mynet.Conn) {
	log.Printf("close id=%d\n", conn.Id())
}

func (s *Service) SendMsg(conn *mynet.Conn, b []byte) error {
	return conn.Send(b, time.Duration(1*time.Second))
}
