package main

import (
	"net"
	"time"
)

type Config struct {
	el    *EventLoop
	data  *Data
	dbCfg *DBConfig
}

func NewConfig(l net.Listener, reapInterval time.Duration) *Config {
	dbCfg := GetInitialDBConfig()

	return &Config{
		el: &EventLoop{
			l: l,
		},
		dbCfg: dbCfg,
		data:  newData(reapInterval),
	}
}
