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

func (cfg *Config) readDatabase() {
	data, err := cfg.dbReader.readDatabase()
	if err != nil {
		return
	}
	for _, d := range data {
		key := d.key
		value := d.value
		cfg.data.setSetData(key, []byte(value), SetArgs{})
	}

	cfg.dbReader.resetFile()
}

func NewConfig(l net.Listener, reapInterval time.Duration) *Config {
	dbCfg := GetInitialDBConfig()
	filepath := dbCfg.configs["dir"] + "/" + dbCfg.configs["dbfilename"]
	rdbReader, err := newRDBReader(filepath)
	if err != nil {
		panic(err)
	}

	cfg := &Config{
		el: &EventLoop{
			l: l,
		},
		dbCfg:    dbCfg,
		data:     newData(reapInterval),
		dbReader: rdbReader,
	}

	cfg.readDatabase()
	return cfg
}
