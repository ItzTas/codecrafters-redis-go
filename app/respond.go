package main

import (
	"fmt"
	"net"
	"strings"
)

func (r *RESP) marshal() (string, error) {
	switch r.st {
	case SimpleString:
		return fmt.Sprintf("+%s\r\n", string(r.data)), nil
	case Integer:
		return fmt.Sprintf(":%s", string(r.data)), nil
	case ErrorType:
		return fmt.Sprintf("-%s\r\n", string(r.data)), nil
	case BulkString:
		if r.data == nil {
			return "$-1\r\n", nil
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", r.count, string(r.data)), nil
	case Array:
		var sb strings.Builder
		for _, item := range r.array {
			s, err := item.marshal()
			if err != nil {
				return "", err
			}
			sb.WriteString(s)
		}
		return fmt.Sprintf("*%d\r\n%s", r.count, sb.String()), nil
	}

	return "", invalidType
}

func respondWithError(conn net.Conn, message string) error {
	_, err := conn.Write([]byte("-" + message + "\r\n"))
	return err
}

func respondToClient(conn net.Conn, payload []*RESP) error {
	var psMarshaled string
	for _, p := range payload {
		if p.st == Array {
			psMarshaled += fmt.Sprintf("*%d\r\n%s", p.count, psMarshaled)
		}
		s, err := p.marshal()
		if err != nil {
			return err
		}

		psMarshaled += s
	}

	_, err := conn.Write([]byte(psMarshaled))
	return err
}
