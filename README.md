# Getting Fast at Go

A step-by-step guide to building a minimal Go API with SQLite, Docker, and tests. This project is a Go implementation of the FastAPI-101 project, demonstrating core Go web development concepts: routing, middleware, database integration, and API design.

---

## What's Included

1. **A minimal Go application** (`main.go`) with root, health, and items API
2. **Go module dependencies** (`go.mod`, `go.sum`) for package management
3. **A Docker image** (`Dockerfile`) that runs the app in a container
4. **Docker Compose** (`docker-compose.yml`) for one-command run
5. **A `.dockerignore`** so unnecessary files stay out of the image
6. **A persistent database** (SQLite + GORM) for items
7. **Data models** (`models.go`) with request/response schemas
8. **API key authentication** (middleware) for protecting endpoints
9. **Service layer** in `main.go` separating business logic from routes

---

## Quick Start

```bash
# Copy environment variables template (optional)
cp .env.example .env

# Start the app
docker compose up --build

# Or run locally (requires Go 1.22+)
go run . 
```

Then open:
- **http://localhost:8000** – API root
- **http://localhost:8000/health** – Health check
- **http://localhost:8000/items** – List items

---

## Project Structure

```
go-101/
├── main.go              # Go application with routes and handlers
├── models.go            # Data models (Item, ItemCreate, ItemUpdate)
├── auth.go              # API key authentication middleware
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── Dockerfile           # How to build the container image
├── docker-compose.yml   # How to run the container
├── .dockerignore        # Files to exclude from Docker build
├── .env.example         # Environment variables template
├── .gitignore           # Git ignore rules
└── README.md            # This file
```

---

## Dependencies

**What we put in `go.mod`:**

| Package | Purpose |
|---------|----------|
| `gin-gonic/gin` | Web framework: routing, middleware, validation |
| `gorm.io/gorm` | ORM: database abstraction, migrations |
| `gorm.io/driver/sqlite` | SQLite driver for GORM |
| `mattn/go-sqlite3` | SQLite C bindings for Go |

---

## The Go App: `main.go`

**What it is:** The main application entry point. Registers routes and handlers.

**Key concepts:**

- **`gin.Default()`** – Creates a router with default middleware
- **`router.GET()`, `router.POST()`, etc.** – Register handlers for HTTP methods
- **`gin.H`** – Map type for JSON responses

**Endpoints:**

| Path | Method | Purpose |
|------|--------|----------|
| `/` | GET | Simple hello message |
| `/health` | GET | Health check |
| `/items` | GET | List items with pagination |
| `/items/{item_id}` | GET | Get a single item |
| `/items` | POST | Create a new item |
| `/items/{item_id}` | PATCH | Update an item (partial) |
| `/items/{item_id}` | DELETE | Delete an item (requires API key) |
| `/items/stats/summary` | GET | Get statistics about items |

---

## Models: `models.go`

Defines the data structures:

- **`Item`** – Database model with GORM tags
- **`ItemCreate`** – Request body schema for creating items
- **`ItemUpdate`** – Request body schema for partial updates

---

## Authentication: `auth.go`

Provides API key authentication via middleware:

- Reads `API_KEY` environment variable (default: `dev-key-123`)
- Middleware checks `X-API-Key` header
- Returns 401 if missing or invalid

---

## Database

Uses **SQLite** with **GORM**:

- **DATABASE_URL** env var (default: `app.db`)
- Auto-migration on startup
- Supports typical SQL operations (SELECT, INSERT, UPDATE, DELETE)

---

## How to Run

### Using Docker Compose (recommended)

```bash
docker compose up --build
```

Stop with `Ctrl+C`.

### Using Docker only (no Compose)

```bash
docker build -t go-101 .
docker run -p 8000:8000 -e DATABASE_URL=app.db -e API_KEY=dev-key-123 go-101
```

### Without Docker (local Go)

```bash
# Requires Go 1.22 or later
go install github.com/cosmtrek/air@latest  # for hot reload (optional)
go run .

# Or with hot reload:
air
```

---

## Environment Variables

Create a `.env` file (or set these directly):

```env
DATABASE_URL=app.db          # SQLite database file path
PORT=8000                     # Server port
API_KEY=dev-key-123          # API key for authentication
```

---

## API Examples

### Create an item

```bash
curl -X POST http://localhost:8000/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Widget",
    "description": "A nice widget",
    "price": 9.99,
    "category": "gadgets"
  }'
```

### List items

```bash
curl http://localhost:8000/items

# With pagination
curl "http://localhost:8000/items?skip=0&limit=5"
```

### Get a single item

```bash
curl http://localhost:8000/items/1
```

### Update an item

```bash
curl -X PATCH http://localhost:8000/items/1 \
  -H "Content-Type: application/json" \
  -d '{"price": 14.99}'
```

### Delete an item (requires API key)

```bash
curl -X DELETE http://localhost:8000/items/1 \
  -H "X-API-Key: dev-key-123"
```

### Get statistics

```bash
curl http://localhost:8000/items/stats/summary
```

---

## Next Steps

1. **Add tests** – Use `testing` package or a test framework like `testify`
2. **Add logging** – Use `slog` or a library like `logrus` or `zap`
3. **Add validation** – Use `playground/validator` for advanced validation
4. **Add migrations** – Use `golang-migrate` for database schema management
5. **Add CI/CD** – Set up GitHub Actions for testing and linting

---

## Comparison: FastAPI vs Go

| Feature | FastAPI | Go |
|---------|---------|----|
| **Language** | Python | Go |
| **Web Framework** | FastAPI | Gin (or net/http) |
| **ORM** | SQLAlchemy | GORM |
| **Database** | SQLite | SQLite |
| **Performance** | Fast (async) | Very fast (concurrent) |
| **Startup Time** | Moderate | Very fast |
| **Binary Size** | N/A (interpreted) | Small (compiled) |
| **Concurrency** | asyncio/await | goroutines/channels |

---

## Resources

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [Go Official Documentation](https://golang.org/doc/)
- [SQLite](https://www.sqlite.org/)
