package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	srv := createServer()

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	var wait time.Duration
	ctx, cancel := context.WithTimeout(context.Background(), wait)

	defer cancel()

	srv.Shutdown(ctx)

	fmt.Println("")
	log.Println("Server shutting down. Goodbye...")
	os.Exit(0)
}
