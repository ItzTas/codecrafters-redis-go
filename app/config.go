package main

import (
	"net"
	"time"
)

type Config struct {
	el       *EventLoop
	data     *Data
	dbCfg    *DBConfig
	dbReader *RDBReader
}

func NewConfig(l net.Listener, reapInterval time.Duration) *Config {
	dbCfg := GetInitialDBConfig()
	filepath := dbCfg.configs["dir"] + "/" + dbCfg.configs["dbfilename"]
	rdbReader, err := newRDBReader(filepath)
	if err != nil {
		panic(err)
	}

	return &Config{
		el: &EventLoop{
			l: l,
		},
		dbCfg:    dbCfg,
		data:     newData(reapInterval),
		dbReader: rdbReader,
	}
}
