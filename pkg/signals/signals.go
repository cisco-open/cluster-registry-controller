// Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func NotifyContext(parentContext context.Context) context.Context {
	ctx, cancel := context.WithCancel(parentContext)
	c := make(chan os.Signal, 1)
	signal.Notify(c, []os.Signal{os.Interrupt, syscall.SIGTERM}...)

	go func() {
		select {
		case <-c:
			cancel()
			<-c
			os.Exit(1) // second signal. Exit directly.
		case <-ctx.Done():
		}
	}()

	return ctx
}
