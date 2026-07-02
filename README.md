# Book Search

Book Search is a full-stack app that lets you search books by author using the Open Library API.

- Backend: Go HTTP API
- Frontend: Angular app served by Nginx
- Orchestration: Docker Compose

## Prerequisites

- Docker Desktop (or Docker Engine + Docker Compose plugin)

Verify installation:

```bash
docker --version
docker compose version
```

## Start frontend and backend with Docker Compose

From the repository root:

```bash
docker compose up --build
```

This command will:

1. Build the backend image from `backend/Dockerfile`
2. Build the frontend image from `frontend/Dockerfile`
3. Start both containers

## URLs

- Frontend UI: http://localhost:3000
- Backend API (direct): http://localhost:8080
- Backend health: http://localhost:8080/health

The frontend calls `/api/*`, and Nginx proxies those requests to the backend container (`backend:8080`).

## Stop the stack

```bash
docker compose down
```

To stop and also remove images built by compose:

```bash
docker compose down --rmi local
```

## Useful commands

Show running services:

```bash
docker compose ps
```

View logs:

```bash
docker compose logs -f
```

Rebuild after code changes:

```bash
docker compose up --build
```

## API quick reference

| Method | Path | Query params | Description |
|--------|------|-------------|-------------|
| GET | `/api/books` | `author` (required), `limit` (default 20, max 100) | Search books by author |
| GET | `/health` | none | Health check |

Example:

```text
GET http://localhost:8080/api/books?author=tolkien&limit=10
```

## Troubleshooting

- Port already in use:
  - If `3000` or `8080` is occupied, stop the conflicting process/container and retry.
- Frontend loads but API fails:
  - Confirm backend is healthy at `http://localhost:8080/health`.
  - Check logs with `docker compose logs -f backend frontend`.
- Build issues:
  - Retry with a clean build: `docker compose build --no-cache` then `docker compose up`.
