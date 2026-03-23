package dto

import "time"

// ---- Auth ----

type AuthResponse struct {
	Token string `json:"token"`
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// ---- Availability ----

type AvailabilityResponse struct {
	ID        uint   `json:"id"`
	CoachID   uint   `json:"coach_id"`
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// ---- Slot ----

type SlotResponse struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// ---- Booking ----

type BookingResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	CoachID   uint      `json:"coach_id"`
	DateTime  time.Time `json:"datetime"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ---- Pagination ----

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

// ---- Error ----

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
