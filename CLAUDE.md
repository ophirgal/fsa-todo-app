# Project Context

This is a full-stack web application built for a coding interview. The goal is to implement the required features as quickly and cleanly as possible.

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
| `backend/internal/handlers/` | One file per resource (e.g., `items.go`) |
| `backend/internal/db/db.go` | `Connect()` and `RunMigrations()` |
| `backend/migrations/` | SQL files run in alphabetical order at startup |
| `frontend/src/api/client.ts` | Pre-wired fetch wrapper — base URL is `/api` |
| `frontend/src/pages/` | Top-level page components |
| `frontend/src/components/` | Reusable UI components |

## Adding a New Resource (End-to-End)

### 1. Migration
Create `backend/migrations/NNN_<name>.sql` (NNN = next number):
```sql
CREATE TABLE IF NOT EXISTS items (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```
Migrations run automatically on next backend start.

### 2. Handler
Create `backend/internal/handlers/items.go`:
```go
package handlers

import (
    "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
)

func ListItems(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        rows, err := db.QueryContext(c.Request.Context(), `SELECT id, name, created_at FROM items ORDER BY id`)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()
        // scan rows...
        c.JSON(http.StatusOK, items)
    }
}
```

### 3. Route
In `backend/internal/router/router.go`, add inside the `api` group:
```go
api.GET("/items", handlers.ListItems(db))
api.POST("/items", handlers.CreateItem(db))
```

Pass `db *sql.DB` to `router.New()` — it already receives it.

### 4. Frontend API call
`frontend/src/api/client.ts` is pre-wired. Use it directly:
```ts
const items = await api.get<Item[]>('/items')
await api.post('/items', { name: 'example' })
```

### 5. Page
Create `frontend/src/pages/ItemsPage.tsx` and add it to the router in `App.tsx`.

## Database Conventions

- All tables use `BIGSERIAL` primary keys named `id`
- Timestamps use `TIMESTAMPTZ` with `DEFAULT NOW()`
- Migration files: `NNN_description.sql` (e.g., `002_add_items.sql`)
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

## Tailwind

Classes are available everywhere in `src/`. No configuration needed — just use utility classes. Common patterns:
- Layout: `flex`, `grid`, `min-h-screen`, `container mx-auto`
- Spacing: `p-4`, `px-6`, `my-8`, `gap-4`
- Typography: `text-xl font-bold text-gray-900`
- Colors: `bg-white`, `text-gray-500`, `border-gray-200`
- Interactive: `hover:bg-blue-600`, `focus:outline-none focus:ring-2`
