package nasa

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos" // фото марсоходов

type Photo struct {
	ID        int    `json:"id"`
	ImgSrc    string `json:"img_src"`
	EarthDate string `json:"earth_date"`
	Camera    struct {
		Name string `json:"name"`
	} `json:"camera"`
}

type photosResponse struct {
	Photos []Photo `json:"photos"`
}

func GetMarsPhotos(apiKey, date string) ([]Photo, error) {
	url := fmt.Sprintf("%s?earth_date=%s&api_key=%s", baseURL, date, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data photosResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Photos, nil
}
