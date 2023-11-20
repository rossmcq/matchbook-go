package model

type Game struct {
	GameID      uint64 `json:"id"`
	EventID     string `json:"event_id"`
	MarketID    string `json:"market_id"`
	Description string `json:"description"`
}
