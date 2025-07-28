package nasa

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const apodURL = "https://api.nasa.gov/planetary/apod"

type APOD struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Url         string `json:"url"`
	MediaType   string `json:"media_type"`
}

func GetAPOD(apiKey, date string) (*APOD, error) {
	url := fmt.Sprintf("%s?api_key=%s", apodURL, apiKey)
	if date != "" {
		url += "&date=" + date
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apod APOD
	if err := json.NewDecoder(resp.Body).Decode(&apod); err != nil {
		return nil, err
	}

	return &apod, nil
}
