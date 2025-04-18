package model

import "time"

type Game struct {
	GameID      string
	EventID     string
	MarketID    int64
	StartAt     time.Time `json:"start"`
	Status      string    `json:"status"`
	HomeTeam    string
	AwayTeam    string
	Description string `json:"name"`
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
