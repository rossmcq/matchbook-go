package model

import (
	uuid "github.com/kevinburke/go.uuid"
)

type Game struct {
	GameID      uuid.UUID
	EventID     string
	MarketID    string
	Description string
}
