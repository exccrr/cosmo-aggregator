package server

import (
	"encoding/json"
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
	if cached, err := cache.Get(cacheKey); err == nil {
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

	jsonData, _ := json.Marshal(photos)
	cache.Set(cacheKey, string(jsonData), 10*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
