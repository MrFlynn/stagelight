package main

import (
	ctx "context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "api",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "address",
				Usage: "The address to listen on.",
				Value: "127.0.0.1",
			},
			&cli.UintFlag{
				Name:  "port",
				Usage: "The port to listen on.",
				Value: 8000,
			},
			&cli.PathFlag{
				Name:     "database",
				Usage:    "Path to bbolt database for program data.",
				Required: true,
			},
		},
		Action: run,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Nick Pleatsikas",
				Email: "nick@pleatsikas.me",
			},
		},
	}

	app.Run(os.Args)
}

func run(context *cli.Context) error {
	srv := createServer(context.String("address"), context.Uint("port"), context.Path("database"))

	go func() {
		log.Printf("Server listening on %s:%d\n", context.String("address"), context.Uint("port"))
		log.Fatal(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	var wait time.Duration
	ctx, cancel := ctx.WithTimeout(ctx.Background(), wait)

	defer cancel()

	srv.Shutdown(ctx)

	fmt.Printf("\n")
	log.Println("Server shutting down. Goodbye...")

	return nil
}
