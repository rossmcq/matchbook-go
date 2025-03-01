package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	matchbook "github.com/rossmcq/matchbook-go/adapter"
	"github.com/rossmcq/matchbook-go/handler"
	"github.com/rossmcq/matchbook-go/postgres"
)

type App struct {
	router          http.Handler
	dbConnection    postgres.DbConnection
	matchbookClient matchbook.Client
	handler         handler.Handler
}

func New(dbConnection postgres.DbConnection, matchbookClient matchbook.Client, handler handler.Handler) (*App, error) {
	app := &App{
		matchbookClient: matchbookClient,
		dbConnection:    dbConnection,
		handler:         handler,
	}

	app.loadRoutes()

	return app, nil
}

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/v1", a.loadOrderRoutes)

	a.router = router
}

func (a *App) loadOrderRoutes(router chi.Router) {

	router.Get("/login", a.handler.Login)
	router.Get("/token", a.handler.GetToken)
	router.Post("/logout", a.handler.Logout)
	router.Get("/health", a.handler.Health)
	router.Post("/event/{id}", a.handler.CreateEvent)
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
