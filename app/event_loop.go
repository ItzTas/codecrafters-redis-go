package main

import (
	"fmt"
	"io"
	"net"
)

type EventLoop struct {
	l net.Listener
}

func (el *EventLoop) runRedis() {
	for {
		conn, err := el.l.Accept()
		if err != nil {
			fmt.Println("Error accepting conection: ", err)
			return
		}

		msgs := make(chan []byte)

		go el.loopEvent(msgs, conn)
		go func() {
			for {
				buffer := make([]byte, 2048)
				_, err := conn.Read(buffer)
				if err != nil {
					defer func() {
						defer closeIt(conn)
						defer close(msgs)
					}()
					if err == io.EOF {
						return
					}
					fmt.Println("Could not read from conection: ", err)
					return
				}

				msgs <- buffer
			}
		}()
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
