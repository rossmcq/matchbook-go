package handler

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"

	"github.com/rossmcq/matchbook-go/matchbook"
	"github.com/rossmcq/matchbook-go/postgres"
)

type Session struct{}

var matchbookToken string

func (s Session) Login(w http.ResponseWriter, r *http.Request) {
	var err error
	matchbookToken, err = matchbook.LoadMatchboookToken()
	if err != nil {
		fmt.Printf("Error loading token %v", err)
	}
	fmt.Printf("Got session token %v", matchbookToken)
}

func (s Session) GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Current session token %v", matchbookToken)
}

func (s Session) Logout(w http.ResponseWriter, r *http.Request) {
	response := matchbook.LogoutMatchbook(&matchbookToken)
	fmt.Printf("Logout of current session %v", response)
}

func (s Session) CreateEvent(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	fmt.Printf("Create event data for Id: %v", idParam)

	marketID, description, err := matchbook.GetMatchOddsMarketId(idParam)
	if err != nil {
		fmt.Printf("Error getting market id %v", err)
	}
	fmt.Printf("MarketId: %f Description: %v", marketID, description)
}

func (s Session) Health(w http.ResponseWriter, r *http.Request) {
	err := postgres.CheckConnection()
	if err != nil {
		fmt.Printf("Error connecting to DB: %v", err)
	}
}
