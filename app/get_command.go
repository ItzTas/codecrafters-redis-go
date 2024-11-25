package main

func (cfg *Config) getCommandDB(key string) ([]*RESP, error) {
	keys, err := cfg.dbReader.readDatabase()
	if err != nil {
		return []*RESP{}, err
	}

	val, ok := getValFromKeys(keys, key)
	if !ok {
		return []*RESP{newNilResp(BulkString)}, nil
	}
	return []*RESP{
		newResp(BulkString, []byte(val)),
	}, nil
}

func (cfg *Config) getCommand(args []*RESP) ([]*RESP, error) {
	if len(args) != 1 {
		return nil, invalidArgsNum
	}

	toGet := args[0].data

	v, exists := cfg.data.getSetData(string(toGet))
	if !exists {
		return cfg.getCommandDB(string(toGet))
	}

	return []*RESP{
		newResp(BulkString, v),
	}, nil
}
