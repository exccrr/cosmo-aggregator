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

func AsteroidsHandler(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start_date")
	end := r.URL.Query().Get("end_date")
	if start == "" {
		start = time.Now().Format("2006-01-02")
	}
	if end == "" {
		end = time.Now().AddDate(0, 0, 2).Format("2006-01-02")
	}

	cacheKey := fmt.Sprintf("nasa:asteroids:%s:%s", start, end)
	if cached, err := cache.Get(cacheKey); err == nil && r.URL.Query().Get("view") != "html" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	apiKey := os.Getenv("NASA_API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	asteroids, err := nasa.GetAsteroids(apiKey, start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("view") == "html" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<h1>Near-Earth Objects (%s - %s)</h1>", start, end)
		fmt.Fprintf(w, "<table border='1' cellpadding='5'><tr><th>Name</th><th>Diameter (km)</th><th>Hazardous</th><th>Close Date</th><th>Distance (km)</th><th>Speed (km/s)</th></tr>")
		for _, neo := range asteroids {
			if len(neo.CloseApproachData) == 0 {
				continue
			}
			ca := neo.CloseApproachData[0]
			fmt.Fprintf(w,
				"<tr><td>%s</td><td>%.3f - %.3f</td><td>%t</td><td>%s</td><td>%s</td><td>%s</td></tr>",
				neo.Name,
				neo.EstimatedDiameter.Kilometers.Min,
				neo.EstimatedDiameter.Kilometers.Max,
				neo.IsHazardous,
				ca.CloseApproachDate,
				ca.MissDistance.Kilometers,
				ca.RelativeVelocity.KmPerSecond,
			)
		}
		fmt.Fprintf(w, "</table>")
		return
	}

	jsonData, _ := json.Marshal(asteroids)
	cache.Set(cacheKey, string(jsonData), 1*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
