package pushover

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func SendNotification(title string, body string) bool {
	pushoverUser, ok := os.LookupEnv("PUSHOVER_USER")
	if !ok {
		log.Fatal("Missing value for PUSHOVER_USER environment variable")
	}
	pushoverToken, ok := os.LookupEnv("PUSHOVER_TOKEN")
	if !ok {
		log.Fatal("Missing value for PUSHOVER_TOKEN environment variable")
	}

	urlStr := "https://api.pushover.net/1/messages.json"

	msgData := url.Values{}
	msgData.Set("token", pushoverToken)
	msgData.Set("user", pushoverUser)
	msgData.Set("title", title)
	msgData.Set("message", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	// req.Header.Add("Accept", "application/json")
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Print("failed to create pushover notification:", err)
		return false
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			// https://support.twilio.com/hc/en-us/articles/223134387-What-is-a-Message-SID-
			// The Message SID is the unique ID for any message successfully created by Twilio’s API.
			// It is a 34 character string that starts with “SM…” for text messages and “MM…” for media messages.
			// log.Print("sent text message (sid: " + data["sid"].(string) + ")")
			return true
		}
	} else {
		log.Print("failed (status: " + resp.Status + ")")
	}

	return false
}
