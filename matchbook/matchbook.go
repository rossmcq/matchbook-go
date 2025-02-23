package matchbook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type eventResponse struct {
	Id               int64    `json:"id"`
	Name             string   `json:"name"`
	SportId          int      `json:"sport-id"`
	Start            string   `json:"start"`
	InRunningFlag    bool     `json:"in-running-flag"`
	AllowLiveBetting bool     `json:"allow-live-betting"`
	CategoryId       []int    `json:"catgory-id"`
	Status           string   `json:"status"`
	Volume           float32  `json:"volume"`
	Markets          []market `json:"markets"`
}

type Client struct {
	Token string
}

type sessionResponse struct {
	SessionToken string `json:"session-token"`
}

type market struct {
	Live bool   `json:"live"`
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func New() (Client, error) {
	// TODO: If session active
	url := "https://api.matchbook.com/bpapi/rest/security/session"
	username := os.Getenv("MATCHBOOK_USER")
	password := os.Getenv("MATCHBOOK_PW")

	payload := strings.NewReader("{\"username\":\"" + username + "\",\"password\":\"" + password + "\"}")
	req, _ := http.NewRequest("POST", url, payload)

	addHeaders(req)
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Client{}, fmt.Errorf("failed making http request: %w", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Client{}, fmt.Errorf("failed io.ReadAll:  %w", err)
	}

	var json_body sessionResponse
	err = json.Unmarshal(body, &json_body)
	if err != nil {
		return Client{}, fmt.Errorf("unable to unmarshal matchbook token response: %s", err)
	}

	sessionToken := json_body.SessionToken

	return Client{Token: sessionToken}, nil
}

func (c Client) LogoutMatchbook(token *string) (string, error) {
	url := "https://api.matchbook.com/bpapi/rest/security/session"

	req, _ := http.NewRequest("DELETE", url, nil)
	addHeaders(req)

	req.Header.Add("session-token", c.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed making http request: %w", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body), nil

}

func (c Client) GetMatchOddsMarketId(eventId string) (int64, string, error) {
	get_event_url := "https://api.matchbook.com/edge/rest/events/" + eventId
	req, _ := http.NewRequest("GET", get_event_url, nil)

	addHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("failed making http request: %w", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var json_body eventResponse
	err = json.Unmarshal(body, &json_body)
	if err != nil {
		return -1, "", fmt.Errorf("unable to unmarshal: %v", err)
	}

	markets := json_body.Markets
	if markets == nil {
		return -1, "", fmt.Errorf("unable to parse markets: %v", err)
	}

	for i := 0; i < len(markets); i++ {
		market := markets[i]
		if market.Name == "Match Odds" {
			return market.Id, json_body.Name + " " + json_body.Start, nil
		}
	}

	return -1, "", fmt.Errorf("no match odds found")

}

func addHeaders(req *http.Request) {
	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "api-doc-test-client")
}
