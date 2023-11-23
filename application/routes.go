package application

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rossmcq/matchbook-go/handler"
	"github.com/rossmcq/matchbook-go/matchbook"
	"github.com/rossmcq/matchbook-go/postgres"

)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/v1", a.loadOrderRoutes)
	sessionToken, err := matchbook.LoadMatchboookToken()
	if err != nil {
		fmt.Printf("Error fetching Matchbook token: %v", err)
	}

	a.matchbookToken = *sessionToken

	a.router = router
}


func (a *App) loadOrderRoutes(router chi.Router) {
	sessionHandler := &handler.Session{
		SessionToken: &a.matchbookToken,
		DbConnection: &postgres.DbConnection{
			Database: a.dbConnection,
		},

	}
	router.Get("/login", sessionHandler.Login)
	router.Get("/token", sessionHandler.GetToken)
	router.Post("/logout", sessionHandler.Logout)
	router.Get("/health", sessionHandler.Health)
	router.Post("/event/{id}", sessionHandler.CreateEvent)
}
