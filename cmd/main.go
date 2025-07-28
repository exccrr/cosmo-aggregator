package main

import (
	"log"
	"net/http"
	"os"

	"github.com/exccrr/cosmo-aggregator/internal/cache"
	"github.com/exccrr/cosmo-aggregator/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	cache.InitRedis(redisAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/spacex/launches", server.SpaceXHandler)
	mux.HandleFunc("/nasa/mars/photos", server.MarsPhotosHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
