package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func loadMatchboookToken() (string, error) {
	url := "https://api.matchbook.com/bpapi/rest/security/session"
	username := os.Getenv("MATCHBOOK_USER")
	password := os.Getenv("MATCHBOOK_PW")

	payload := strings.NewReader("{\"username\":\"" + username + "\",\"password\":\"" + password + "\"}")
	req, _ := http.NewRequest("POST", url, payload)

	addHeaders(req)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed io.ReadAll  %w", err)

	}
	var json_body map[string]string
	json.Unmarshal(body, &json_body)

	return json_body["session-token"], nil
}

func logoutMatchbook(token *string) string {
	url := "https://api.matchbook.com/bpapi/rest/security/session"

	req, _ := http.NewRequest("DELETE", url, nil)
	addHeaders(req)

	req.Header.Add("session-token", *token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body)

}

func getMatchOddsMarketId(eventId string) (float64, string, error) {
	get_event_url := "https://api.matchbook.com/edge/rest/events/" + eventId
	fmt.Println("url: ", get_event_url)
	req, _ := http.NewRequest("GET", get_event_url, nil)

	addHeaders(req)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var json_body interface{}
	err := json.Unmarshal(body, &json_body)
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
