package handler

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

type Session struct{}

var matchbooktoken string

func (s Session) Login(w http.ResponseWriter, r *http.Request) {
	var err error
	matchbooktoken, err = loadMatchboookToken()
	if err != nil {
		fmt.Printf("Error loading token %v", err)
	}
	fmt.Printf("Got session token %v", matchbooktoken)
}

func (s Session) GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Current session token %v", matchbooktoken)
}

func (s Session) Logout(w http.ResponseWriter, r *http.Request) {
	response := logoutMatchbook(&matchbooktoken)
	fmt.Printf("Logout of current session %v", response)
}

func (s Session) CreateEvent(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	fmt.Printf("Create event data for Id: %v", idParam)

	marketID, description, err := getMatchOddsMarketId(idParam)
	if err != nil {
		fmt.Printf("Error getting market id %v", err)
	}
	fmt.Printf("MarketId: %f Description: %v", marketID, description)
}
