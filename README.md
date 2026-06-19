# Go Gin Hexagonal Architecture CRUD API

A small CRUD REST API built with Gin and a Hexagonal Architecture layout.

## Structure

```text
.
├── cmd/
│   └── api/                            # application entrypoint
├── db/                                 # database init scripts (mysql, postgres)
├── internal/
│   ├── adapter/
│   │   ├── inbound/
│   │   │   ├── grpc/                   # inbound gRPC adapter
│   │   │   └── http/                   # inbound HTTP adapter (Handlers, Router)
│   │   └── outbound/
│   │       ├── clock/                  # system clock adapter
│   │       ├── event/                  # outbound event/API adapter with circuit breaker
│   │       ├── id/                     # ID generator adapter
│   │       ├── repository/             # outbound database adapters (mysql, postgres)
│   │       └── storage/                # outbound storage adapter (minio)
│   ├── core/
│   │   ├── domain/                     # business entities
│   │   ├── port/                       # ports / interfaces
│   │   └── service/                    # use cases / business logic
│   └── infrastructure/
│       ├── config/                     # runtime configuration
│       ├── database/                   # database pool connections
│       └── observability/              # logger and tracer setup
├── pkg/                                # reusable public packages, when needed
├── docs/                               # generated Swagger docs
└── docker-compose.yml                  # local infrastructure orchestration
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
OUTBOUND_API_BASE_URL=
```

## Endpoints

```text
GET    /health
GET    /metrics
POST   /api/v1/users
GET    /api/v1/users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
```

## Observability

The HTTP router includes the Phase 1 observability pieces:

```text
Recovery middleware
Request ID middleware
Context-aware zap logger
OpenTelemetry request tracing
Prometheus metrics middleware
/metrics endpoint
```

## Wiring

Application dependencies are wired in [cmd/api/main.go](cmd/api/main.go) in this order:

```text
load config
init logger and tracer
init metrics registry
init database pool (MySQL/PostgreSQL)
init ID generator and clock
inject pool into repository adapter
inject repository and outbound adapter into usecase
inject usecase into HTTP handler/router
start HTTP server with graceful shutdown
```

## Outbound API Adapter

The project includes an outbound HTTP API adapter with an embedded circuit breaker:

```text
internal/core/port/user_event.go
internal/adapter/outbound/event/httpclient
```

The core depends on the business port `UserEventPublisher`; HTTP request details stay inside the outbound adapter.
The circuit opens after repeated failed external calls, blocks requests while open, then allows a half-open probe after the timeout.
When `OUTBOUND_API_BASE_URL` is empty, the app injects a no-op outbound adapter so local CRUD remains self-contained.

## Example

```bash
curl -X POST http://localhost:8080/api/v1/users   -H 'Content-Type: application/json'   -d '{"name":"Jane Doe","email":"jane@example.com"}'
```
