# BookMySlot API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All protected endpoints require a JWT token in the `Authorization` header:
```
Authorization: Bearer <jwt_token>
```

---

## Endpoints

### Auth

#### POST `/auth/register`
Register a new user or coach.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secure123",
  "role": "user"  // or "coach"
}
```

**Response (201):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user"
}
```

---

#### POST `/auth/login`
Login as a user or coach.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "secure123",
  "role": "user"
}
```

**Response (200):** Same as register response.

---

### Coach Endpoints (Requires `coach` role)

#### POST `/coaches/availability`
Set weekly availability.

**Request Body:**
```json
{
  "coach_id": 1,
  "day_of_week": "Monday",
  "start_time": "10:00",
  "end_time": "15:00"
}
```

**Response (201):**
```json
{
  "id": 1,
  "coach_id": 1,
  "day_of_week": "Monday",
  "start_time": "10:00",
  "end_time": "15:00"
}
```

---

#### GET `/coaches/availability`
Get the authenticated coach's availability.

**Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "coach_id": 1,
      "day_of_week": "Monday",
      "start_time": "10:00",
      "end_time": "15:00"
    }
  ]
}
```

---

### User Endpoints (Requires `user` role)

#### GET `/users/slots?coach_id=1&date=2025-10-28`
Get available 30-minute slots for a coach on a specific date.

**Query Parameters:**
| Param | Type | Required | Description |
|-------|------|----------|-------------|
| `coach_id` | int | Yes | Coach ID |
| `date` | string | Yes | Date (YYYY-MM-DD) |

**Response (200):**
```json
{
  "coach_id": 1,
  "date": "2025-10-28",
  "slots": [
    {
      "start_time": "2025-10-28T10:00:00Z",
      "end_time": "2025-10-28T10:30:00Z"
    },
    {
      "start_time": "2025-10-28T10:30:00Z",
      "end_time": "2025-10-28T11:00:00Z"
    }
  ],
  "count": 2
}
```

---

#### POST `/users/bookings`
Book an appointment slot.

**Headers:**
| Header | Required | Description |
|--------|----------|-------------|
| `Idempotency-Key` | No | Prevents duplicate bookings |

**Request Body:**
```json
{
  "user_id": 1,
  "coach_id": 1,
  "datetime": "2025-10-28T10:00:00Z"
}
```

**Response (201):**
```json
{
  "id": 1,
  "user_id": 1,
  "coach_id": 1,
  "datetime": "2025-10-28T10:00:00Z",
  "status": "confirmed",
  "created_at": "2025-10-25T12:00:00Z"
}
```

**Error (409 Conflict):** Slot already booked.

---

#### GET `/users/bookings?page=1&page_size=10`
Get authenticated user's bookings (paginated).

**Query Parameters:**
| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 10 | Items per page (max 100) |

**Response (200):**
```json
{
  "data": [...],
  "page": 1,
  "page_size": 10,
  "total_items": 25,
  "total_pages": 3
}
```

---

#### DELETE `/users/bookings/:id`
Cancel a booking.

**Response (200):**
```json
{
  "message": "booking cancelled successfully"
}
```

---

## Error Responses
All errors follow this format:
```json
{
  "error": "error_code",
  "message": "Human-readable description"
}
```

| Status | Error Code | Description |
|--------|-----------|-------------|
| 400 | `validation_error` | Invalid request body or parameters |
| 401 | `unauthorized` | Missing or invalid auth token |
| 403 | `forbidden` | Insufficient permissions |
| 404 | `not_found` | Resource not found |
| 409 | `booking_error` | Slot already booked (double booking) |
| 500 | `server_error` | Internal server error |
