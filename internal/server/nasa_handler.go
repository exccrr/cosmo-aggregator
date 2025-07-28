package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/exccrr/cosmo-aggregator/internal/cache"
	"github.com/exccrr/cosmo-aggregator/internal/nasa"
)

func MarsPhotosHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = "2022-12-01" // дефолтная дата
	}

	cacheKey := "nasa:mars:" + date
	if cached, err := cache.Get(cacheKey); err == nil && r.URL.Query().Get("view") != "html" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	apiKey := os.Getenv("NASA_API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	photos, err := nasa.GetMarsPhotos(apiKey, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 5
	}
	if len(photos) > limit {
		photos = photos[:limit]
	}

	if r.URL.Query().Get("view") == "html" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<h1>Mars Rover Photos (%s)</h1>", date)
		for _, p := range photos {
			fmt.Fprintf(w, `<div style="margin:10px"><img src="%s" width="300"><p>%s (%s)</p></div>`,
				p.ImgSrc, p.Camera.Name, p.EarthDate)
		}
		return
	}

	jsonData, _ := json.Marshal(photos)
	cache.Set(cacheKey, string(jsonData), 10*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
