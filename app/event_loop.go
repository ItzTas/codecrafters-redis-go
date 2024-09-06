package main

import (
	"fmt"
	"io"
	"net"
	"strings"
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
	for m := range msgs {

		r := NewReader(m)

		_, err := r.readResp()
		if err != nil {
			fmt.Println(err)
			return
		}

		commandStr := r.getCommand()

		commandStr = strings.ToLower(commandStr)

		command, exists := getCommands()[commandStr]
		if !exists {
			err := respondWithError(conn, fmt.Sprintf("Comand: %s does not exist", commandStr))
			if err != nil {
				fmt.Println(err)
				return
			}
			continue
		}

		args := r.getArgs()

		toWir, err := command(args)
		if err != nil {
			err := respondWithError(conn, err.Error())
			if err != nil {
				fmt.Println(err)
				return
			}
			continue
		}

		err = respondToClient(conn, toWir)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}
