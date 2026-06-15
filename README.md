# Go Gin Hexagonal Architecture CRUD API

A small CRUD REST API built with Gin and a Hexagonal Architecture layout.

## Structure

```text
cmd/api                     application entrypoint
internal/core/domain         business entities
internal/core/port           outbound ports
internal/core/service        use cases / inbound port
internal/adapter/http        Gin handlers and routes
internal/adapter/repository  outbound repository adapters
internal/config              runtime configuration
```

## Run

Start PostgreSQL with Docker:

```bash
docker compose up -d
```

PostgreSQL is available on host port `5439`.

```bash
go mod tidy
go run ./cmd/api
```

The API listens on `:8080` by default. Set `PORT` to use another port.

Swagger UI is available at:

```text
http://localhost:8080/swagger/index.html
```

Regenerate Swagger docs after changing API annotations:

```bash
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

Default database settings:

```text
DATABASE_HOST=localhost
DATABASE_PORT=5439
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=hexagonal_architecture
DATABASE_SSLMODE=disable
```

## Endpoints

```text
GET    /health
POST   /api/v1/users
GET    /api/v1/users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
```

## Example

```bash
curl -X POST http://localhost:8080/api/v1/users   -H 'Content-Type: application/json'   -d '{"name":"Jane Doe","email":"jane@example.com"}'
```
