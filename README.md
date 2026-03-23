# 📅 BookMySlot

A production-ready RESTful API for an appointment booking system where users can book 30-minute sessions with coaches. Built with **Go**, **Gin**, **PostgreSQL**, and **GORM** following clean architecture principles.

## ✨ Features

- **Dynamic Slot Generation** — 30-minute slots generated from weekly coach availability
- **Concurrency-Safe Booking** — DB transactions + row-level locking (`SELECT ... FOR UPDATE`) + unique constraints
- **JWT Authentication** — Role-based access control (user/coach)
- **Idempotent Bookings** — Duplicate-safe via `Idempotency-Key` header
- **Booking Cancellation** — Soft cancellation that frees slots for rebooking
- **Pagination** — All list endpoints support `page` and `page_size`
- **UTC Timezone Handling** — All times stored and compared in UTC
- **Clean Architecture** — Controller → Service → Repository → DB
- **Docker Ready** — One-command setup with `docker-compose`

## 🏗️ Architecture

```
cmd/server/main.go          → Entry point, DI wiring
internal/
├── config/                 → Environment configuration
├── database/               → PostgreSQL connection & migrations
├── models/                 → GORM models (User, Coach, Availability, Booking)
├── repository/             → Data access layer
├── service/                → Business logic layer
├── handler/                → HTTP handlers (controllers)
├── dto/                    → Request/Response DTOs with validation
├── middleware/             → Auth (JWT) & error recovery
└── router/                 → Route definitions & grouping
tests/                      → Unit tests
docs/                       → API documentation
```

## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- PostgreSQL 16+
- Docker & Docker Compose (optional)

### Option 1: Docker (Recommended)

### 🛠️ Using Makefile (Recommended)

The easiest way to manage the project is using the provided `Makefile`.

```bash
# Start everything (Docker + App)
make docker-up

# Rebuild and start
make docker-rebuild

# Run tests
make test

# View logs
make docker-logs

# Stop everything
make docker-down
```

### Option 2: Docker Compose (Direct)

```bash
# Start everything
docker compose up --build
```

### Option 3: Local Setup (Manual)

```bash
# 1. Clone and enter the project
git clone https://github.com/MirMonajir244/BookMySlot.git
cd BookMySlot

# 2. Create PostgreSQL database
createdb bookmyslot

# 3. Configure environment
cp .env.example .env
# Edit .env with your database credentials

# 4. Run the server
go run cmd/server/main.go
```

## 🔌 API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/auth/register` | ❌ | Register user or coach |
| `POST` | `/api/v1/auth/login` | ❌ | Login, returns JWT |
| `POST` | `/api/v1/coaches/availability` | 🔒 Coach | Set weekly availability |
| `GET` | `/api/v1/coaches/availability` | 🔒 Coach | View own availability |
| `GET` | `/api/v1/users/slots` | 🔒 User | Get available slots |
| `POST` | `/api/v1/users/bookings` | 🔒 User | Book appointment |
| `GET` | `/api/v1/users/bookings` | 🔒 User | List bookings (paginated) |
| `DELETE` | `/api/v1/users/bookings/:id` | 🔒 User | Cancel booking |

📄 See [full API documentation](docs/api.md) for request/response examples.

### 🧪 Running Tests

```bash
make test
```

## 🧠 Design Decisions

1. **Dynamic Slot Generation**: Slots are not stored in the DB. They're computed on-the-fly from weekly availability schedules, making the system flexible for future slot duration changes.

2. **Double Booking Prevention** (3 layers):
   - **Application-level**: Availability validation before booking
   - **Row-level locking**: `SELECT ... FOR UPDATE` in a transaction
   - **DB constraint**: Unique index on `(coach_id, datetime)`

3. **Idempotency**: The `Idempotency-Key` header ensures that retried requests (e.g., network failures) don't create duplicate bookings.

4. **Soft Cancellation**: Bookings are cancelled by setting `status = 'cancelled'` rather than deleting, preserving audit history and freeing the slot.

5. **Clean Architecture**: Business logic lives in the service layer, independent of HTTP handlers and DB implementation.

## 📝 Assumptions

- Coach availability is recurring weekly (same schedule every week)
- All times are in UTCq
- A session is exactly 30 minutes
- Users and coaches register with separate roles
- Each coach can only set their own availability
- Users can only book/cancel their own appointments

## 📦 Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.23 |
| Framework | Gin |
| Database | PostgreSQL 16 |
| ORM | GORM |
| Auth | JWT (golang-jwt/jwt) |
| Containerization | Docker + Docker Compose |

## 🤖 AI Usage Disclosure

AI tools were used to assist in the development of this project. Key areas of AI assistance:
- Initial project scaffolding and boilerplate code
- Code structure following clean architecture patterns
- Unit test creation

All code has been reviewed, understood, and can be fully explained by the developer.

## 📄 Author

Mir Monajir