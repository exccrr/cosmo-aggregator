package server

import "net/http"

func FrontendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>CosmoHub</title>
        <meta charset="utf-8"/>
        <style>
            body { font-family: Arial, sans-serif; margin: 20px; background: #0b0c10; color: #fff; }
            button { margin: 5px; padding: 10px 15px; background: #1f2833; color: #fff; border: none; cursor: pointer; }
            button:hover { background: #45a29e; }
            #content { margin-top: 20px; padding: 10px; background: #1f2833; border-radius: 8px; }
            img { max-width: 100%; border-radius: 6px; }
            input, select { padding: 5px; margin: 5px; }
        </style>
    </head>
    <body>
        <h1>🚀 CosmoHub</h1>
        <div>
            <button onclick="showAPOD()">NASA APOD</button>
            <button onclick="showMars()">Mars Photos</button>
            <button onclick="showAsteroids()">Asteroids</button>
            <button onclick="showISS()">ISS Map</button>
        </div>
        <div id="content">Select a category to view data...</div>

        <script>
            function showAPOD() {
                fetch('/nasa/apod/dates')
                    .then(r => r.json())
                    .then(dates => {
                        let options = dates.map(d => "<option value='" + d + "'>" + d + "</option>").join("");
                        document.getElementById('content').innerHTML =
                            "<h2>NASA APOD</h2>" +
                            "<select id='apodDate'>" + options + "</select>" +
                            "<button onclick='fetchAPOD()'>Load</button>" +
                            "<div id='apodResult'></div>";
                    });
            }

            function fetchAPOD() {
                let date = document.getElementById('apodDate').value;
                let url = '/nasa/apod?view=html';
                if (date) url += '&date=' + date;
                fetch(url)
                    .then(r => r.text())
                    .then(html => document.getElementById('apodResult').innerHTML = html);
            }

            function showMars() {
                fetch('/nasa/mars/dates')
                    .then(r => r.json())
                    .then(dates => {
                        let options = dates.map(d => "<option value='" + d + "'>" + d + "</option>").join("");
                        document.getElementById('content').innerHTML =
                            "<h2>Mars Photos</h2>" +
                            "<select id='marsDate'>" + options + "</select>" +
                            "<input type='number' id='marsLimit' placeholder='Limit (default 5)'>" +
                            "<button onclick='fetchMars()'>Load</button>" +
                            "<div id='marsResult'></div>";
                    });
            }

            function fetchMars() {
                let date = document.getElementById('marsDate').value;
                let limit = document.getElementById('marsLimit').value;
                let url = '/nasa/mars/photos?view=html';
                if (date) url += '&date=' + date;
                if (limit) url += '&limit=' + limit;
                fetch(url)
                    .then(r => r.text())
                    .then(html => document.getElementById('marsResult').innerHTML = html);
            }

            function showAsteroids() {
                document.getElementById('content').innerHTML =
                    "<h2>Near-Earth Asteroids</h2>" +
                    "Start: <input type='date' id='startDate'>" +
                    "End: <input type='date' id='endDate'>" +
                    "<button onclick='fetchAsteroids()'>Load</button>" +
                    "<div id='asteroidsResult'></div>";
            }

            function fetchAsteroids() {
                let start = document.getElementById('startDate').value;
                let end = document.getElementById('endDate').value;
                let url = '/nasa/asteroids?view=html';
                if (start) url += '&start_date=' + start;
                if (end) url += '&end_date=' + end;
                fetch(url)
                    .then(r => r.text())
                    .then(html => document.getElementById('asteroidsResult').innerHTML = html);
            }

            function showISS() {
                document.getElementById('content').innerHTML =
                    "<iframe src='/iss/map' width='100%' height='500px' style='border:none;'></iframe>";
            }
        </script>
    </body>
    </html>
    `
	w.Write([]byte(html))
}
