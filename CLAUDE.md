# search-mm2

Commercial real estate property scraper with web UI.

## Architecture

Monorepo with two services:
- **backend/** — Go API server + scraper (Colly, Chi, pgx, JWT)
- **frontend/** — React + TypeScript + Vite

## Getting Started

```bash
docker compose up              # Start all services (db, backend, frontend)
```

Or run services individually for development:

```bash
docker compose up -d db        # Start PostgreSQL only
cd backend && go run ./cmd/server  # Start API on :8080
cd frontend && npm install && npm run dev  # Start UI on :5173
```

Default credentials: admin / admin

## Key Details

- Frontend proxies `/api` to backend via Vite dev server
- Auth: stateless JWT, admin credentials from env vars
- Scraper: Colly-based, currently supports LoopNet
- DB migrations run automatically on backend startup
- Scraper runs on a configurable interval (default 30m) for active searches

## Testing

```bash
cd backend && go test ./...
cd frontend && npm test
```

## Environment Variables

See `.env.example` for all config options.
