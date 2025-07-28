package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/exccrr/cosmo-aggregator/internal/cache"
	"github.com/exccrr/cosmo-aggregator/internal/nasa"
)

func MarsDatesHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 7
	}

	cacheKey := fmt.Sprintf("nasa:dates:%d", limit)
	if cached, err := cache.Get(cacheKey); err == nil && r.URL.Query().Get("view") != "html" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	apiKey := os.Getenv("NASA_API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	// ласт 30 дней
	dates := make(map[string]bool)
	today := time.Now()
	for i := 0; i < 30; i++ {
		date := today.AddDate(0, 0, -i).Format("2006-01-02")
		photos, err := nasa.GetMarsPhotos(apiKey, date)
		if err == nil && len(photos) > 0 {
			dates[date] = true
		}
		if len(dates) >= limit {
			break
		}
	}

	var result []string
	for d := range dates {
		result = append(result, d)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(result)))

	if r.URL.Query().Get("view") == "html" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<h1>Available Mars Rover Dates</h1>")
		for _, d := range result {
			fmt.Fprintf(w, `<div><a href="/nasa/mars/photos?date=%s&view=html">%s</a></div>`, d, d)
		}
		return
	}

	jsonData, _ := json.Marshal(result)
	cache.Set(cacheKey, string(jsonData), 1*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
