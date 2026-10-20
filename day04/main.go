package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Satellite struct {
	ID          string  `json:"id"`
	Temperature float64 `json:"temperature"`
	Battery     float64 `json:"battery"`
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
	SatelliteID string    `json:"satellite_id"`
	Temp        float64   `json:"temp"`
	Battery     float64   `json:"battery"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

type SatelliteResponse struct {
	ID          string  `json:"id"`
	Temperature float64 `json:"temperature"`
	Battery     float64 `json:"battery"`
	Health      string  `json:"health"`
}

var (
	satellites = map[string]*Satellite{
		"SAT-01": {ID: "SAT-01", Temperature: 70, Battery: 50},
		"SAT-02": {ID: "SAT-02", Temperature: 45, Battery: 10},
		"SAT-03": {ID: "SAT-03", Temperature: 41, Battery: 90},
	}
	mu        sync.Mutex
	startTime = time.Now()
)

func main() {
	fmt.Println("🚀 Mission Control - Day 04: HTTP API")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"uptime": fmt.Sprintf("%.0f segundos", time.Since(startTime).Seconds()),
		})
	})

	http.HandleFunc("/telemetry/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/telemetry/"):]
		mu.Lock()
		sat, ok := satellites[id]
		mu.Unlock()
		if !ok {
			http.Error(w, "satellite not found", http.StatusNotFound)
			return
		}
		tel := Telemetry{
			SatelliteID: sat.ID, Temp: sat.Temperature, Battery: sat.Battery,
			Status: sat.CheckHealth(), Timestamp: time.Now(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tel)
	})

	// Um único handler pra /satellites que resolve lista + heat
	http.HandleFunc("/satellites/", func(w http.ResponseWriter, r *http.Request) {
		// GET /satellites -> lista
		if r.URL.Path == "/satellites" || r.URL.Path == "/satellites/" {
			if r.Method == http.MethodGet {
				mu.Lock()
				var list []SatelliteResponse
				for _, sat := range satellites {
					list = append(list, SatelliteResponse{ID: sat.ID, Temperature: sat.Temperature, Battery: sat.Battery, Health: sat.CheckHealth()})
				}
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(list)
				return
			}
		}

		// POST /satellites/SAT-01/heat
		if r.Method == http.MethodPost {
			path := r.URL.Path
			if len(path) < len("/satellites//heat") {
				http.Error(w, "invalid path", http.StatusBadRequest)
				return
			}
			id := path[len("/satellites/") : len(path)-len("/heat")]

			mu.Lock()
			sat, ok := satellites[id]
			if ok {
				sat.Temperature += 5
			}
			mu.Unlock()

			if !ok {
				http.Error(w, "satellite not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message":  fmt.Sprintf("%s heated by 5C", id),
				"new_temp": sat.Temperature,
				"status":   sat.CheckHealth(),
			})
			return
		}

		http.Error(w, "not found", http.StatusNotFound)
	})

	fmt.Println("📡 Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}