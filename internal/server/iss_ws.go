package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var lastISS *ISSLocation

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

	if lastISS != nil {
		data, _ := json.Marshal(lastISS)
		conn.WriteMessage(websocket.TextMessage, data)
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
	conn.Close()
	log.Println("ISS WebSocket client disconnected")
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
	resp, err := http.Get("https://api.wheretheiss.at/v1/satellites/25544")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResp struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timestamp int64   `json:"timestamp"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}
	log.Printf("ISS API response: %s", body)

	return &ISSLocation{
		Latitude:  fmt.Sprintf("%.6f", apiResp.Latitude),
		Longitude: fmt.Sprintf("%.6f", apiResp.Longitude),
		Timestamp: apiResp.Timestamp,
	}, nil

}

func broadcastISS(loc *ISSLocation) {
	lastISS = loc

	data, _ := json.Marshal(loc)
	clientsMu.Lock()
	defer clientsMu.Unlock()

	log.Printf("Sending ISS coords to %d clients", len(clients))
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
            body { font-family: Arial, sans-serif; background: #0b0c10; color: #fff; }
            #map { height: 80vh; width: 100%; margin-bottom: 10px; }
            #info { padding: 10px; background: #1f2833; border-radius: 6px; }
        </style>
    </head>
    <body>
        <h1>ISS Location in Real Time</h1>
        <div id="map"></div>
        <div id="info">Waiting for coordinates...</div>

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
            var firstUpdate = true;

			const ws = new WebSocket("ws://localhost:8080/ws/iss");

			ws.onmessage = function(event) {
				var data = JSON.parse(event.data);
				console.log("ISS data:", data);

				var lat = parseFloat(data.latitude);
				var lon = parseFloat(data.longitude);

				marker.setLatLng([lat, lon]);

				if (firstUpdate) {
					map.setView([lat, lon], 4);
					firstUpdate = false;
				}

				// map.setView([lat, lon], 4);

				var ts = new Date(data.timestamp * 1000).toLocaleTimeString();
				document.getElementById("info").innerText =
					"Latitude: " + lat.toFixed(4) +
					" | Longitude: " + lon.toFixed(4) +
					" | Time: " + ts;
			};

        </script>
    </body>
    </html>
    `
	w.Write([]byte(html))
}
