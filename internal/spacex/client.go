package spacex

import (
	"encoding/json"
	"net/http"
	"time"
)

const baseURL = "https://api.spacexdata.com/v4/launches"

type Launch struct {
	Name    string    `json:"name"`
	DateUtc time.Time `json:"date_utc"`
	Success bool      `json:"success"`
	Details string    `json:"details"`
}

func GetLatestLaunches(limit int) ([]Launch, error) {
	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var launches []Launch
	if err := json.NewDecoder(resp.Body).Decode(&launches); err != nil {
		return nil, err
	}

	if len(launches) > limit {
		launches = launches[len(launches)-limit:]
	}
	return launches, nil
}
