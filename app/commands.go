package main

import "errors"

var invalidArgsNum = errors.New("invalid number of arguments")

type commandType func([]*RESP) ([]*RESP, error)

func getCommands() map[string]commandType {
	return map[string]commandType{
		"echo": echoCommand,
		"ping": pingCommand,
	}
}

func pingCommand(_ []*RESP) ([]*RESP, error) {
	data := []byte("PONG")
	return []*RESP{
		{
			st:    SimpleString,
			data:  data,
			count: len(data),
		},
	}, nil
}

func echoCommand(args []*RESP) ([]*RESP, error) {
	if len(args) != 1 {
		return []*RESP{}, invalidArgsNum
	}

	return []*RESP{
		args[0],
	}, nil
}
