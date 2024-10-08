package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rossmcq/matchbook-go/postgres"
)

type App struct {
	router         http.Handler
	dbConnection   postgres.DbConnection
	matchbookToken string
}

func New(dbConnection postgres.DbConnection, matchbookToken string) (*App, error) {
	app := &App{
		matchbookToken: matchbookToken,
		dbConnection:   dbConnection,
	}

	app.loadRoutes()

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	// Test connection to psql
	err := a.dbConnection.Database.Ping()
	if err != nil {
		return fmt.Errorf("can't ping Db: %v", err)
	}

	defer func() {
		if err := a.dbConnection.Database.Close(); err != nil {
			fmt.Println("failed to close dB", err)
		}
	}()

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to listen to server %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
