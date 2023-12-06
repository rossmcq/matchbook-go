package matchbook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func LoadMatchboookToken() (*string, error) {
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
		return nil, fmt.Errorf("failed making http request: %w", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		var x *string
		return x, fmt.Errorf("failed io.ReadAll:  %w", err)

	}
	var json_body map[string]string
	json.Unmarshal(body, &json_body)
	sessionToken := json_body["session-token"]

	return &sessionToken, nil
}

func LogoutMatchbook(token *string) (string, error) {
	url := "https://api.matchbook.com/bpapi/rest/security/session"

	req, _ := http.NewRequest("DELETE", url, nil)
	addHeaders(req)

	req.Header.Add("session-token", *token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed making http request: %w", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body), nil

}

func GetMatchOddsMarketId(eventId string) (float64, string, error) {
	get_event_url := "https://api.matchbook.com/edge/rest/events/" + eventId
	fmt.Println("url: ", get_event_url)
	req, _ := http.NewRequest("GET", get_event_url, nil)

	addHeaders(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("failed making http request: %w", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// TODO Refactor this mess and pass into a struct
	var json_body interface{}
	err = json.Unmarshal(body, &json_body)
	if err != nil {
		return -1, "", fmt.Errorf("Unable to unmarshal: %v", err)
	}

	m := json_body.(map[string]interface{})

	markets := m["markets"]
	if markets == nil {
		return -1, "", fmt.Errorf("Unable to unmarshal: %v", err)
	}
	v := markets.([]interface{})

	for i := 0; i < len(v); i++ {
		x := v[i].(map[string]interface{})
		if x["name"] == "Match Odds" {
			return x["id"].(float64), m["name"].(string), nil
		}
	}

	return -1, "", fmt.Errorf("No match odds found")

}

func addHeaders(req *http.Request) {
	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "api-doc-test-client")
}
