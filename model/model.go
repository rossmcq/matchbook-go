package model

import "github.com/google/uuid"

type Game struct {
	GameID   uint64
	EventID  uuid.UUID
	MarketID uuid.UUID
}
