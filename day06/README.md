# Day 06 - Integração com Portfólio + Simulação Automática

Conexão entre Go e HTML/CSS puro.

## O que foi feito

- **Goroutine de simulação**: 
  ```go
  go func() {
    for {
      time.Sleep(10s)
      sat.Temperature += 1
      if sat.Temperature > 90 { sat.Temperature = 60 }
    }
  }()