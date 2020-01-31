package grpc_quick

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SignalContext does cool magic see github.com/rodaine/grpc-chat.git
func SignalContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		signal.Stop(sigs)
		close(sigs)
		cancel()
	}()

	return ctx
}