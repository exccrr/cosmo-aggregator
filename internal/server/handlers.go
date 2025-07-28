package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/exccrr/cosmo-aggregator/internal/spacex"
)

func SpaceXHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 5
	}

	launches, err := spacex.GetLatestLaunches(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(launches)
}
