# aegis-demo-go

Gin + cookie session. Target for Aegis-generated scanner pipelines.

## Run

```bash
go mod tidy
go run .
# open http://localhost:8080/login
# demo creds: demo@example.com / demo1234
```

## Endpoints

- `GET /` — public landing
- `GET /login` / `POST /login` — form auth
- `GET /dashboard` — requires session
- `GET /api/items` — requires session
- `POST /logout`
- `GET /healthz`
