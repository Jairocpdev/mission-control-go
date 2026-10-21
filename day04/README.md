### 4. Cria `mission-control-go/day04/README.md`
```md
# Day 04 - HTTP API REST

Primeira API em Go puro.

Endpoints:
- `GET /health` → `{"status":"ok","uptime":"89 segundos"}`
- `GET /satellites` → lista com health
- `GET /telemetry/{id}` → dado de um satélite
- `POST /satellites/{id}/heat` → aquece e muda status

Destaque:
- `sync.Mutex` protegendo mapa em handlers concorrentes
- JSON com `encoding/json`
- Status que muda ao vivo: 70 NOMINAL → 80 WARNING

Rodar:
```powershell
go run ./day04

# Terminal 2
Invoke-RestMethod -Uri http://localhost:8080/health
Invoke-RestMethod -Method POST -Uri http://localhost:8080/satellites/SAT-01/heat