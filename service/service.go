//go:generate mockgen -destination=./mocks/mock_matchbook.go -package=mocks github.com/rossmcq/matchbook-go/service MatchbookClient
//go:generate mockgen -destination=./mocks/mock_postgres.go -package=mocks github.com/rossmcq/matchbook-go/service DbConnection

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/rossmcq/matchbook-go/model"
)

type Service struct {
	MatchbookClient MatchbookClient
	DbConnection    DbConnection
}

type MatchbookClient interface {
	GetEvent(eventId string) (model.EventResponse, error)
	LogoutMatchbook() (string, error)
	GetMatchbookToken() string
}

type DbConnection interface {
	CreateGame(ctx context.Context, game model.Game) error
	CheckConnection() error
}

func New(matchbookClient MatchbookClient, dbConnection DbConnection) (Service, error) {
	if matchbookClient == nil {
		return Service{}, errors.New("matchbook client is nil")
	}

	if dbConnection == nil {
		return Service{}, errors.New("database client is nil")
	}

	return Service{
		MatchbookClient: matchbookClient,
		DbConnection:    dbConnection,
	}, nil
}

func (s Service) CreateEvent(id string) error {
	event, err := s.MatchbookClient.GetEvent(id)
	if err != nil {
		return errors.New("error getting market id")
	}

	// Find the Match Odds market for the event
	markets := event.Markets
	if markets == nil {
		return fmt.Errorf("no markets found for event %s", id)
	}
	marketID := int64(0)
	var description string
	for i := 0; i < len(markets); i++ {
		market := markets[i]
		if market.Name == "Match Odds" {
			marketID = market.Id
			// TODO split timestamp out from here
			description = event.Name + " " + event.Start
		}
	}

	if marketID == 0 {
		return errors.New("no match odds found")
	}

	game := model.Game{
		GameID:      id,
		EventID:     id,
		MarketID:    marketID,
		Description: description,
	}

	fmt.Printf("marketId: %v, description: %v \n", game.MarketID, game.Description)

	err = s.DbConnection.CreateGame(context.TODO(), game)
	if err != nil {
		return errors.Join(errors.New("error inserting into postgres"), err)
	}
	return nil
}

func (s Service) LogoutMatchbook() error {
	response, err := s.MatchbookClient.LogoutMatchbook()
	if err != nil {
		return errors.Join(errors.New("logout error"), err)
	}

	fmt.Printf("logged out of current session %v \n", response)
	return nil
}

func (s Service) Health() error {
	err := s.DbConnection.CheckConnection()
	if err != nil {
		return errors.Join(errors.New("error connecting to db"), err)
	}
	return nil
}

func (s Service) GetMatchbookToken() string {
	return s.MatchbookClient.GetMatchbookToken()
}
