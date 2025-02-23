package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/rossmcq/matchbook-go/application"
	"github.com/rossmcq/matchbook-go/handler"
	"github.com/rossmcq/matchbook-go/matchbook"
	"github.com/rossmcq/matchbook-go/postgres"
)

func main() {
	dbConnection, err := postgres.New()
	if err != nil {
		log.Fatalf("unable to create db connection: %e", err)
	}

	matchbookClient, err := matchbook.New()
	if err != nil {
		log.Fatalf("unable to fetch matchbook token: %e", err)
	}

	handler := handler.New(matchbookClient, &dbConnection)
	// Check error

	//init app passing in dependancies
	app, err := application.New(dbConnection, matchbookClient, handler)
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
