package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"

	"github.com/rossmcq/matchbook-go/matchbook"
	"github.com/rossmcq/matchbook-go/model"
	"github.com/rossmcq/matchbook-go/postgres"
)

type Session struct {
	matchbookClient matchbook.Client
	DbConnection    *postgres.DbConnection
}

func New(m matchbook.Client, db *postgres.DbConnection) Session {
	return Session{
		matchbookClient: m,
		DbConnection:    db,
	}
}

func (s *Session) Login(w http.ResponseWriter, r *http.Request) {
	var err error
	s.matchbookClient, err = matchbook.New()
	if err != nil {
		fmt.Printf("Error loading token %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Printf("Got session token %v \n", s.matchbookClient.Token)
}

func (s *Session) GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Current session token %v \n", s.matchbookClient.Token)
}

func (s *Session) Logout(w http.ResponseWriter, r *http.Request) {
	response, err := s.matchbookClient.LogoutMatchbook(&s.matchbookClient.Token)
	if err != nil {
		fmt.Printf("logout error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Printf("Logout of current session %v \n", response)
}

func (s *Session) CreateEvent(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	fmt.Printf("Create event data for Id: %v \n", idParam)

	marketID, description, err := s.matchbookClient.GetMatchOddsMarketId(idParam)
	if err != nil {
		fmt.Printf("Error getting market id %v \n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	marketIDStr := strconv.FormatInt(marketID, 10)
	game := model.Game{
		GameID:      rand.Uint64(),
		EventID:     idParam,
		MarketID:    marketIDStr,
		Description: description}

	fmt.Printf("MarketId: %v Description: %v \n", game.MarketID, game.Description)

	err = s.DbConnection.InsertOrReturnGameID(r.Context(), game)

	if err != nil {
		fmt.Printf("Error from postgres.insert: %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Session) Health(w http.ResponseWriter, r *http.Request) {
	err := s.DbConnection.CheckConnection()

	if err != nil {
		fmt.Printf("Error connecting to DB: %v \n", err)
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
