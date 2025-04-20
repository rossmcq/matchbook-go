package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	matchbook "github.com/rossmcq/matchbook-go/adapter"
	"github.com/rossmcq/matchbook-go/application"
	"github.com/rossmcq/matchbook-go/handler"
	"github.com/rossmcq/matchbook-go/service"
	"github.com/rossmcq/matchbook-go/service/poller"
	"github.com/rossmcq/matchbook-go/store"

	_ "github.com/golang/mock/mockgen/model"
)

func main() {
	dbConnection, err := store.New()
	if err != nil {
		log.Fatalf("unable to create db connection: %e", err)
	}

	matchbookClient, err := matchbook.New()
	if err != nil {
		log.Fatalf("unable to fetch matchbook token: %e", err)

	}

	service, err := service.New(matchbookClient, &dbConnection)
	if err != nil {
		log.Fatalf("unable to create service: %e", err)
	}

	handler, err := handler.New(service)
	if err != nil {
		log.Fatalf("unable to create handler: %e", err)
	}

	app, err := application.New(dbConnection, matchbookClient, handler)
	if err != nil {
		log.Fatalf("failed initiating app %e", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	//start polling
	EventPoller := poller.New(poller.Poller{}, service.RecordMatchOdds, func(err error) {
		log.Printf("error polling event: %v", err)
	})
	EventPoller.Start(ctx, 10*time.Second)

	err = app.Start(ctx)
	if err != nil {
		log.Fatalf("failed to listen to server %e", err)
	}
}
