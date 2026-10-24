# 🚀 Mission Control Go

API de telemetria de satélites em **Go puro**, construída do zero em 7 dias. Simula uma frota em tempo real com aquecimento automático, histórico e dashboard ao vivo no portfólio.

**Live:** 
- API: https://mission-control-go.onrender.com
- Portfólio: https://jairo-andrade.vercel.app
- Código: https://github.com/Jairocpdev/mission-control-go

## 📡 O que faz

Frota de 3 satélites (SAT-01, SAT-02, SAT-03) com:

- `CheckHealth()` → NOMINAL / WARNING (>75°C ou <30% bat) / CRITICAL (>85°C ou <15% bat)
- `sync.Mutex` para concorrência segura
- Goroutine auto-heating: +1°C a cada 10s, reset em 90°C → 60°C
- Histórico de 20 pontos por satélite em **horário BRT (Nilópolis/RJ)**
- Gráfico Chart.js ao vivo no portfólio

## 🔌 Endpoints

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/` | Swagger HTML |
| GET | `/health` | Uptime |
| GET | `/satellites` | Lista ordenada SAT-01,02,03 com status |
| GET | `/telemetry/{id}` | Telemetria + timestamp |
| POST | `/satellites/{id}/heat` | Aquece +5°C |
| GET | `/history/{id}` | Últimos 20 pontos `{time: "12:00:52", temp: 81}` BRT |

curl https://mission-control-go.onrender.com/satellites
curl -X POST https://mission-control-go.onrender.com/satellites/SAT-01/heat
curl https://mission-control-go.onrender.com/history/SAT-01
