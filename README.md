# Prosel

Prosel is a small, personal blog system with a Go API backend and a Next.js frontend.

## Services

- Backend API: http://localhost:8080
- Frontend: http://localhost:3000
- PostgreSQL: localhost:5432
- Redis: localhost:6379

## Local development

1. Copy environment defaults if you want to run services directly:

   ```bash
   cp .env.example .env
   ```

2. Start all services with Docker Compose:

   ```bash
   docker compose up --build
   ```

3. Check backend health:

   ```bash
   curl http://localhost:8080/api/v1/health
   ```

4. Check public settings:

   ```bash
   curl http://localhost:8080/api/v1/settings/public
   ```

## Backend commands

```bash
cd backend
go test ./...
go build ./cmd/api
```

The backend runs SQL migrations from `backend/migrations` on startup.

## Frontend commands

```bash
cd frontend
npm install
npm run typecheck
npm run build
```

## Environment variables

See `.env.example` for all supported variables. The Docker Compose defaults use PostgreSQL user/database/password `prosel` and Redis without a password.
