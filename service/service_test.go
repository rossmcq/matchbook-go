package service_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/rossmcq/matchbook-go/model"
	"github.com/rossmcq/matchbook-go/service"
	mock_service "github.com/rossmcq/matchbook-go/service/mocks"
)

func TestService_New_Success(t *testing.T) {
	t.Log("TestService_Success")

	gomockCtrl := gomock.NewController(t)

	matchbookClient := mock_service.NewMockMatchbookClient(gomockCtrl)
	dbConnection := mock_service.NewMockDbConnection(gomockCtrl)

	s, err := service.New(matchbookClient, dbConnection)
	assert.NotNil(t, s)
	assert.NoError(t, err)
}

func TestService_New_Failure(t *testing.T) {
	t.Log("TestService_Failure")

	gomockCtrl := gomock.NewController(t)

	matchbookClient := mock_service.NewMockMatchbookClient(gomockCtrl)
	dbConnection := mock_service.NewMockDbConnection(gomockCtrl)

	_, err := service.New(nil, dbConnection)
	assert.EqualError(t, err, "matchbook client is nil")

	_, err = service.New(matchbookClient, nil)
	assert.EqualError(t, err, "database client is nil")
}

func TestService_CreateEvent_Sucess(t *testing.T) {
	t.Log("TestService_CreateEvent_Success")

	eventID := "123"
	marketID := int64(55)
	description := "Team A vs Team B"

	gomockCtrl := gomock.NewController(t)

	matchbookClient := mock_service.NewMockMatchbookClient(gomockCtrl)
	dbConnection := mock_service.NewMockDbConnection(gomockCtrl)

	s, err := service.New(matchbookClient, dbConnection)
	assert.NotNil(t, s)
	assert.NoError(t, err)

	matchbookClient.EXPECT().GetEvent(eventID).Return(model.EventResponse{
		Name:    description,
		Start:   "2025-04-18T14:00:00.000Z",
		Status:  "open",
		Markets: []model.Market{{Id: marketID, Name: "Match Odds"}},
	}, nil)
	dbConnection.EXPECT().CreateGame(gomock.Any(), model.Game{
		GameID:      eventID,
		EventID:     eventID,
		MarketID:    marketID,
		StartAt:     time.Date(2025, 4, 18, 14, 0, 0, 0, time.UTC),
		Status:      "open",
		HomeTeam:    "Team A",
		AwayTeam:    "Team B",
		Description: description}).Return(nil)

	err = s.CreateEvent(eventID)
	assert.NoError(t, err)
}
