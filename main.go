package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/rossmcq/matchbook-go/application"
	"github.com/rossmcq/matchbook-go/postgres"
)

func main() {
	//init db
	dbConnection, err := postgres.New()
	if err != nil {
		log.Fatalf("unable to create db connection: %e", err)
	}

	//init matchbook connection

	//init handlers

	//init app passing in dependancies
	app, err := application.New(dbConnection)
	if err != nil {
		fmt.Println("failed initiating app %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = app.Start(ctx)
	if err != nil {
		fmt.Println("failed to listen to server %w", err)
	}
}
