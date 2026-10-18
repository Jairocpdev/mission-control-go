package main

import (
	"fmt"
	"time"
)

// Reaproveita o que você já fez - mas agora com método que simula telemetria
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

// NOVO DIA 2: Método que vai rodar em uma goroutine
// Ele recebe um channel (canal) pra enviar dados de volta
func (s *Satellite) StreamTelemetry(ch chan string) {
	for i := 0; i < 5; i++ {
		s.Temperature += 2.5
		msg := fmt.Sprintf("[%s] Temp: %.1f°C | Status: %s", s.ID, s.Temperature, s.CheckHealth())
		ch <- msg
		time.Sleep(500 * time.Millisecond)
	}
	close(ch)
}

func main() {
	fmt.Println("🚀 Mission Control - Day 02: Concorrência com Goroutines")

	// TODO 1: Crie 3 satélites
	
	sat1 := Satellite{ID: "SAT-01", Temperature: 70, Battery: 50}
	sat2 := Satellite{ID: "SAT-02", Temperature: 45, Battery: 10}
	sat3 := Satellite{ID: "SAT-03", Temperature: 41, Battery: 90}

	// TODO 2: Crie 3 channels
	ch1 := make (chan string)
	ch2 := make (chan string)
	ch3 := make (chan string)

	// TODO 3: Dispare 3 goroutines (é só colocar 'go' na frente)
	go sat1.StreamTelemetry(ch1)
	go sat2.StreamTelemetry(ch2)
	go sat3.StreamTelemetry(ch3)
	
	// TODO 4: Leia os channels em paralelo (dica: for range)
	
	for msg := range ch1 {
	fmt.Println(msg)
	 }

	for msg := range ch2 {
	fmt.Println(msg)
	 }

	for msg := range ch3 {
	fmt.Println(msg)
	 }
}