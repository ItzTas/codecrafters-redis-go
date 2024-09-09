package main

import (
	"errors"
	"fmt"
	"sync"
)

type Data struct {
	mu   *sync.RWMutex
	sets map[string][]byte
}

func newData() *Data {
	return &Data{
		mu:   &sync.RWMutex{},
		sets: make(map[string][]byte),
	}
}

var invalidArgsNum = errors.New("invalid number of arguments")

type commandType func([]*RESP) ([]*RESP, error)

func (dat *Data) getCommands() map[string]commandType {
	return map[string]commandType{
		"echo": echoCommand,
		"ping": pingCommand,
		"set":  dat.setCommand,
		"get":  dat.getCommand,
	}
}

func pingCommand(_ []*RESP) ([]*RESP, error) {
	data := []byte("PONG")
	return []*RESP{
		newResp(SimpleString, data),
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

func (d *Data) setCommand(args []*RESP) ([]*RESP, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("%v: want: <key> <value>", invalidArgsNum)
	}

	set := string(args[0].data)

	d.mu.Lock()
	defer d.mu.Unlock()

	d.sets[set] = args[1].data

	toRet := []*RESP{
		newResp(SimpleString, []byte("OK")),
	}

	return toRet, nil
}

func (d *Data) getCommand(args []*RESP) ([]*RESP, error) {
	if len(args) != 1 {
		return nil, invalidArgsNum
	}

	toGet := args[0].data

	d.mu.RLock()
	defer d.mu.RUnlock()

	v, exists := d.sets[string(toGet)]
	if !exists {
		return []*RESP{newNilResp(BulkString)}, nil
	}

	return []*RESP{
		newResp(BulkString, v),
	}, nil
}
