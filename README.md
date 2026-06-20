# TODO App

A simple full-stack TODO application.

**Stack:** Go + Gin · React + TypeScript + TailwindCSS · PostgreSQL · Nginx · Docker Compose

---

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) + Docker Compose v2
- `make`

---

## Quick Start

```bash
cp .env.example .env
make up
```

- App: http://localhost:3000
- API health: http://localhost:3000/api/health

---

## Architecture

```
Browser → nginx:3000
              ├── /api/*  → backend:8080  (Go + Gin, hot reload via Air)
              └── /*      → frontend:5173 (React + Vite, HMR)
                                │
                           postgres:5432
```

All traffic flows through a single Nginx entry point — no CORS configuration needed.

---

## Development

| Command | Description |
|---|---|
| `make up` | Build and start all services |
| `make down` | Stop all services |
| `make logs` | Tail logs from all services |
| `make psql` | Open a psql shell in the database |
| `make clean` | Stop services and wipe volumes |

Both the backend (Air) and frontend (Vite) support hot reload — file saves are reflected immediately without restarting containers.

---

## Project Structure

```
fsa-todo-app/
├── backend/
│   ├── internal/
│   │   ├── config/      # env var loading
│   │   ├── db/          # connection pool + migration runner
│   │   ├── handlers/    # Gin handler functions
│   │   ├── middleware/  # CORS, auth, etc.
│   │   └── router/      # route registration
│   └── migrations/      # *.sql files run in order at startup
├── frontend/
│   └── src/
│       ├── api/         # fetch wrapper (base URL: /api)
│       ├── components/  # shared UI components
│       └── pages/       # route-level page components
└── nginx/
    ├── nginx.conf       # dev proxy config
    └── nginx.prod.conf  # prod static file server + proxy
```

---

## Adding a New Resource

**1. Migration** — `backend/migrations/002_<name>.sql`
```sql
CREATE TABLE IF NOT EXISTS todos (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT        NOT NULL,
    done       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**2. Handler** — `backend/internal/handlers/todos.go`
```go
func ListTodos(c *gin.Context) { ... }
func CreateTodo(c *gin.Context) { ... }
```

**3. Route** — add to `backend/internal/router/router.go`
```go
api.GET("/todos", handlers.ListTodos)
api.POST("/todos", handlers.CreateTodo)
```

**4. API call** — `frontend/src/api/client.ts` is pre-wired
```ts
const todos = await api.get<Todo[]>('/todos')
```

**5. Page** — `frontend/src/pages/TodosPage.tsx`

---

## Production

```bash
make build-prod   # builds static frontend, runs backend binary
```

In production, Nginx serves the compiled React assets from `dist/` and proxies `/api/*` to the Go backend. No Vite or Node process runs.
