package main

import (
	"flag"
	"sync"
)

type DBConfig struct {
	configs map[string]string
	mu      *sync.RWMutex
}

const (
	defaultDir        = "/var/lib/redis"
	defaultDBFilename = "dump.rdb"
)

func (dbConfig *DBConfig) getConfigCamp(camp string) (c string, exists bool) {
	dbConfig.mu.RLock()
	defer dbConfig.mu.RUnlock()
	c, exists = dbConfig.configs[camp]
	return c, exists
}

func GetInitialDBConfig() *DBConfig {
	dir := flag.String("dir", defaultDir, "the path to the directory where the RDB file is stored")
	dbFilename := flag.String("dbfilename", defaultDBFilename, "the name of the RDB file")
	flag.Parse()

	dbConfig := &DBConfig{
		configs: make(map[string]string),
		mu:      &sync.RWMutex{},
	}

	dbConfig.mu.Lock()
	defer dbConfig.mu.Unlock()

	dbConfig.configs["dir"] = *dir
	dbConfig.configs["dbfilename"] = *dbFilename

	return dbConfig
}
