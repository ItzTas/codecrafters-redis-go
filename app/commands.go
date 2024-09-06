package main

type commandType func([]*RESP) []*RESP

func getCommands() map[string]commandType {
	return map[string]commandType{
		"echo": echoCommand,
		"ping": pingCommand,
	}
}

func pingCommand(_ []*RESP) []*RESP {
	data := []byte("PONG")
	return []*RESP{
		{
			st:    SimpleString,
			data:  data,
			count: len(data),
		},
	}
}

func echoCommand(args []*RESP) []*RESP {
	return []*RESP{
		args[0],
	}
}
