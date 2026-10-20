# Day 01 - Structs & Health Check

Primeiro contato com Go.

- Struct `Satellite` com `ID`, `Temperature`, `Battery`
- Method `CheckHealth() string`
    - CRITICAL: temp > 85 ou battery < 15
    - WARNING: temp > 75 ou battery < 30
    - NOMINAL: resto
- Uso de `fmt.Printf`

Rodar:
```bash
go run ./day01