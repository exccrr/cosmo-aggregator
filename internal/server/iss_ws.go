package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type issAPIResponse struct {
	Timestamp   int64 `json:"timestamp"`
	ISSPosition struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"iss_position"`
}

type ISSLocation struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Timestamp int64  `json:"timestamp"`
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func ISSWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	log.Println("New ISS WebSocket client connected")

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
		log.Println("ISS WebSocket client disconnected")
	}()
}

func StartISSUpdater() {
	go func() {

		loc, err := fetchISSLocation()
		if err == nil {
			broadcastISS(loc)
		}

		for {
			time.Sleep(10 * time.Second)
			loc, err := fetchISSLocation()
			if err != nil {
				log.Println("ISS fetch error:", err)
			} else {
				broadcastISS(loc)
			}
		}
	}()
}

func fetchISSLocation() (*ISSLocation, error) {
	resp, err := http.Get("http://api.open-notify.org/iss-now.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResp issAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	return &ISSLocation{
		Latitude:  apiResp.ISSPosition.Latitude,
		Longitude: apiResp.ISSPosition.Longitude,
		Timestamp: apiResp.Timestamp,
	}, nil
}

func broadcastISS(loc *ISSLocation) {
	data, _ := json.Marshal(loc)

	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("WebSocket write error:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

func ISSMapHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>ISS Tracker</title>
        <meta charset="utf-8" />
        <link rel="stylesheet" href="https://unpkg.com/leaflet/dist/leaflet.css"/>
        <style>
            #map { height: 90vh; width: 100%; }
        </style>
    </head>
    <body>
        <h1>ISS Location in Real Time</h1>
        <div id="map"></div>

        <script src="https://unpkg.com/leaflet/dist/leaflet.js"></script>
        <script>
            var map = L.map('map').setView([0, 0], 2);
            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                maxZoom: 19
            }).addTo(map);

            var issIcon = L.icon({
                iconUrl: 'https://upload.wikimedia.org/wikipedia/commons/d/d0/International_Space_Station.svg',
                iconSize: [50, 32]
            });

            var marker = L.marker([0, 0], {icon: issIcon}).addTo(map);

            const ws = new WebSocket("ws://localhost:8080/ws/iss");
            ws.onmessage = function(event) {
                var data = JSON.parse(event.data);
                var lat = parseFloat(data.latitude);
                var lon = parseFloat(data.longitude);
                marker.setLatLng([lat, lon]);
                map.setView([lat, lon], map.getZoom());
            };
        </script>
    </body>
    </html>
    `
	w.Write([]byte(html))
}
