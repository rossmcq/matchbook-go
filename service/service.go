package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	uuid "github.com/kevinburke/go.uuid"

	matchbook "github.com/rossmcq/matchbook-go/adapter"
	"github.com/rossmcq/matchbook-go/model"
	"github.com/rossmcq/matchbook-go/postgres"
)

type Service struct {
	MatchbookClient matchbook.Client
	DbConnection    *postgres.DbConnection
}

func New(matchbookClient matchbook.Client, dbConnection *postgres.DbConnection) Service {
	return Service{
		MatchbookClient: matchbookClient,
		DbConnection:    dbConnection,
	}
}

func (s Service) CreateEvent(id string) error {
	marketID, description, err := s.MatchbookClient.GetMatchOddsMarketId(id)
	if err != nil {
		return errors.New("error getting market id")
	}
	marketIDStr := strconv.FormatInt(marketID, 10)
	game := model.Game{
		GameID:      uuid.NewV4(),
		EventID:     id,
		MarketID:    marketIDStr,
		Description: description}

	fmt.Printf("marketId: %v, description: %v \n", game.MarketID, game.Description)

	err = s.DbConnection.InsertOrReturnGameID(context.TODO(), game)
	if err != nil {
		return errors.Join(errors.New("error inserting into postgres"), err)
	}
	return nil
}

func (s Service) LogoutMatchbook() error {
	response, err := s.MatchbookClient.LogoutMatchbook(&s.MatchbookClient.Token)
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
