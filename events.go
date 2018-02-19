package dexcomClient

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

var EVENTS_URL = "v1/users/self/events"



func GetEvents(token string, c *Config) []Event {
	req, _ := http.NewRequest("GET",
		UrlWithDateRange(c, EVENTS_URL, "2015-09-19T00:00:00", "2015-11-10T00:00:00"), nil)
	req.Header.Add("authorization", "Bearer " + token)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var response EventResponse
	json.Unmarshal(body, &response)

	return response.Events
}