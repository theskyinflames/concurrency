package config

import (
	"os"
	"strconv"
)

const (
	RPCCRMServerAddr     = "CSVREADER_CRMINTEGRATOR_RPC_SERVER_ADDRESS"
	CSVRecordsBlock      = "CSVREADER_RECORDS_BLOCK_SIZE"
	MaxAttempsPerRecord  = "CSVREADER_MAX_ATTEMPTS_PER_RECORD"
	MaxConcurrentSenders = "CSVREADER_MAX_CONCURRENTS_SENDERS"
	FilePath             = "CSVREADER_FILE_PATH"
	PostgreSQLConnStr    = "CSVREADER_POSTGRESQL_CONN_STR" // Example: postgres://pqgotest:password@localhost/pqgotest
)

type (
	Config struct {
		RPCCRMServerAddr     string
		CSVRecordsBlock      int64
		MaxAttempsPerRecord  int
		MaxConcurrentSenders int
		FilePath             string
		PostgreSQLConnStr    string
	}
)

func (c *Config) Load() (err error) {
	c.RPCCRMServerAddr = getEnv(RPCCRMServerAddr)
	c.CSVRecordsBlock, err = strconv.ParseInt(getEnv(CSVRecordsBlock), 10, 64)
	if err == nil {
		c.MaxAttempsPerRecord, err = strconv.Atoi(getEnv(MaxAttempsPerRecord))
	}
	if err == nil {
		c.MaxConcurrentSenders, err = strconv.Atoi(getEnv(MaxConcurrentSenders))
	}
	if err == nil {
		c.FilePath = getEnv(FilePath)
	}
	if err == nil {
		c.PostgreSQLConnStr = getEnv(PostgreSQLConnStr)
	}
	return
}

func getEnv(env string) (value string) {
	value = os.Getenv(env)
	if len(value) == 0 {
		panic("environment variable " + env + " does not exist")
	}
	return
}
