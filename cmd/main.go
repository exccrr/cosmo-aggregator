package main

import (
	"log"
	"net/http"

	"github.com/exccrr/cosmo-aggregator/internal/cache"
	"github.com/exccrr/cosmo-aggregator/internal/server"
)

func main() {
	cache.InitRedis("localhost:6379") // Redis

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/spacex/launches", server.SpaceXHandler)
	mux.HandleFunc("/nasa/mars/photos", server.MarsPhotosHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
