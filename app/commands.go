package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	invalidArgsNum = errors.New("invalid number of arguments")
	invalidArg     = errors.New("invalid arg")
)

type commandType func([]*RESP) ([]*RESP, error)

func (cfg *Config) getCommands() map[string]commandType {
	return map[string]commandType{
		"echo":   echoCommand,
		"ping":   pingCommand,
		"set":    cfg.setCommand,
		"get":    cfg.getCommand,
		"config": cfg.executeConfigCommands,
		"keys":   cfg.keysCommand,
	}
}

func (cfg *Config) executeConfigCommands(args []*RESP) ([]*RESP, error) {
	if len(args) == 0 {
		return []*RESP{}, fmt.Errorf("%v: want CONFIG + command", invalidArgsNum)
	}
	commandStr := string(args[0].data)
	commandStr = strings.ToLower(commandStr)
	command, exists := cfg.getConfigCommands()[commandStr]
	if !exists {
		return []*RESP{}, fmt.Errorf("%v: %s does not exist", invalidArg, commandStr)
	}

	return command(args[1:])
}

func (cfg *Config) getConfigCommands() map[string]commandType {
	return map[string]commandType{
		"get": cfg.dbCfg.getCommand,
		// "set": cfg.configSetCommand,
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

func (cfg *Config) setCommand(args []*RESP) ([]*RESP, error) {
	setAtgs, err := parseSetArgs(args)
	if err != nil {
		return []*RESP{}, err
	}

	if len(args) == 0 {
		return []*RESP{}, invalidArgsNum
	}

	set := string(args[0].data)

	cfg.data.setSetData(set, args[1].data, setAtgs)

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
