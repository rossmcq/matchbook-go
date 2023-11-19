package application

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rossmcq/matchbook-go/handler"
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
	sessionHandler := &handler.Session{}
	router.Get("/login", sessionHandler.Login)
	router.Get("/token", sessionHandler.GetToken)
	router.Post("/logout", sessionHandler.Logout)
	router.Get("/health", sessionHandler.Health)
	router.Post("/event/{id}", sessionHandler.CreateEvent)
}
