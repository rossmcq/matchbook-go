package handler

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"

	matchbook "github.com/rossmcq/matchbook-go/adapter"
	"github.com/rossmcq/matchbook-go/service"
)

type Handler struct {
	service service.Service
}

func New(s service.Service) Handler {
	return Handler{
		service: s,
	}
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var err error
	h.service.MatchbookClient, err = matchbook.New()
	if err != nil {
		fmt.Printf("Error loading token %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Printf("Got session token %v \n", h.service.MatchbookClient.Token)
}

func (h Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Current session token %v \n", h.service.MatchbookClient.Token)
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.service.LogoutMatchbook()
	if err != nil {
		fmt.Printf("logout error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	fmt.Printf("handle_request: create event data for id: %v \n", idParam)

	err := h.service.CreateEvent(idParam)
	if err != nil {
		fmt.Printf("error creating event: %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Handler) Health(w http.ResponseWriter, r *http.Request) {
	err := s.service.Health()
	if err != nil {
		fmt.Printf("Error connecting to DB: %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
