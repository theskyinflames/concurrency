package config

import (
	"os"
)

const (
	RPCServerAddr = "CRMINTEGRATOR_RPC_SERVER_ADDRESS"
)

type (
	Config struct {
		RPCServerAddr string
	}
)

func (c *Config) Load() (err error) {
	c.RPCServerAddr = getEnv(RPCServerAddr)
	return
}

func getEnv(env string) (value string) {
	value = os.Getenv(env)
	if len(value) == 0 {
		panic("environment variable " + env + " does not exist")
	}
	return
}
