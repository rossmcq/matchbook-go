package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/rossmcq/matchbook-go/service"
	mock_service "github.com/rossmcq/matchbook-go/service/mocks"
)

func TestService_Success(t *testing.T) {
	t.Log("TestService_Success")

	gomockCtrl := gomock.NewController(t)

	matchbookClient := mock_service.NewMockMatchbookClient(gomockCtrl)
	dbConnection := mock_service.NewMockDbConnection(gomockCtrl)

	s := service.New(matchbookClient, dbConnection)
	assert.NotNil(t, s)
}
