package main

import (
	"log"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/context"
	"github.com/nicksdlc/macaw/ctl"
)

func main() {
	cfg := readConfig()

	ctx, err := context.BuildContext(cfg)
	if err != nil {
		log.Panic(err.Error())
	}

	runCtl(cfg, ctx)

	ctx.Run()

}

func readConfig() config.Configuration {
	return config.Read("config")
}

func runCtl(cfg config.Configuration, ctx *context.Context) {
	ctl := ctl.New(cfg.Control.OnPort, ctx)
	go ctl.Serve()
}
