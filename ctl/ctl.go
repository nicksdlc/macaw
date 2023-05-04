package ctl

import (
	"fmt"
	"net/http"

	"github.com/nicksdlc/macaw/context"
)

type Control struct {
	port uint16
	ctx  *context.Context
}

func New(port uint16, ctx *context.Context) *Control {
	return &Control{
		port: port,
		ctx:  ctx,
	}
}

func (c *Control) Serve() {
	healthcheck := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", healthcheck)
	http.ListenAndServe(fmt.Sprintf(":%d", c.port), mux)
}

func (c *Control) AddResponse(response string) {
	// listen on the port for post request with new response
	// add response to context
	// respond with ok
	// TODO: implement
	panic("not implemented")
}
