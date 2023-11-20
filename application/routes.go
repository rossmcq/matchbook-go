package application

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rossmcq/matchbook-go/handler"
	"github.com/rossmcq/matchbook-go/matchbook"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/v1", loadOrderRoutes)

	return router
}

func loadOrderRoutes(router chi.Router) {
	sessionToken, err := matchbook.LoadMatchboookToken()
	if err != nil {
		fmt.Printf("Error fetching Matchbook token: %v", err)
	}
	sessionHandler := &handler.Session{
		SessionToken: sessionToken,
	}
	router.Get("/login", sessionHandler.Login)
	router.Get("/token", sessionHandler.GetToken)
	router.Post("/logout", sessionHandler.Logout)
	router.Get("/health", sessionHandler.Health)
	router.Post("/event/{id}", sessionHandler.CreateEvent)
}
