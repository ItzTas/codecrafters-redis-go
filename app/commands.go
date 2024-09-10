package main

import (
	"errors"
	"strconv"
	"time"
)

var (
	invalidArgsNum = errors.New("invalid number of arguments")
	invalidArg     = errors.New("invalid arg")
)

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
	setAtgs, err := parseSetArgs(args)
	if err != nil {
		return []*RESP{}, err
	}

	set := string(args[0].data)

	d.setSetData(set, args[1].data, setAtgs)

	toRet := []*RESP{
		newResp(SimpleString, []byte("OK")),
	}

	return toRet, nil
}

func parseSetArgs(args []*RESP) (SetArgs, error) {
	var setArgs SetArgs
	for i, arg := range args {
		if string(arg.data) == "px" {
			if i+1 == len(args) {
				return SetArgs{}, invalidArgsNum
			}

			pxInt, err := strconv.Atoi(string(args[i+1].data))
			if err != nil {
				return SetArgs{}, err
			}

			setArgs.expiry = time.Duration(pxInt) * time.Millisecond
			continue
		}

		if string(arg.data) == "ex" {
			if i+1 == len(args) {
				return SetArgs{}, invalidArgsNum
			}

			pxInt, err := strconv.Atoi(string(args[i+1].data))
			if err != nil {
				return SetArgs{}, err
			}

			setArgs.expiry = time.Duration(pxInt) * time.Second
			continue

		}
	}
	return setArgs, nil
}

func (d *Data) getCommand(args []*RESP) ([]*RESP, error) {
	if len(args) != 1 {
		return nil, invalidArgsNum
	}

	toGet := args[0].data

	d.mu.RLock()
	defer d.mu.RUnlock()

	v, exists := d.getSetData(string(toGet))
	if !exists {
		return []*RESP{newNilResp(BulkString)}, nil
	}

	return []*RESP{
		newResp(BulkString, v),
	}, nil
}
