package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/exccrr/cosmo-aggregator/internal/cache"
	"github.com/exccrr/cosmo-aggregator/internal/nasa"
)

func APODHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	cacheKey := "nasa:apod:" + date
	if cached, err := cache.Get(cacheKey); err == nil && r.URL.Query().Get("view") != "html" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	apiKey := os.Getenv("NASA_API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	apod, err := nasa.GetAPOD(apiKey, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if apod.Url == "" {
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		apod, err = nasa.GetAPOD(apiKey, yesterday)
		if err != nil || apod.Url == "" {
			http.Error(w, "NASA APOD is not available for this date yet", http.StatusNotFound)
			return
		}
		date = yesterday
		cacheKey = "nasa:apod:" + date
	}

	if r.URL.Query().Get("view") == "html" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<h1>%s (%s)</h1>", apod.Title, apod.Date)
		if apod.MediaType == "image" {
			fmt.Fprintf(w, `<img src="%s" width="600"><p>%s</p>`, apod.Url, apod.Explanation)
		} else if apod.MediaType == "video" {
			fmt.Fprintf(w, `<iframe src="%s" width="800" height="450" frameborder="0" allowfullscreen></iframe><p>%s</p>`, apod.Url, apod.Explanation)
		} else {
			fmt.Fprintf(w, `<p><a href="%s" target="_blank">Open Media</a></p><p>%s</p>`, apod.Url, apod.Explanation)
		}
		return
	}

	jsonData, _ := json.Marshal(apod)
	cache.Set(cacheKey, string(jsonData), 6*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
