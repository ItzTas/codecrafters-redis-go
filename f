package main

import (
	"fmt"
	"net"

	"github.com/golangci/golangci-lint/pkg/golinters/funlen"
)

type EventLoop struct {
	l net.Listener
}

func (el *EventLoop) runRedis() {
	for {
		conn, err := el.l.Accept()
		if err != nil {
			fmt.Println("Error accepting conection", err)
			return
		}

		msgs := make(chan []byte)

		go el.loopEvent(msgs, conn)
		go func() {}()
	}
}

func (el *EventLoop) loopEvent(msgs <-chan []byte, conn net.Conn) {
	for range msgs {
		_, err := conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Could not write to conection: ", err)
			continue
		}
	}
}
