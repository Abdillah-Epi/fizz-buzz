# FizzBuzz API

A small, layered Go REST API that serves FizzBuzz sequences and records every request in
ClickHouse so it can report usage statistics.

| Endpoint | Purpose |
| --- | --- |
| `GET /health` | Liveness probe |
| `GET /api/v1/fizzbuzz` | Generate a FizzBuzz sequence |
| `GET /api/v1/stats` | Usage analytics aggregated in ClickHouse |

Stack: [Go](https://go.dev) 1.26 · [chi](https://github.com/go-chi/chi) router ·
[ClickHouse](https://clickhouse.com) 25.8 · [golang-migrate](https://github.com/golang-migrate/migrate) ·
Docker Compose.

---

## Table of contents

- [Quick start](#quick-start)
- [Configuration](#configuration)
- [API reference](#api-reference)
- [Architecture](#architecture)
- [Project layout](#project-layout)
- [Database & migrations](#database--migrations)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [Maintenance conventions](#maintenance-conventions)

---

## Quick start

> Requires Docker with the Compose plugin. Building and testing locally also needs Go 1.26+.

```bash
# 1. Clone
git clone https://github.com/Abdillah-Epi/fizz-buzz.git
cd fizz-buzz

# 2. Create your environment file and fill in the values (see Configuration)
cp .env.example .env

# 3. Build and start the API + ClickHouse
docker compose up -d --build

# 4. Apply database migrations
make tools        # one-off: installs the golang-migrate CLI (with the clickhouse tag)
make migrate-up

# 5. Smoke test (replace 8080 with your API_PORT)
curl -s http://localhost:8080/health
curl -s 'http://localhost:8080/api/v1/fizzbuzz?int1=3&int2=5&limit=15&str1=Fizz&str2=Buzz'
curl -s http://localhost:8080/api/v1/stats
```

Stop the stack with `docker compose down` (add `-v` to also delete the ClickHouse volumes).

---

## Configuration

All configuration comes from environment variables, loaded by `internal/config` at startup.
Docker Compose reads them from `.env`; `make` reads the same file via `-include .env`.

Start from the template:

```bash
cp .env.example .env
```

Then fill it in with concrete values, for example:

```ini
API_PORT=8080
FIZZBUZZ_MAX_LIMIT=10000

CLICKHOUSE_DB=fizzbuzz
CLICKHOUSE_USER=fizzbuzz
CLICKHOUSE_PASSWORD=fizzbuzz
CLICKHOUSE_PORT=9000
CLICKHOUSE_HTTP_PORT=8123
CLICKHOUSE_MIGRATION_HOST=localhost
```

| Variable | Read by | Default | Purpose |
| --- | --- | --- | --- |
| `FIZZBUZZ_MAX_LIMIT` | API (`internal/config`) | `10000` | Upper bound accepted for the `limit` param |
| `CLICKHOUSE_MIGRATION_HOST` | Makefile migrations | `localhost` | ClickHouse host **as seen from your machine** |

---

## API reference

### `GET /health`

Returns `200 OK` with an empty body once the server is up.

### `GET /api/v1/fizzbuzz`

All five query parameters are **required**.

| Param | Type | Rules |
| --- | --- | --- |
| `int1` | int | `> 0` — the first divisor |
| `int2` | int | `> 0` — the second divisor |
| `limit` | int | `> 0` and `<= FIZZBUZZ_MAX_LIMIT` |
| `str1` | string | non-empty — printed for multiples of `int1` |
| `str2` | string | non-empty — printed for multiples of `int2` |

```bash
curl -s 'http://localhost:8080/api/v1/fizzbuzz?int1=3&int2=5&limit=15&str1=Fizz&str2=Buzz'
```

```json
{ "values": ["1","2","Fizz","4","Buzz","Fizz","7","8","Fizz","Buzz","11","Fizz","13","14","FizzBuzz"] }
```

Invalid input returns `400 Bad Request` with a plain-text message, for example
`invalid int1`, `limit must be greater than 0`, or
`limit must be less than or equal to 10000`.

Every valid request is recorded in `request_events`. **Analytics must never break the
core endpoint:** if the write fails it is logged and the sequence is still returned.

### `GET /api/v1/stats`

```bash
curl -s http://localhost:8080/api/v1/stats
```

```json
{
  "total_requests": 1255,
  "requests_last_hour": 1255,
  "requests_last_24h": 1255,
  "unique_configurations": 6,
  "most_used": {
    "int1": 15,
    "int2": 25,
    "limit": 350,
    "str1": "Fezz",
    "str2": "Byzz",
    "hits": 450,
    "share": 35.85657370517929,
    "first_seen_at": "2026-09-23T08:33:28Z",
    "last_seen_at": "2026-09-23T08:33:36Z"
  }
}
```

| Field | Computed in | Meaning |
| --- | --- | --- |
| `total_requests` | ClickHouse | Total rows in `request_events` |
| `requests_last_hour` | ClickHouse | Requests in the last hour |
| `requests_last_24h` | ClickHouse | Requests in the last 24 hours |
| `unique_configurations` | ClickHouse | Distinct `(int1,int2,limit,str1,str2)` tuples |
| `most_used` | ClickHouse | The most frequent tuple, with `hits` and first/last seen timestamps |
| `share` | Go | `hits / total_requests * 100`, derived in `internal/analytics/repository/clickhouse.go` |

When no events have been recorded yet, the endpoint returns `200 OK` with every field
zeroed (`hits`, `share`, and all counters are `0`). A `500 failed to get statistics`
response therefore means a genuine ClickHouse failure, not an empty table.

---

## Architecture

Requests flow through one direction only, from transport down to storage:

```
HTTP request
   │
   ▼
router  (internal/router)                       chi middleware + route table
   │
   ▼
handler (internal/*/handler)                    parse input, map DTO ⇄ model
   │
   ▼
service (internal/*/service)                    business rules / use cases
   │
   ▼
repository (internal/analytics/repository)      SQL against ClickHouse
   │
   ▼
infrastructure/clickhouse                       connection lifecycle (open/ping/close)
```

- Layers depend **downward only**; no layer imports a layer above it.
- Each feature keeps its own `dto` (wire format) and `model` (domain) types, so HTTP and
  storage shapes can change independently.
- All dependency wiring lives in one place: `internal/server/server.go`. It builds the
  ClickHouse client, repository, services and handlers, then hands them to `router.New`.
  Add an endpoint by adding a handler and a route file — not by editing `main.go`.

---

## Project layout

```
.
├── cmd/api/main.go                 entrypoint: config → server → graceful shutdown
├── internal/
│   ├── config/config.go            env parsing, validation, defaults
│   ├── server/server.go            composition root (dependency wiring)
│   ├── router/                     chi middleware + route registration
│   ├── infrastructure/clickhouse/  ClickHouse client (connect / ping / close)
│   ├── fizzbuzz/
│   │   ├── handler/                HTTP handler + query param parsing
│   │   ├── service/                generation + validation rules
│   │   ├── dto/                    JSON request/response shapes
│   │   └── model/                  domain types
│   └── analytics/
│       ├── handler/                /stats HTTP handler
│       ├── service/                analytics use cases
│       ├── repository/             SQL against request_events
│       ├── dto/                    JSON response shape
│       └── model/                  Stats, RequestEvent
├── migrations/clickhouse/          golang-migrate SQL migrations
├── script01.sh … script05.sh       sample load generators
├── stats.sh                        prints /stats only
├── docker-compose.yml / Dockerfile
├── Makefile                        migration helpers
└── .github/workflows/ci.yml        build → lint → test
```

---

## Database & migrations

The single table, `request_events`:

| Column | Type | Notes |
| --- | --- | --- |
| `timestamp` | `DateTime64(3, 'UTC')` | Event time, set by the API |
| `int1`, `int2` | `UInt32` | Request divisors |
| `limit_value` | `UInt32` | Request `limit` |
| `str1`, `str2` | `String` | Request strings |

Engine `MergeTree`, `ORDER BY (timestamp, int1, int2)`.

Migrations live in `migrations/clickhouse/` as paired files named
`<version>_<name>.up.sql` and `<version>_<name>.down.sql`, applied with golang-migrate
(the ClickHouse driver is behind a build tag, which `make tools` handles for you).

```bash
make tools           # install the migrate CLI (clickhouse tag)
make migrate-up      # apply all pending migrations
make migrate-down    # roll back the last migration
make migrate-version # show the current schema version
```

The Makefile builds its connection URL from `.env`, so migrations run from your machine
against the published native port (`CLICKHOUSE_PORT`).

**Adding a migration**

1. Create the next-numbered `NNN_description.up.sql` and the matching `.down.sql`.
2. Run `make migrate-up`, then confirm with `make migrate-version`.
3. Ship the up and down files in the same commit and keep them reversible.

---
