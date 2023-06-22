package main

import (
	"log"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/context"
)

func main() {
	cfg := readConfig()

	ctx, err := context.BuildContext(cfg)
	if err != nil {
		log.Panic(err.Error())
	}

	ctx.Run()
}

func readConfig() config.Configuration {
	return config.Read("config_rabbit")
}
