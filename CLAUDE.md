# Project Context

This is a simple full-stack TODO application.

## Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.23 + Gin |
| Database | PostgreSQL 16 (via `database/sql` + pgx driver) |
| Frontend | React 18 + TypeScript + TailwindCSS 3 |
| Dev server | Vite 5 (HMR enabled) |
| Reverse proxy | Nginx (single entry point) |
| Orchestration | Docker Compose |

## Running the Project

```bash
cp .env.example .env   # first time only
make up                # build + start all services
make logs              # watch logs
make psql              # database shell
make down              # stop
make clean             # stop + wipe volumes
```

App is at **http://localhost:3000**. API is at **http://localhost:3000/api/**.

## Architecture

```
Browser → nginx:80 (host: 3000)
              ├── /api/*  → backend:8080
              └── /*      → frontend:5173 (dev) / static files (prod)
```

Both backend (Air) and frontend (Vite) hot-reload on file save. No container restarts needed during development.

## Key Files

| File | Purpose |
|---|---|
| `backend/main.go` | Entry point: load config, connect DB, run migrations, start server |
| `backend/internal/router/router.go` | All route registrations go here |
| `backend/internal/handlers/` | One file per resource (e.g., `todos.go`) |
| `backend/internal/db/db.go` | `Connect()` and `RunMigrations()` |
| `backend/migrations/` | SQL files run in alphabetical order at startup |
| `frontend/src/api/client.ts` | Pre-wired fetch wrapper — base URL is `/api` |
| `frontend/src/pages/` | Top-level page components |
| `frontend/src/components/` | Reusable UI components |

## Adding a New Resource (End-to-End)

### 1. Migration
Create `backend/migrations/NNN_<name>.sql` (NNN = next number):
```sql
CREATE TABLE IF NOT EXISTS todos (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT        NOT NULL,
    done       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```
Migrations run automatically on next backend start.

### 2. Handler
Create `backend/internal/handlers/todos.go`:
```go
package handlers

import (
    "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
)

func ListTodos(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        rows, err := db.QueryContext(c.Request.Context(), `SELECT id, title, done, created_at FROM todos ORDER BY id`)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()
        // scan rows...
        c.JSON(http.StatusOK, todos)
    }
}
```

### 3. Route
In `backend/internal/router/router.go`, add inside the `api` group:
```go
api.GET("/todos", handlers.ListTodos(db))
api.POST("/todos", handlers.CreateTodo(db))
```

Pass `db *sql.DB` to `router.New()` — it already receives it.

### 4. Frontend API call
`frontend/src/api/client.ts` is pre-wired. Use it directly:
```ts
const todos = await api.get<Todo[]>('/todos')
await api.post('/todos', { title: 'example' })
```

### 5. Page
Create `frontend/src/pages/TodosPage.tsx` and add it to the router in `App.tsx`.

## Database Conventions

- All tables use `BIGSERIAL` primary keys named `id`
- Timestamps use `TIMESTAMPTZ` with `DEFAULT NOW()`
- Migration files: `NNN_description.sql` (e.g., `002_add_todos.sql`)
- Use `IF NOT EXISTS` guards so migrations are idempotent

## Environment Variables

| Variable | Description |
|---|---|
| `DATABASE_URL` | Full postgres connection string |
| `PORT` | Backend port (default: 8080) |
| `POSTGRES_USER` | DB user (used by postgres container) |
| `POSTGRES_PASSWORD` | DB password |
| `POSTGRES_DB` | DB name |

## Go Module

Module name: `fsa-boilerplate/backend`

Import paths follow `fsa-boilerplate/backend/internal/<package>`.

> Note: The Go module name remains `fsa-boilerplate/backend` for historical reasons — use this exact string in all import paths.

## Tailwind

Classes are available everywhere in `src/`. No configuration needed — just use utility classes. Common patterns:
- Layout: `flex`, `grid`, `min-h-screen`, `container mx-auto`
- Spacing: `p-4`, `px-6`, `my-8`, `gap-4`
- Typography: `text-xl font-bold text-gray-900`
- Colors: `bg-white`, `text-gray-500`, `border-gray-200`
- Interactive: `hover:bg-blue-600`, `focus:outline-none focus:ring-2`

## Autonomy
- When asked to implement a plan, execute it fully end-to-end — including running tests and verification — without asking for permission at each step.

## General Coding Guidelines
- When adding code for a new feature, make minimal changes needed to implement the feature while still sticking to good practices (e.g. do not hardcode things just because it's more minimal). 
- Do not overengineer and add more files than are needed for the feature.
- For backend features, write tests first and use those to validate the feature's implementation (TDD). One happy path test and one error case test is good enough, no need to go crazy with tests.