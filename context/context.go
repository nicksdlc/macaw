package context

import (
	"fmt"

	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
)

// Context is a context with which macaw is running
type Context struct {
	communicator communicators.Communicator
	runner       runner
	cfg          *config.Configuration
}

// BuildContext creates new context for macaw
func BuildContext(cfg config.Configuration) (*Context, error) {
	ctx := Context{}
	return ctx.Build(&cfg)
}

// Build creates a context of macaw that is able to be executed
func (ctx *Context) Build(cfg *config.Configuration) (*Context, error) {
	ctx.cfg = cfg
	if ctx.cfg.ConnectThrough == "" {
		return nil, fmt.Errorf("connection profile is not defined")
	}

	c, err := BuildCommunicator(ctx.cfg)
	if err != nil {
		return nil, err
	}

	ctx.communicator = c
	ctx.runner = get(ctx.communicator, *ctx.cfg)

	return ctx, nil
}

// Run the context
func (ctx *Context) Run() {
	ctx.runner(ctx.communicator, *ctx.cfg)
}
