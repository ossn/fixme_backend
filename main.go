package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ossn/fixme_backend/actions"
	"github.com/ossn/fixme_backend/workers/worker_github"
	"github.com/ossn/fixme_backend/workers/worker_gitlab"
)

// main is the starting point to your Buffalo application.
func main() {
	app := actions.App()

	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	// Go routine to listen for os messages
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	// Start workers
	go worker_github.WorkerInst.Init(ctx, c)
	go worker_gitlab.WorkerInst.Init(ctx, c)


	// Start app serve
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
