package worker

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func gracefulStop(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		var signals = make(chan os.Signal)
		defer close(signals)

		signal.Notify(signals, syscall.SIGINT)
		signal.Notify(signals, syscall.SIGQUIT)
		signal.Notify(signals, syscall.SIGTERM)
		defer signalStop(signals)

		<-signals
		defer cancel()
	}()
}

// Stops signals channel
func signalStop(c chan<- os.Signal) {
	signal.Stop(c)
}
