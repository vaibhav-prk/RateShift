# WardGate

**WardGate** is a lightweight, multi-tenant API Gateway and Rate Limiting service built with Go, Redis, and Chi. It acts as a reverse proxy that authenticates incoming client requests, evaluates rate limits on a per-tenant basis, and forwards permitted traffic to backend services.

---

## Architecture Overview

```text
                    +----------------------+
                    |     Client / User    |
                    +----------+-----------+
                               |
                 X-API-Key Header / HTTP Request
                               |
                               v
                    +----------------------+
                    |    WardGate Gateway  |
                    |      (Port 8080)     |
                    +----------+-----------+
                               |
          +--------------------+--------------------+
          |                                         |
          v                                         v
+---------------------+                   +---------------------+
|   Auth Middleware   |                   | RateLimit Middleware|
| (Validates API Key) |                   | (Checks Redis limit)|
+----------+----------+                   +----------+----------+
          |                                         |
          +--------------------+--------------------+
                               |
                      Allowed Traffic Only
                               |
                               v
                    +----------------------+
                    |     Mock Backend     |
                    |      (Port 8081)     |
                    +----------------------+
```

---

## Key Features

* **Multi-Tenant Architecture**: Configure rate limits, burst sizes, and algorithms per tenant independently in standard YAML.
* **API Key Authentication**: Header-based `X-API-Key` authentication matching keys directly to tenant profiles.
* **Extensible Rate Limiter Interface**: Core design supports multiple rate-limiting algorithms including `token_bucket`, `leaky_bucket`, and `sliding_window`.
* **Dynamic Headers**: Standardized response headers including `X-RateLimit-Remaining` and `Retry-After` on limit breaches (`HTTP 429`).
* **Reverse Proxy**: Native HTTP reverse proxy forwarding authorized requests to designated upstream backend services.
* **Graceful Shutdown**: Handles `SIGINT` and `SIGTERM` signals for clean server termination without dropping active connections.
* **Container-Ready**: Built-in multi-stage Docker builds and Docker Compose orchestration with Redis.

---

## Project Structure

```text
.
├── cmd/
│   └── gateway/
│       └── main.go                 # Gateway server entry point & routing
├── deploy/
│   ├── mock-backend/
│   │   ├── Dockerfile              # Dockerfile for mock backend service
│   │   └── main.go                 # Lightweight upstream mock service
│   ├── Dockerfile                  # Gateway multi-stage Docker build
│   └── docker-compose.yml          # Full stack Docker Compose deployment
├── internal/
│   ├── auth/
│   │   └── auth.go                 # Tenant API key authentication engine
│   ├── config/
│   │   └── config.go               # YAML configuration loader
│   ├── limiter/
│   │   └── interface.go            # Rate limiter interfaces & types
│   └── middleware/
│       ├── auth.go                 # HTTP authentication middleware
│       └── ratelimit.go            # HTTP rate-limiting middleware
├── config.yaml                     # Tenant profiles & limit rules
├── go.mod                          # Go module definition
└── go.sum                          # Go module dependencies checksums
```

---

## Configuration

Tenants and their rate limits are defined inside `config.yaml`:

```yaml
tenants:
  tenant-a:
    tenant_id: tenant-a
    api_key: tenant-a-key
    rate_limit: 10
    burst_size: 20
    algorithm: token_bucket

  tenant-b:
    tenant_id: tenant-b
    api_key: tenant-b-key
    rate_limit: 5
    burst_size: 10
    algorithm: token_bucket
```

### Supported Limiter Algorithms

* `token_bucket`
* `leaky_bucket`
* `sliding_window`

---

## Getting Started

### Prerequisites

* [Go](https://golang.org/) 1.26+ installed locally
* [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
* [Redis](https://redis.io/) running locally or via Docker

---

## Running with Docker Compose

Docker Compose is the recommended way to run the complete stack.

### 1. Clone the Repository

```bash
git clone https://github.com/vaibhav-prk/WardGate.git
cd WardGate
```

### 2. Start All Services

```bash
docker-compose -f deploy/docker-compose.yml up --build
```

The services will start on:

| Service      | Address                 |
| ------------ | ----------------------- |
| Gateway      | `http://localhost:8080` |
| Mock Backend | `http://localhost:8081` |
| Redis        | `localhost:6379`        |

---

## Running Locally

### 1. Start Redis

```bash
docker run -p 6379:6379 -d redis:alpine
```

### 2. Run the Mock Backend

```bash
go run deploy/mock-backend/main.go
```

### 3. Run the Gateway

```bash
REDIS_URL="localhost:6379" BACKEND_URL="http://localhost:8081" go run cmd/gateway/main.go
```

---

# API Endpoints & Testing

## Health Check

Public endpoint — authentication is not required.

```bash
curl -i http://localhost:8080/health
```

### Response

```json
{
  "status": "healthy"
}
```

---

## Proxied API Request

Pass a valid `X-API-Key` configured in `config.yaml`:

```bash
curl -i -H "X-API-Key: tenant-a-key" http://localhost:8080/api/users
```

### Successful Response

```http
HTTP/1.1 200 OK
Content-Type: application/json
X-RateLimit-Remaining: 19
```

```json
{
  "headers": {},
  "method": "GET",
  "path": "/api/users",
  "status": "ok",
  "time": "2026-10-05T18:46:27Z"
}
```

---

## Rate Limit Exceeded Response

When a tenant exceeds their configured rate limit:

```bash
curl -i -H "X-API-Key: tenant-b-key" http://localhost:8080/api/test
```

### Response

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
X-RateLimit-Remaining: 0
Retry-After: 5
```

```json
{
  "error": "rate limit exceeded"
}
```

---

## Error Responses

| Scenario               | HTTP Code               | Response Body                      |
| ---------------------- | ----------------------- | ---------------------------------- |
| Missing API Key Header | `401 Unauthorized`      | `{"error": "missing api key"}`     |
| Invalid API Key        | `401 Unauthorized`      | `{"error": "invalid api key"}`     |
| Rate Limit Breached    | `429 Too Many Requests` | `{"error": "rate limit exceeded"}` |

---

## Environment Variables

| Variable      | Description                   | Default Value           |
| ------------- | ----------------------------- | ----------------------- |
| `REDIS_URL`   | Redis host and port           | `localhost:6379`        |
| `BACKEND_URL` | Upstream backend proxy target | `http://localhost:8081` |

---

## Dependencies

* https://github.com/go-chi/chi — Lightweight HTTP router
* https://github.com/redis/go-redis — Redis Go client
* https://gopkg.in/yaml.v3 — YAML configuration parser

---

## License

This project is available for educational and development purposes.
