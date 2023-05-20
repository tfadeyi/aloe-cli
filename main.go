package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tfadeyi/aloe-cli/cmd"
	"github.com/tfadeyi/aloe-cli/internal/logging"
)

// @aloe name aloe_cli
// @aloe url https://tfadeyi.github.io
// @aloe version v0.0.1
// @aloe description Aloe CLI application

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	log := logging.NewStandardLogger()
	ctx = logging.ContextWithLogger(ctx, log)

	cmd.Execute(ctx)
}
