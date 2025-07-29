package nasa

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const neoFeedURL = "https://api.nasa.gov/neo/rest/v1/feed"

type NEOFeed struct {
	NearEarthObjects map[string][]NEO `json:"near_earth_objects"`
}

type NEO struct {
	Name              string `json:"name"`
	EstimatedDiameter struct {
		Kilometers struct {
			Min float64 `json:"estimated_diameter_min"`
			Max float64 `json:"estimated_diameter_max"`
		} `json:"kilometers"`
	} `json:"estimated_diameter"`
	IsHazardous       bool `json:"is_potentially_hazardous_asteroid"`
	CloseApproachData []struct {
		CloseApproachDate string `json:"close_approach_date"`
		MissDistance      struct {
			Kilometers string `json:"kilometers"`
		} `json:"miss_distance"`
		RelativeVelocity struct {
			KmPerSecond string `json:"kilometers_per_second"`
		} `json:"relative_velocity"`
	} `json:"close_approach_data"`
}

func GetAsteroids(apiKey, startDate, endDate string) ([]NEO, error) {
	url := fmt.Sprintf("%s?start_date=%s&end_date=%s&api_key=%s", neoFeedURL, startDate, endDate, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var feed NEOFeed
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	var result []NEO
	for _, list := range feed.NearEarthObjects {
		result = append(result, list...)
	}

	return result, nil
}
