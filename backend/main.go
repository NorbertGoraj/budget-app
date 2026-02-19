package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"budget-app/appcontext"
)

func main() {
	log.Println("Starting budget-app.")

	ctx, err := appcontext.NewContext()
	if err != nil {
		log.Fatal("failed to initialize application context: ", err)
	}
	defer ctx.CancelF()

	go waitShutdown(ctx.CancelF)

	<-ctx.Done()
	ctx.WaitGroup.Wait()
	log.Println("see you!")
}

func waitShutdown(cancelF context.CancelFunc) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT)

	s := <-sigint
	log.Printf("os signal received: %d (%s)\n", s, s)
	cancelF()
}
