package main

import "fmt"

func (dbCfg *DBConfig) getCommand(args []*RESP) ([]*RESP, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%v: want: get + <config camp>", invalidArgsNum)
	}

	camp := string(args[0].data)

	cfg, exists := dbCfg.getConfigCamp(camp)
	if !exists {
		return nil, fmt.Errorf("%v: %s not found", invalidArg, camp)
	}

	r := stringToArrayOfBulkResp(camp + " " + cfg)

	return []*RESP{r}, nil
}
