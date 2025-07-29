# CosmoHub

CosmoHub is a Go-based space data aggregator providing real-time ISS tracking, NASA APOD, Mars Rover Photos, and Near-Earth Asteroids information.  
It includes a simple web frontend and REST API with optional HTML views.

---

## Features
- **NASA APOD** – Astronomy Picture of the Day with dynamic date selection.
- **Mars Rover Photos** – View Mars photos by available dates with Redis caching.
- **Near-Earth Asteroids (NEO)** – List upcoming asteroids with size, distance, and hazard flag.
- **ISS Tracker** – Real-time ISS position via WebSocket + live map.
- **SpaceX Launches** – (Available in API, HTML view planned).
- **Redis caching** for API responses.
- **Web frontend** with dynamic date selection and integrated ISS map.

---

## Installation

## 1. Clone the repository:
```bash
git clone https://github.com/yourusername/cosmo-aggregator.git
cd cosmo-aggregator
```
## 2. Install dependencies:
```
go mod tidy
```
## 3. Configure environment:

Create ```.env``` file:
```
NASA_API_KEY=DEMO_KEY
REDIS_ADDR=localhost:6379
```
(You can get a free NASA API key here: https://api.nasa.gov/)

## 4. Run Redis:
```
docker run -d -p 6379:6379 redis
```
## 5. Run the server:
```
go run cmd/main.go
```
Server will start at http://localhost:8080

## API Endpoints:

### _**NASA APOD**:_
- ```JSON:``` /nasa/apod?date=YYYY-MM-DD
- ```HTML:``` /nasa/apod?date=YYYY-MM-DD&view=html```

---

### _**Mars Rover Photos**:_
- ```JSON:``` /nasa/mars/photos?date=YYYY-MM-DD&limit=N
- ```HTML:``` /nasa/mars/photos?date=YYYY-MM-DD&limit=N&view=html
- ```Available dates:``` /nasa/mars/dates

---

### _**Near-Earth Asteroids**:_
- ```JSON:``` /nasa/asteroids?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
- ```HTML:``` /nasa/asteroids?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&view=html

---

### _**SpaceX Launches**:_
- ```JSON:``` /spacex/launches

---

### _**ISS Tracker**:_
- ```WebSocket:``` ws://localhost:8080/ws/iss
- ```Live Map:``` /iss/map

---

## Web Frontend

Visit: http://localhost:8080/

### Features:
- Select APOD date from available list.
- Choose Mars Rover photos by date.
- View Near-Earth Asteroids in HTML table.
- Real-time ISS map integrated via WebSocket.

---

## Technologies
- Go 1.21+
- Redis (for caching)
- NASA Open APIs
- SpaceX API
- Gorilla WebSocket
- Leaflet.js (for map rendering)

---