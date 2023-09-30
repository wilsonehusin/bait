package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"go.husin.dev/bait"
)

var (
	configFilePath = os.Getenv("BAIT_CONFIG")
)

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	conf, err := bait.NewConfigFromFile(configFilePath)
	if err != nil {
		return err
	}

	done := make(chan bool)
	srv := bait.NewServer(ctx, conf.Config)
	go func() {
		srv.ListenAndServe()
		log.Printf("server stopped")
		done <- true
	}()

	<-ctx.Done()
	log.Printf("shutting down server")
	srv.Shutdown(context.Background())
	<-done
	return ctx.Err()
}

func main() {
	if err := run(); err != nil {
		log.Print(err.Error())
	}
}
