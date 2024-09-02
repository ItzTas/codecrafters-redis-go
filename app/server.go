package main

import (
	"fmt"
	"io"
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
	defer closeIt(l)

	el := &EventLoop{
		l: l,
	}

	el.runRedis()
}

func closeIt(c io.Closer, msg ...string) {
	err := c.Close()
	m := " "
	for _, ms := range msg {
		m += " " + ms
	}
	if err != nil {
		fmt.Println(err, m)
	}
}
