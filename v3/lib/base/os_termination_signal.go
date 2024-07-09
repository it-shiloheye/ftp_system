package base

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var baseContext context.Context = context.Background()

func BaseContext() (ctx context.Context, cancel context.CancelFunc) {
	// Create a context to signal cancellation to the network goroutines
	return signal.NotifyContext(baseContext, os.Interrupt, syscall.SIGTERM)

}
