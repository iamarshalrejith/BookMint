# BookMint 🎬

A concurrent cinema ticket booking REST API built with **Go**, demonstrating how real-world systems prevent race conditions when multiple users try to book the same seat simultaneously.

## The Problem This Solves

Ever wondered what happens when two users click "Book Seat A3" at the exact same millisecond? Without proper concurrency handling, both could succeed — and now two people show up for the same seat. BookMint solves this using Go's concurrency primitives and Redis atomic operations.

## How It Works

BookMint uses a **two-phase booking lifecycle**:

```
User clicks seat → HOLD (2 min TTL) → User confirms → CONFIRMED (permanent)
                                    ↘ Time expires / User cancels → RELEASED
```

**The secret weapon against double-booking:** Redis `SET NX` (Set if Not Exists). This is an atomic operation — only one goroutine in the entire distributed system can win the race. The loser gets `ErrSeatAlreadyBooked` immediately.

## Tech Stack

- **Go 1.25** — HTTP server using the standard library (`net/http`)
- **Redis 7** — Distributed locking via atomic `SET NX` with TTL
- **Docker Compose** — One command to spin up Redis + Redis Commander (UI)
- **UUID** — Session ID generation per booking

## Architecture

```
cmd/
└── main.go                  # Entry point, route registration

internal/
├── booking/
│   ├── domain.go            # Core types: Booking, BookingStore interface, errors
│   ├── service.go           # Business logic layer
│   ├── handler.go           # HTTP handlers (hold, confirm, release, list)
│   ├── redis_store.go       # Redis-backed store (production)
│   ├── concurrent_store.go  # In-memory store with mutex (learning reference)
│   ├── memory_store.go      # Simple in-memory store
│   └── service_test.go      # Unit tests
├── adapters/redis/
│   └── redis.go             # Redis client setup
└── utils/
    └── utils.go             # Shared JSON helpers

static/
└── index.html               # Browser UI for seat selection
```

## API Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/movies` | List all available movies |
| `GET` | `/movies/{movieID}/seats` | List seat availability for a movie |
| `POST` | `/movies/{movieID}/seats/{seatID}/hold` | Hold a seat (returns session ID) |
| `PUT` | `/sessions/{sessionID}/confirm` | Confirm a held seat |
| `DELETE` | `/sessions/{sessionID}` | Release a held seat |

### Hold a Seat
```bash
curl -X POST http://localhost:8080/movies/inception/seats/A1/hold \
  -H "Content-Type: application/json" \
  -d '{"user_id": "alice"}'
```
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "movie_id": "inception",
  "seat_id": "A1",
  "expires_at": "2025-07-26T10:02:00Z"
}
```

### Confirm a Booking
```bash
curl -X PUT http://localhost:8080/sessions/{session_id}/confirm \
  -H "Content-Type: application/json" \
  -d '{"user_id": "alice"}'
```

### Release a Seat
```bash
curl -X DELETE http://localhost:8080/sessions/{session_id} \
  -H "Content-Type: application/json" \
  -d '{"user_id": "alice"}'
```

## Redis Key Design

```
seat:{movieID}:{seatID}   →  Booking JSON  [TTL=2min if held, no TTL if confirmed]
session:{sessionID}        →  seat key      [reverse lookup, TTL=2min]
```

The TTL on the seat key is what auto-releases abandoned holds — no cron job needed.

## Getting Started

**Prerequisites:** Go 1.25+, Docker

```bash
# 1. Clone the repo
git clone https://github.com/iamarshalrejith/BookMint.git
cd BookMint

# 2. Start Redis
docker compose up -d

# 3. Run the server
go run ./cmd/main.go

# 4. Open the UI
open http://localhost:8080
```

Redis Commander (visual Redis browser) is available at **http://localhost:8081**

## Concurrency: Two Layers

**Layer 1 — In-process (ConcurrentStore):** Uses `sync.RWMutex` to protect the in-memory map. Multiple goroutines can read simultaneously; writes get an exclusive lock.

**Layer 2 — Distributed (RedisStore):** Uses Redis `SET NX` for atomic seat claiming across any number of server instances. This is the production implementation.

## Running Tests

```bash
go test ./internal/booking/...
```

## Available Movies (Seeded)

| ID | Title | Layout |
|----|-------|--------|
| `inception` | Inception | 5 rows × 8 seats |
| `dune` | Dune: Part Two | 4 rows × 6 seats |

## What I Learned Building This

- How `SET NX` provides atomic distributed locking without a separate lock key
- Why TTLs are a cleaner solution than manual cleanup for temporary holds
- The difference between in-process mutex locking and distributed coordination
- Clean architecture with interface-driven design (`BookingStore` interface lets you swap Redis for memory in tests)
- Go's pattern of `context` propagation for cancellation and timeouts

---

Built with ❤️ by [iamarshalrejith](https://github.com/iamarshalrejith)
