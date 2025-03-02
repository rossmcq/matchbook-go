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

type EventResponse struct {
	Id               int64    `json:"id"`
	Name             string   `json:"name"`
	SportId          int      `json:"sport-id"`
	Start            string   `json:"start"`
	InRunningFlag    bool     `json:"in-running-flag"`
	AllowLiveBetting bool     `json:"allow-live-betting"`
	CategoryId       []int    `json:"catgory-id"`
	Status           string   `json:"status"`
	Volume           float32  `json:"volume"`
	Markets          []Market `json:"markets"`
}

type Market struct {
	Live bool   `json:"live"`
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
