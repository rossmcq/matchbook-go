package application

import (
	"fmt"

	"github.com/rossmcq/matchbook-go/matchbook"
)

type Config struct {
	matchbookToken string
}

func LoadConfig() Config {
	matchbookToken, err := matchbook.LoadMatchboookToken()
	if err != nil {
		fmt.Printf("Error loading matchbook token %s", err)
	}
	return Config{
		matchbookToken: *matchbookToken}
}
