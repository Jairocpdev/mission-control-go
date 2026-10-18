package main

import "fmt"

// Satellite é nossa struct - note a letra maiúscula = pública
type Satellite struct {
	ID          string
	Temperature float64
	Battery     float64
}

// CheckHealth usa receiver (s Satellite) - não altera, só lê
func (s Satellite) CheckHealth() string {
	if s.Temperature > 85 || s.Battery < 15 {
		return "CRITICAL"
	}
	if s.Temperature > 75 || s.Battery < 30 {
		return "WARNING"
	}
	return "NOMINAL"
}

// NewSatellite é uma factory - retorna (Satellite, error) - padrão Go
func NewSatellite(id string, temp float64, battery float64) (Satellite, error) {
	if id == "" {
		return Satellite{}, fmt.Errorf("id não pode ser vazio")
	}
	return Satellite{
		ID:          id,
		Temperature: temp,
		Battery:     battery,
	}, nil
}

func main() {
	sat, err := NewSatellite("SAT-01", 72.5, 45.0)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	fmt.Printf("Satélite %s - Status: %s\n", sat.ID, sat.CheckHealth())
}