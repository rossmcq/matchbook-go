package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

type App struct {
	router         http.Handler
	dbConnection   *sql.DB
	matchbookToken string
	config         Config
}

func New() (*App, error) {
	app := &App{
		config: LoadConfig(),
	}

	dbConnection, err := sql.Open("postgres", app.config.dbConnectionString)
	if err != nil {
		return &App{}, fmt.Errorf("Can't open DB: %v", err)
	}

	app.dbConnection = dbConnection

	app.loadRoutes()

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	// connect to psql
	err := a.dbConnection.Ping()
	if err != nil {
		return fmt.Errorf("Can't ping Db: %v", err)
	}

	defer func() {
		if err := a.dbConnection.Close(); err != nil {
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
