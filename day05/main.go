package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Satellite struct {
	ID          string `json:"id"`
	Temperature int    `json:"temperature"`
	Battery     int    `json:"battery"`
}

type HistoryPoint struct {
	Time string `json:"time"`
	Temp int    `json:"temp"`
}

var loc = time.FixedZone("BRT", -3*60*60)

func (s Satellite) CheckHealth() string {
	if s.Temperature > 85 || s.Battery < 15 {
		return "CRITICAL"
	}
	if s.Temperature > 75 || s.Battery < 30 {
		return "WARNING"
	}
	return "NOMINAL"
}

type Telemetry struct {
	Satellite Satellite `json:"satellite"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	satellites = map[string]*Satellite{
		"SAT-01": {ID: "SAT-01", Temperature: 70, Battery: 85},
		"SAT-02": {ID: "SAT-02", Temperature: 65, Battery: 10},
		"SAT-03": {ID: "SAT-03", Temperature: 80, Battery: 50},
	}
	history    = map[string][]HistoryPoint{}
	mu         sync.Mutex
	startTime  = time.Now()
	order      = []string{"SAT-01", "SAT-02", "SAT-03"}
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}
		next(w, r)
	}
}

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`
	<html><head><title>Mission Control API</title></head>
	<body style="font-family:monospace;padding:40px;background:#0a0a0a;color:#00ff88;max-width:800px">
	<h1>🚀 Mission Control Go — Day 07</h1>
	<p>Base: https://%s</p>
	<hr>
	<h2>GET /health</h2><p>Uptime</p>
	<h2>GET /satellites</h2><p>Lista frota</p>
	<h2>GET /telemetry/{id}</h2><p>Ex: /telemetry/SAT-02</p>
	<h2>POST /satellites/{id}/heat</h2><p>Aquece +5C</p>
	<h2>GET /history/{id}</h2><p>Histórico 20 pontos em BRT</p>
	<hr><p>by Jairo Andrade — <a style="color:#00ff88" href="https://github.com/Jairocpdev/mission-control-go">【entity-GitHub¦canonical_name=GitHub】</a></p>
	</body></html>`, r.Host)
	w.Write([]byte(html))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"uptime": fmt.Sprintf("%.0f segundos", time.Since(startTime).Seconds()),
	})
}

func satellitesHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	type Resp struct {
		Satellite
		Status string `json:"status"`
	}
	list := []Resp{}
	for _, id := range order {
		s := satellites[id]
		list = append(list, Resp{Satellite: *s, Status: s.CheckHealth()})
	}
	json.NewEncoder(w).Encode(list)
}

func telemetryHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/telemetry/")
	mu.Lock()
	defer mu.Unlock()
	sat, ok := satellites[id]
	if !ok {
		http.Error(w, "satellite not found", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Telemetry{Satellite: *sat, Status: sat.CheckHealth(), Timestamp: time.Now()})
}

func heatHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "satellites" || parts[2] != "heat" {
		http.NotFound(w, r)
		return
	}
	id := parts[1]
	mu.Lock()
	defer mu.Unlock()
	sat, ok := satellites[id]
	if !ok {
		http.Error(w, "satellite not found", 404)
		return
	}
	sat.Temperature += 5
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"satellite": sat,
		"status":    sat.CheckHealth(),
		"message":   fmt.Sprintf("%s aquecido para %d°C", id, sat.Temperature),
	})
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/history/")
	mu.Lock()
	defer mu.Unlock()
	h, ok := history[id]
	if !ok {
		h = []HistoryPoint{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h)
}

func main() {
	// inicia histórico
	for _, id := range order {
		history[id] = []HistoryPoint{{Time: time.Now().In(loc).Format("15:04:05"), Temp: satellites[id].Temperature}}
	}

	http.HandleFunc("/", swaggerHandler)
	http.HandleFunc("/health", corsMiddleware(healthHandler))
	http.HandleFunc("/satellites", corsMiddleware(satellitesHandler))
	http.HandleFunc("/telemetry/", corsMiddleware(telemetryHandler))
	http.HandleFunc("/satellites/", corsMiddleware(heatHandler))
	http.HandleFunc("/history/", corsMiddleware(historyHandler))

	go func() {
		for {
			time.Sleep(10 * time.Second)
			mu.Lock()
			for _, id := range order {
				sat := satellites[id]
				sat.Temperature += 1
				sat.Battery -= 1
				if sat.Temperature > 90 {
					sat.Temperature = 60
					sat.Battery = 85
				}
				if sat.Battery < 0 {
					sat.Battery = 0
				}
				history[id] = append(history[id], HistoryPoint{
					Time: time.Now().In(loc).Format("15:04:05"),
					Temp: sat.Temperature,
				})
				if len(history[id]) > 20 {
					history[id] = history[id][1:]
				}
			}
			mu.Unlock()
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("📡 Mission Control Go Day07 rodando em :%s BRT\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}