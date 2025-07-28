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

func APODDatesHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 7
	}

	cacheKey := fmt.Sprintf("nasa:apod:dates:%d", limit)
	if cached, err := cache.Get(cacheKey); err == nil && r.URL.Query().Get("view") != "html" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	apiKey := os.Getenv("NASA_API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	var dates []string
	today := time.Now()
	for i := 0; i < limit; i++ {
		date := today.AddDate(0, 0, -i).Format("2006-01-02")
		apod, err := nasa.GetAPOD(apiKey, date)
		if err == nil && apod.Url != "" {
			dates = append(dates, date)
		}
	}

	if r.URL.Query().Get("view") == "html" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<h1>Available APOD Dates</h1>")
		for _, d := range dates {
			fmt.Fprintf(w, `<div><a href="/nasa/apod?date=%s&view=html">%s</a></div>`, d, d)
		}
		return
	}

	jsonData, _ := json.Marshal(dates)
	cache.Set(cacheKey, string(jsonData), 12*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
