### 🧠 Objective

Design and build a RESTful API for a simple appointment booking system. The platform connects coaches with users who want to book 30-minute appointment slots. This task will assess ability to design a system, model data, build APIs, and handle core business logic.

---

### 📋 Core Requirements

#### Functional Requirements

1. **Coach Availability**: A coach must be able to define their available hours for any day of the week. For example, Coach A is available on Mondays from 10:00 AM to 3:00 PM and Wednesdays from 9:00 AM to 12:00 PM.
2. **View Available Slots**: A user should be able to query the system to see all available 30-minute booking slots for a specific coach on a given day.
3. **Book an Appointment**: A user must be able to book one of the available 30-minute slots with a coach.
4. **No Double Booking**: Once a slot is booked, it must no longer appear as available. The system should prevent the same slot from being booked by two different users.
5. **View Booked Appointments**: A user should be able to see a list of all the appointments they have booked.

---

### 🔌 API Endpoints

#### 1. Set Coach Availability

`POST /coaches/availability`

* Action: Allows a coach to set their weekly availability.
* Request Body:

```json
{
  "coach_id": 1,
  "day_of_week": "Tuesday",
  "start_time": "09:00",
  "end_time": "14:00"
}
```

---

#### 2. Get Available Slots

`GET /users/slots?coach_id=1&date=2025-10-28`

* Action: Fetches all available 30-minute slots for a given coach on a specific day.
* Return: Array of available time slots

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

* Return all upcoming bookings with details about the coach and time.

---

### ⭐ Bonus Features

* **Cancellation**: Add an endpoint for a user to cancel a booking.
* **Timezones**: Handle timezones gracefully. Assume the coach and user might be in different timezones. Store everything in UTC.
* **Concurrency**: Handle two users trying to book the exact same slot at the exact same time. Implement using database transactions with row-level locking (`SELECT ... FOR UPDATE`) and unique constraints.
* **Testing**: Write unit or integration tests for API endpoints and business logic.
* **Pagination**: Add pagination for list APIs.
* **Authentication**: Add basic authentication (JWT or mock).
* **Idempotency**: Add idempotency for booking API.

---

### 🏗️ Architecture Expectations

* Use clean architecture / layered design:

  * Controller → Service → Repository → DB
* Keep business logic separate from handlers
* Use DTOs and validation layer
* Add meaningful error responses
* Code should be well-structured, clean, and maintainable

---

### ⚙️ Tech Stack

* Language: Golang
* Framework: Gin
* Database: PostgreSQL
* ORM: GORM

---

### 📄 Deliverables

1. GitHub repository with clean commit history
2. README including:

   * Design choices and assumptions
   * Step-by-step setup instructions (installing dependencies, database setup, starting server)
3. API documentation:
   * Swagger / Postman / Markdown
4. Docker setup
5. Makefile for common tasks

---

### 🧠 Evaluation Criteria

* **Correctness**: Does the application meet all core functional requirements?
* **Problem Solving**: How was data modeled? How are available slots generated based on coach's schedule and existing bookings?
* **Code Quality**: Is the code clean, well-organized, and easy to understand?
* **API Design**: Are the APIs well-designed, RESTful, and consistent? Is input validation and error handling implemented correctly?
* **Documentation**: Is the project easy to understand, set up, and use?

---

### ⚡ Important Notes

* Prioritize correctness over premature optimization
* Avoid hardcoding slots — generate dynamically
* Ensure system is extendable (e.g., different slot durations in future)

---

### 🤖 AI Usage Disclosure

AI tools were used during development. The complete conversation log is available upon request.

---

### 🎯 Goal

The final system should simulate a **real-world booking platform backend** with production-level thinking around consistency, concurrency, and clean code.
