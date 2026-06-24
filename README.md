# Book Search

A full-stack SPA that searches books by author using the [Open Library](https://openlibrary.org) API.

- **Backend**: Go (standard library only) — proxies Open Library and exposes a clean REST API
- **Frontend**: Angular 19 with standalone components and signals

## Project structure

```
book-search/
├── backend/          Go HTTP server
│   ├── main.go
│   └── go.mod
└── frontend/         Angular SPA
    └── src/app/
        ├── app.ts                           Root component (search + state)
        ├── book.service.ts                  HTTP service
        ├── models/book.model.ts             Interfaces
        └── components/
            ├── book-card/                   Individual book card
            └── book-list/                   Responsive card grid
```

## Running locally

### 1. Start the Go backend

```bash
cd backend
go run main.go
# Listening on :8080
```

### 2. Start the Angular dev server

```bash
cd frontend
npm start
# App available at http://localhost:4200
```

### API

| Method | Path | Query params | Description |
|--------|------|-------------|-------------|
| GET | `/api/books` | `author` (required), `limit` (default 20, max 100) | Search books by author |
| GET | `/health` | — | Health check |

Example:

```
GET http://localhost:8080/api/books?author=tolkien&limit=10
```

```json
{
  "books": [
    {
      "key": "/works/OL45804W",
      "title": "The Fellowship of the Ring",
      "author_names": ["J. R. R. Tolkien"],
      "first_publish_year": 1954,
      "cover_url": "https://covers.openlibrary.org/b/id/8473352-M.jpg"
    }
  ],
  "total": 142
}
```
