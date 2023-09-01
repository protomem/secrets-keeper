package config

import (
	"os"
)

type Config struct {
	BindAddr  string
	LogLevel  string
	Database  string
	CertsName string
}

func New() (Config, error) {
	const op = "config.New"
	var (
		exist bool
		conf  Config
	)

	conf.BindAddr, exist = os.LookupEnv("BIND_ADDR")
	if !exist {
		conf.BindAddr = ":8080"
	}

	conf.LogLevel, exist = os.LookupEnv("LOG_LEVEL")
	if !exist {
		conf.LogLevel = "debug"
	}

	conf.Database, exist = os.LookupEnv("DATABASE")
	if !exist {
		conf.Database = ":memory:"
	}

	conf.CertsName, exist = os.LookupEnv("CERTS_NAME")
	if !exist {
		conf.CertsName = "localhost"
	}

	return conf, nil
}
