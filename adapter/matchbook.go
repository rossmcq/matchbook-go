package matchbook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rossmcq/matchbook-go/model"
)

type Client struct {
	Token string
}

type sessionResponse struct {
	SessionToken string `json:"session-token"`
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

	var sessionResponse sessionResponse
	err = json.Unmarshal(body, &sessionResponse)
	if err != nil {
		return Client{}, fmt.Errorf("unable to unmarshal matchbook token response: %s", err)
	}

	sessionToken := sessionResponse.SessionToken

	return Client{Token: sessionToken}, nil
}

func (c Client) LogoutMatchbook() (string, error) {
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

func (c Client) GetEvent(eventId string) (model.EventResponse, error) {
	getEventURL := "https://api.matchbook.com/edge/rest/events/" + eventId
	req, _ := http.NewRequest("GET", getEventURL, nil)

	addHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.EventResponse{}, fmt.Errorf("failed making http request: %w", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var eventResponse model.EventResponse
	err = json.Unmarshal(body, &eventResponse)
	if err != nil {
		return eventResponse, fmt.Errorf("unable to unmarshal: %v", err)
	}

	return eventResponse, nil
}

func (c Client) GetMarket(ctx context.Context, eventId string, marketId int64) (model.MarketResponse, error) {
	getMarketURL := fmt.Sprintf("https://api.matchbook.com/edge/rest/events/%s/markets/%d", eventId, marketId)
	req, _ := http.NewRequestWithContext(ctx, "GET", getMarketURL, nil)

	addHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.MarketResponse{}, fmt.Errorf("failed making http request: %w", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var marketResponse model.MarketResponse
	err = json.Unmarshal(body, &marketResponse)
	if err != nil {
		return marketResponse, fmt.Errorf("unable to unmarshal: %v", err)
	}

	return marketResponse, nil
}

func addHeaders(req *http.Request) {
	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "api-doc-test-client")
}

func (c Client) GetMatchbookToken() string {
	return c.Token
}
