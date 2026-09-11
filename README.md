# Shipment gRPC Microservice

A gRPC microservice for managing shipments and tracking status changes in a Transportation Management System (TMS).

---

## Architecture

The project follows **Clean Architecture** principles with clear separation of concerns:
```
internal/
├── domain/           # Business entities, rules, repository interfaces
├── application/      # Use cases (orchestrates domain logic)
├── infrastructure/   # PostgreSQL repository implementations
└── transport/grpc/   # gRPC handlers and proto mappers
```

### Layer responsibilities

- **Domain** — Shipment entity, status transitions, business rules. No external dependencies.
- **Application** — Use cases: create shipment, get shipment, add event, get history.
- **Infrastructure** — PostgreSQL implementations of domain repository interfaces.
- **Messaging** — Kafka publisher for shipment status events.
- **Transport** — gRPC handlers that map proto types to domain types and call application services.

---

## Shipment Status Lifecycle

### Successful Flow
```
PENDING → PICKED_UP → IN_TRANSIT → DELIVERED
   ↓           ↓            ↓
CANCELLED  CANCELLED   CANCELLED
```

### Kafka Publish Failure
```
PENDING → PICKED_UP → IN_TRANSIT → DELIVERED
              │
              ├─ Shipment update saved to PostgreSQL
              ├─ Event saved to PostgreSQL
              ├─ Kafka publish fails
              └─ gRPC returns INTERNAL error
```

Rules:
- Every shipment starts with `pending` status
- Status can only move forward — no going back
- `delivered` and `cancelled` are final states
- Invalid transitions are rejected with gRPC status code `INVALID_ARGUMENT`
- Shipment and event data are saved to PostgreSQL before the event is published to Kafka
- If Kafka is unavailable, the database changes remain saved but the request returns an `INTERNAL` error

---

## Design Decisions

- **Clean Architecture** was chosen to keep business logic independent from gRPC and PostgreSQL
- **Repository interfaces** are defined in the domain layer — infrastructure implements them
- **Status transitions** are enforced in the domain entity, not in the database or handler
- **Domain errors** are mapped to gRPC status codes in the transport layer
- **UUID** is used for all entity IDs

---

## Assumptions

- `reference_number` must be unique per shipment
- `amount` and `driver_revenue` cannot be negative
- Status history is append-only — events are never deleted or updated
- Initial `pending` event is automatically recorded on shipment creation

---

## Prerequisites

- Go 1.22+
- Docker and Docker Compose

---

## How to Run the Service

### With Docker (recommended)
```bash
docker-compose up --build
```

The gRPC server will be available at `localhost:50051`.
PostgreSQL migrations run automatically on first start.

### Without Docker

Start PostgreSQL manually, then:
```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/shipments?sslmode=disable
export GRPC_PORT=50051
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=shipment-events

go run ./cmd/server
```

---

## How to Run Tests
```bash
go test ./internal/domain/... -v
```

Tests cover:
- Shipment creation and validation
- Valid status transitions
- Invalid status transitions
- Boundary cases (final states, negative values)

---

## gRPC Methods

| Method | Description |
|---|---|
| `CreateShipment` | Create a new shipment (starts with `pending`) |
| `GetShipment` | Get shipment details by ID |
| `AddShipmentEvent` | Add a status update (validated against lifecycle rules) |
| `GetShipmentHistory` | Get full event history for a shipment |

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/shipments?sslmode=disable` | PostgreSQL connection string |
| `GRPC_PORT` | `50051` | gRPC server port |
| `KAFKA_BROKERS` | `localhost:9092` | Comma-separated Kafka broker addresses |
| `KAFKA_TOPIC` | `shipment-events` | Topic for shipment events |

---

## Project Structure
```
shipment-service/
├── cmd/
│   └── server/              # Entry point (main.go)
├── gen/                     # Generated gRPC code (from proto)
│   ├── shipment.pb.go
│   └── shipment_grpc.pb.go
├── internal/
│   ├── domain/              # Core business logic
│   │   ├── shipment.go
│   │   ├── repository.go
│   │   └── shipment_test.go
│   ├── application/         # Use cases
│   │   └── shipment_service.go
│   ├── infrastructure/      # PostgreSQL repositories
│   │   └── postgres/
│   │       ├── shipment_repo.go
│   │       └── event_repo.go
│   └── transport/           # gRPC handlers
│       └── grpc/
│           ├── handler.go
│           └── mapper.go
├── migrations/
│   └── 001_create_shipments.sql
├── proto/
│   └── shipment.proto
├── Dockerfile
├── docker-compose.yml
└── README.md
```