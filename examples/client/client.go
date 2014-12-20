package main

import (
	"bufio"
	"log"
	"mynet"
	"os"
	"time"
)

var (
	buf    []byte
	client *mynet.Client
)

func main() {
	buf = make([]byte, 1024)

	var err error
	client, err = mynet.DialTimeout(
		"tcp",
		"127.0.0.1:8760",
		time.Duration(100*time.Millisecond),
		mynet.NewSimpleProtocol(1024, 1024),
	)
	if err != nil {
		panic(err.Error())
	}

	go client.ReadLoop(handler)

	rd := bufio.NewReader(os.Stdin)
	for {
		l, isPrefix, err := rd.ReadLine()
		if isPrefix || err != nil {
			break
		}
		err = send(l)
		log.Println("send msg", string(l), err)
	}
}

func handler(b []byte) {
	log.Println("get msg", string(b))
}

func send(b []byte) error {
	return client.Send(b, buf)
}
