package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer closeIt(l, "")

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	for {

		buf := make([]byte, 2048)

		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				conn, err = l.Accept()
				if err != nil {
					fmt.Println("Error accepting connection: ", err.Error())
					os.Exit(1)
				}
				continue
			}
			log.Fatalf("Error reading from conection: %v", err)
		}

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			log.Fatalf("Error writing to conection: %v", err)
		}
	}
}

func closeIt(c io.Closer, msg string) {
	err := c.Close()
	if err != nil {
		fmt.Println(err, " ", msg)
	}
}
