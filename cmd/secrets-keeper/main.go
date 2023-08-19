package main

import (
	"os"

	"github.com/protomem/secrets-keeper/internal/api"
	"github.com/protomem/secrets-keeper/internal/config"
)

func main() {
	var err error

	conf, err := config.New()
	if err != nil {
		os.Exit(1)
	}

	server, err := api.New(conf)
	if err != nil {
		os.Exit(1)
	}

	err = server.Run()
	if err != nil {
		os.Exit(1)
	}
}
