package main

func (cfg *Config) keysCommand(args []*RESP) ([]*RESP, error) {
	if string(args[0].data) != "*" {
		return []*RESP{}, invalidArg
	}

	if cfg.dbReader.noFile {
		return []*RESP{newNilResp(BulkString)}, nil
	}
	keys, err := cfg.dbReader.readDatabase()
	if err != nil {
		return []*RESP{}, err
	}

	var result string

	for _, key := range keys {
		result += " " + key.key
	}

	res := stringToArrayOfBulkResp(result)
	return []*RESP{res}, nil
}
