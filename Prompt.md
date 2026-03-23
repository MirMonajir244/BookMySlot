### 🧠 Objective

Design and implement a scalable, production-ready RESTful API for an appointment booking system where users can book 30-minute sessions with coaches. The system should ensure correctness, prevent double booking, and be designed with clean architecture and extensibility in mind.

---

### 🏗️ System Design Expectations

Please approach this problem with a **backend engineer mindset**, covering:

1. **Data Modeling**

   * Design relational database schema (PostgreSQL preferred)
   * Entities should include (but not limited to):

     * Users
     * Coaches
     * Availability (recurring weekly schedule)
     * Bookings
   * Clearly define relationships, constraints, and indexes

2. **Core Logic**

   * Generate 30-minute slots dynamically from coach availability
   * Exclude already booked slots
   * Ensure no double booking using strong consistency guarantees

3. **Concurrency Handling**

   * Prevent race conditions when multiple users try to book the same slot
   * Use techniques such as:

     * DB transactions
     * Row-level locking (`SELECT ... FOR UPDATE`)
     * Unique constraints (coach_id + datetime)

---

### 🔌 API Design (RESTful & Clean)

Implement the following endpoints with proper validation and error handling:

#### 1. Set Coach Availability

`POST /coaches/availability`

* Define weekly recurring availability
* Input:

```json
{
  "coach_id": 1,
  "day_of_week": "Monday",
  "start_time": "10:00",
  "end_time": "15:00"
}
```

---

#### 2. Get Available Slots

`GET /users/slots?coach_id=1&date=2025-10-28`

* Return all **available 30-minute slots**
* Must:

  * Map weekday → availability
  * Subtract already booked slots

---

#### 3. Book Appointment

`POST /users/bookings`

```json
{
  "user_id": 101,
  "coach_id": 1,
  "datetime": "2025-10-28T09:30:00Z"
}
```

* Must:

  * Validate slot exists in availability
  * Prevent double booking
  * Use transaction

---

#### 4. Get User Bookings

`GET /users/bookings?user_id=101`

* Return all upcoming bookings

---

### ⭐ Bonus (Strong Signal for Interview)

* Implement booking cancellation
* Add timezone handling (store everything in UTC)
* Add pagination for APIs
* Add basic authentication (JWT or mock)
* Write unit tests for:

  * Slot generation
  * Booking logic
* Add idempotency for booking API

---

### 🧩 Architecture Expectations

* Use clean architecture / layered design:

  * Controller → Service → Repository → DB
* Keep business logic separate from handlers
* Use DTOs and validation layer
* Add meaningful error responses

---

### ⚙️ Tech Stack (Preferred)

* Language: Golang (preferred) / Node.js / Python
* Framework: Gin (Go) / Express / FastAPI
* Database: PostgreSQL
* ORM: Optional (GORM / Prisma / SQLAlchemy)

---

### 📄 Deliverables

1. GitHub repository with clean commit history
2. README including:

   * Setup steps
   * Design decisions
   * Assumptions
3. API documentation:

   * Swagger / Postman / Markdown
4. (Optional) Docker setup

---

### 🧠 Evaluation Focus

* Correctness of booking logic
* Data modeling decisions
* Concurrency handling
* Code quality and structure
* API design clarity

---

### ⚡ Important Notes

* Prioritize correctness over premature optimization
* Avoid hardcoding slots — generate dynamically
* Ensure system is extendable (e.g., different slot durations in future)

---

### 🎯 Goal

The final system should simulate a **real-world booking platform backend** with production-level thinking around consistency, concurrency, and clean code.
