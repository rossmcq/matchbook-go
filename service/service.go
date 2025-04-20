//go:generate mockgen -destination=./mocks/mock_matchbook.go -package=mocks github.com/rossmcq/matchbook-go/service MatchbookClient
//go:generate mockgen -destination=./mocks/mock_postgres.go -package=mocks github.com/rossmcq/matchbook-go/service Store

package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rossmcq/matchbook-go/model"
)

type Service struct {
	MatchbookClient MatchbookClient
	Store           Store
}

type MatchbookClient interface {
	GetEvent(eventId string) (model.EventResponse, error)
	LogoutMatchbook() (string, error)
	GetMatchbookToken() string
	GetMarket(ctx context.Context, eventId string, marketId int64) (model.MarketResponse, error)
}

type Store interface {
	CreateGame(ctx context.Context, game model.Game) error
	CheckConnection() error
	GetOpenGames(ctx context.Context) ([]model.Game, error)
	InsertMatchOdds(ctx context.Context, matchOdds model.MatchOdds) error
}

func New(matchbookClient MatchbookClient, dbConnection Store) (Service, error) {
	if matchbookClient == nil {
		return Service{}, errors.New("matchbook client is nil")
	}

	if dbConnection == nil {
		return Service{}, errors.New("database client is nil")
	}

	return Service{
		MatchbookClient: matchbookClient,
		Store:           dbConnection,
	}, nil
}

func (s Service) CreateEvent(ctx context.Context, id string) error {
	event, err := s.MatchbookClient.GetEvent(id)
	if err != nil {
		return errors.New("error getting market id")
	}

	// Find the Match Odds market for the event
	markets := event.Markets
	if markets == nil {
		return fmt.Errorf("no markets found for event %s", id)
	}

	var matchOddsMarket model.Market
	for _, market := range markets {
		if market.Name == "Match Odds" {
			matchOddsMarket = market
			break
		}
	}

	if matchOddsMarket.Id == 0 {
		return errors.New("no match odds found")
	}

	gameStart, err := time.Parse(time.RFC3339, event.Start)
	if err != nil {
		return errors.Join(errors.New("error parsing time"), err)
	}

	game := model.Game{
		GameID:      id,
		EventID:     id,
		MarketID:    matchOddsMarket.Id,
		StartAt:     gameStart,
		Status:      event.Status,
		HomeTeam:    strings.Split(event.Name, " vs ")[0],
		AwayTeam:    strings.Split(event.Name, " vs ")[1],
		Description: event.Name,
	}

	log.Printf("marketId: %v, description: %v \n", game.MarketID, game.Description)

	err = s.Store.CreateGame(ctx, game)
	if err != nil {
		return errors.Join(errors.New("error inserting into postgres"), err)
	}
	return nil
}

func (s Service) RecordMatchOdds(ctx context.Context) (any, error) {
	log.Printf("recording match odds for game...")
	games, err := s.Store.GetOpenGames(ctx)
	if err != nil {
		return nil, errors.Join(errors.New("error getting open games"), err)
	}

	for _, game := range games {
		Market, err := s.MatchbookClient.GetMarket(ctx, game.EventID, game.MarketID)
		if err != nil {
			return nil, errors.Join(errors.New("error getting market"), err)
		}
		matchOdds := marketResponseToMatchOdds(game, Market)
		err = s.Store.InsertMatchOdds(ctx, matchOdds)
		if err != nil {
			return nil, errors.Join(errors.New("error inserting match odds"), err)
		}
	}

	return "", nil

}

func (s Service) LogoutMatchbook() error {
	response, err := s.MatchbookClient.LogoutMatchbook()
	if err != nil {
		return errors.Join(errors.New("logout error"), err)
	}

	log.Printf("logged out of current session %v \n", response)
	return nil
}

func (s Service) Health() error {
	err := s.Store.CheckConnection()
	if err != nil {
		return errors.Join(errors.New("error connecting to db"), err)
	}
	return nil
}

func (s Service) GetMatchbookToken() string {
	return s.MatchbookClient.GetMatchbookToken()
}

func marketResponseToMatchOdds(game model.Game, marketResponse model.MarketResponse) model.MatchOdds {
	matchOdds := model.MatchOdds{}

	log.Printf("marketResponse: %v", marketResponse)

	// for _, runner := range marketResponse.Runners {
	// 	runnerIsHome := runner.Name == game.HomeTeam
	// 	bestBackOdds := 0.0
	// 	for _, price := range runner.Prices {
	// 		matchOdds.HomeWinBackOdds = price.Odds
	// 	}
	// }

	return matchOdds
}
