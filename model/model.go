package model

import "time"

type Game struct {
	GameID      string
	EventID     string
	MarketID    int64
	StartAt     time.Time
	Status      string
	HomeTeam    string
	AwayTeam    string
	Description string
}

type MatchOdds struct {
	GameID                  string
	PriceUpdateAt           time.Time
	HomeWinBackOdds         float32
	HomeWinBackAmount       float32
	HomeWinLayOdds          float32
	HomeWinLayAmount        float32
	DrawBackOdds            float32
	DrawBackAmount          float32
	DrawLayOdds             float32
	DrawLayAmount           float32
	AwayWinBackOdds         float32
	AwayWinBackAmount       float32
	AwayWinLayOdds          float32
	AwayWinLayAmount        float32
	HomeWinBackOddsSecond   float32
	HomeWinBackAmountSecond float32
	HomeWinLayOddsSecond    float32
	HomeWinLayAmountSecond  float32
	DrawBackOddsSecond      float32
	DrawBackAmountSecond    float32
	DrawLayOddsSecond       float32
	DrawLayAmountSecond     float32
	AwayWinBackOddsSecond   float32
	AwayWinBackAmountSecond float32
	AwayWinLayOddsSecond    float32
	AwayWinLayAmountSecond  float32
}

type EventResponse struct {
	Id               string   `json:"id"`
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

type MarketResponse struct {
	Id      int64    `json:"id"`
	Name    string   `json:"name"`
	Status  string   `json:"status"`
	Runners []Runner `json:"runners"`
	Errors  []Error  `json:"errors"`
}

type Error struct {
	Messages []string `json:"messages"`
}

type Runner struct {
	Id               int64   `json:"id"`
	Name             string  `json:"name"`
	LastPriceUpateAt string  `json:"last-price-update-time"`
	Prices           []Price `json:"prices"`
}

type Price struct {
	Odds            float32 `json:"odds"`
	Side            string  `json:"side"`
	AvailableAmount float32 `json:"available-amount"`
}

type Market struct {
	Live bool   `json:"live"`
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
