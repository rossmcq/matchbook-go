package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

func loadMatchboookToken() string {

	url := "https://api.matchbook.com/bpapi/rest/security/session"
	username := os.Getenv("MATCHBOOK_USER")
	password := os.Getenv("MATCHBOOK_PW")

	payload := strings.NewReader("{\"username\":\"" + username + "\",\"password\":\"" + password + "\"}")
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "api-doc-test-client")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	// if err != nil {
	// 	return "DUMMAY" fmt.Errorf("failed io.ReadAll  %w", err)

	// }
	var json_body map[string]string
	json.Unmarshal(body, &json_body)

	return json_body["session-token"]
}
