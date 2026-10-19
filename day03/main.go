package main

import (
	"fmt"
	"sync"
	"time"
)

type Satellite struct {
	ID          string
	Temperature float64
	Battery     float64
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
	SatelliteID string
	Temp        float64
	Battery     float64
	Status      string
	Timestamp   time.Time
}

func (s *Satellite) StreamTelemetry(ch chan Telemetry, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)

	for i := 0; i < 5; i++ {
		s.Temperature += 2.5

		// TODO 1 RESOLVIDO: consumo de bateria
		s.Battery -= 1.5
		if s.Battery < 0 {
			s.Battery = 0
		}

		t := Telemetry{
			SatelliteID: s.ID,
			Temp:        s.Temperature,
			Battery:     s.Battery,
			Status:      s.CheckHealth(),
			Timestamp:   time.Now(),
		}
		ch <- t
		time.Sleep(400 * time.Millisecond)
	}
}

func main() {
	fmt.Println("🚀 Mission Control - Day 03 TURBO: Structs + Select + Metrics")

	sat1 := &Satellite{ID: "SAT-01", Temperature: 70, Battery: 50}
	sat2 := &Satellite{ID: "SAT-02", Temperature: 45, Battery: 10}
	sat3 := &Satellite{ID: "SAT-03", Temperature: 41, Battery: 90}

	ch1 := make(chan Telemetry, 5)
	ch2 := make(chan Telemetry, 5)
	ch3 := make(chan Telemetry, 5)

	var wg sync.WaitGroup
	wg.Add(3)

	go sat1.StreamTelemetry(ch1, &wg)
	go sat2.StreamTelemetry(ch2, &wg)
	go sat3.StreamTelemetry(ch3, &wg)

	go func() {
		wg.Wait()
	}()

	// TODO 2 RESOLVIDO: contadores
	criticalCount := 0
	warningCount := 0

	ch1Closed, ch2Closed, ch3Closed := false, false, false
	for !ch1Closed || !ch2Closed || !ch3Closed {
		select {
		case t, ok := <-ch1:
			if !ok {
				ch1Closed = true
				continue
			}
			if t.Status == "CRITICAL" {
				criticalCount++
			}
			if t.Status == "WARNING" {
				warningCount++
			}
			fmt.Printf("[%s] %.1f°C | %s | Bat: %.1f%% | %s\n", t.SatelliteID, t.Temp, t.Status, t.Battery, t.Timestamp.Format("15:04:05"))

		case t, ok := <-ch2:
			if !ok {
				ch2Closed = true
				continue
			}
			if t.Status == "CRITICAL" {
				criticalCount++
			}
			if t.Status == "WARNING" {
				warningCount++
			}
			fmt.Printf("[%s] %.1f°C | %s | Bat: %.1f%% | %s\n", t.SatelliteID, t.Temp, t.Status, t.Battery, t.Timestamp.Format("15:04:05"))

		case t, ok := <-ch3:
			if !ok {
				ch3Closed = true
				continue
			}
			if t.Status == "CRITICAL" {
				criticalCount++
			}
			if t.Status == "WARNING" {
				warningCount++
			}
			fmt.Printf("[%s] %.1f°C | %s | Bat: %.1f%% | %s\n", t.SatelliteID, t.Temp, t.Status, t.Battery, t.Timestamp.Format("15:04:05"))

		// TODO 3 RESOLVIDO: timeout
		case <-time.After(1 * time.Second):
			fmt.Println("⚠️ Timeout - perda de sinal")
		}
	}

	// TODO 4 RESOLVIDO: resumo
	fmt.Printf("\n📊 Resumo: %d CRITICAL, %d WARNING\n", criticalCount, warningCount)
	if criticalCount > 0 {
		fmt.Println("🚨 AÇÃO NECESSÁRIA: Verificar satélites críticos")
	}

	fmt.Println("\n✅ Mission Control encerrando")
}