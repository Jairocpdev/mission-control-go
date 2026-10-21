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
	mu        sync.Mutex
	startTime = time.Now()
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
	<h1>🚀 Mission Control Go — Day 05</h1>
	<p>Base: https://%s</p>
	<hr>
	<h2>GET /health</h2><p>Uptime da missão</p><code>curl /health</code>
	<h2>GET /satellites</h2><p>Lista frota com health</p>
	<h2>GET /telemetry/{id}</h2><p>Ex: /telemetry/SAT-02</p>
	<h2>POST /satellites/{id}/heat</h2><p>Aquece +5°C</p><code>curl -X POST /satellites/SAT-01/heat</code>
	<hr><p>by Jairo Andrade — <a style="color:#00ff88" href="https://github.com/Jairocpdev/mission-control-go">GitHub</a> | <a style="color:#00ff88" href="https://jairo-andrade.vercel.app">Portfolio</a></p>
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
	for _, s := range satellites {
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
	json.NewEncoder(w).Encode(Telemetry{
		Satellite: *sat,
		Status:    sat.CheckHealth(),
		Timestamp: time.Now(),
	})
}

func heatHandler(w http.ResponseWriter, r *http.Request) {
	// path: /satellites/{id}/heat
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

func main() {
	// Rotas
	http.HandleFunc("/", swaggerHandler)
	http.HandleFunc("/health", corsMiddleware(healthHandler))
	http.HandleFunc("/satellites", corsMiddleware(satellitesHandler))
	http.HandleFunc("/telemetry/", corsMiddleware(telemetryHandler))
	http.HandleFunc("/satellites/", corsMiddleware(heatHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("📡 Mission Control Go rodando em :%s (Nilópolis, RJ)\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}