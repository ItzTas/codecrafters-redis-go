package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379", err)
		os.Exit(1)
	}
	defer closeIt(l)

	cfg := NewConfig(l, 1*time.Second)

	cfg.runRedis()
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
