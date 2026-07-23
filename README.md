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
10. **Tests** with `testify` mirroring the FastAPI-101 test suite
11. **Structured logging** with `slog` and request middleware
12. **Request validation** with `playground/validator`
13. **SQL migrations** with `golang-migrate`
14. **CI/CD** with GitHub Actions for testing and linting

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
├── main.go              # Application entry point
├── internal/app/        # Application code (handlers, models, auth, etc.)
├── migrations/          # SQL migration files (golang-migrate)
├── tests/
│   ├── testcase/        # Laravel-style TestCase base with HTTP helpers
│   ├── feature/         # HTTP / integration tests
│   └── unit/            # Unit tests (auth, validation)
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── Makefile             # Development commands (test, lint, build)
├── .github/workflows/   # GitHub Actions CI
├── Dockerfile           # How to build the container image
├── Dockerfile.dev       # Dev image for tests and linting
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
| `gorm.io/gorm` | ORM: database abstraction |
| `gorm.io/driver/sqlite` | SQLite driver for GORM |
| `mattn/go-sqlite3` | SQLite C bindings for Go |
| `golang-migrate/migrate` | Database schema migrations |
| `go-playground/validator` | Advanced request validation |
| `stretchr/testify` | Test assertions |

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

Uses **SQLite** with **GORM** and **golang-migrate**:

- **DATABASE_URL** env var (default: `app.db`)
- SQL migrations in `migrations/` applied on startup
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

## Testing

Tests follow a Laravel-style layout:

| Directory | Purpose |
|-----------|---------|
| `tests/feature/` | HTTP integration tests (like Laravel Feature tests) |
| `tests/unit/` | Isolated unit tests (like Laravel Unit tests) |
| `tests/testcase/` | Shared `TestCase` base with `Get`, `Post`, `ResetDatabase`, etc. |

Run the test suite (mirrors the FastAPI-101 test coverage):

```bash
make test

# With local Go installed:
make test-local
```

## Development

Development commands run inside Docker by default (no local Go install required):

```bash
make lint          # Run golangci-lint
make test-coverage # Run tests with coverage report
make build         # Build binary
make docker-shell  # Interactive dev container shell
```

With local Go and golangci-lint installed, use `make test-local`, `make lint-local`, or `make install-tools`.

CI runs on every push and pull request to `main` via GitHub Actions (`.github/workflows/ci.yml`).

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

## *-101 Family

### API backends

| Repo | Port | Type | Stack |
|------|------|------|-------|
| [fastAPI-101](https://github.com/iammikek/fastAPI-101) | 8000 | API-only | FastAPI, SQLAlchemy |
| [django-101](https://github.com/iammikek/django-101) | 8001 | Monolith | Django + DRF + shop |
| [symfony-101](https://github.com/iammikek/symfony-101) | 8002 | Monolith | Symfony + shop |
| [laravel-101](https://github.com/iammikek/laravel-101) | 8003 | Monolith | Laravel + shop |
| [framework-x-101](https://github.com/iammikek/framework-x-101) | 8004 | Monolith | Framework X + shop |
| [orchestr-101](https://github.com/iammikek/orchestr-101) | 8005 | Monolith | Orchestr + shop |
| [nest-101](https://github.com/iammikek/nest-101) | 8006 | API-only | NestJS, TypeScript |
| [express-101](https://github.com/iammikek/express-101) | 8007 | API-only | Express, Vitest |
| [**go-101**](https://github.com/iammikek/go-101) | 8000* | API-only | Gin, GORM |
| [fortran-101](https://github.com/iammikek/fortran-101) | 8008 | API-only | Fortran, fpm |
| [java-101](https://github.com/iammikek/java-101) | 8009 | API-only | Spring Boot, JPA, Flyway |
| [dotNet-101](https://github.com/iammikek/dotNet-101) | 8010 | API-only | ASP.NET Core, xUnit |
| [flask-101](https://github.com/iammikek/flask-101) | 8011 | API-only | Flask, pytest |
| [rails-101](https://github.com/iammikek/rails-101) | 8012 | Monolith | Rails + shop |
| [geblang-101](https://github.com/iammikek/geblang-101) | 8013 | API-only | Geblang, SQLite |
| [gebweb-101](https://github.com/iammikek/gebweb-101) | 8014 | API-only | Geblang + Gebweb |
\* go-101 also uses port 8000 — run one backend at a time, or change port in config.

### Other clients

| Repo | Platform | Stack |
|------|----------|-------|
| [flutter-101](https://github.com/iammikek/flutter-101) | Mobile / desktop | Flutter (iOS, macOS, Android) |
| [react-101](https://github.com/iammikek/react-101) | Web browser | React 19, Vite, Vitest |
| [vue-101](https://github.com/iammikek/vue-101) | Web browser | Vue 3, Vite, Pinia |
| [alpine-101](https://github.com/iammikek/alpine-101) | Web browser | Alpine.js, Vite, Vitest |

### Suggested pairing

- **Compare compiled backends:** go-101 vs [fastAPI-101](https://github.com/iammikek/fastAPI-101), [fortran-101](https://github.com/iammikek/fortran-101) (8008), or [java-101](https://github.com/iammikek/java-101) (8009)
- **Pair with a client:** [react-101](https://github.com/iammikek/react-101), [vue-101](https://github.com/iammikek/vue-101), [alpine-101](https://github.com/iammikek/alpine-101), or [flutter-101](https://github.com/iammikek/flutter-101)

Catalogue: [automica.io/learning-101](https://automica.io/learning-101.html)

---

## Resources

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [Go Official Documentation](https://golang.org/doc/)
- [SQLite](https://www.sqlite.org/)
