package dto

import "time"

// ---- Auth ----

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=user coach"` // "user" or "coach"
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=user coach"`
}

// ---- Availability ----

type SetAvailabilityRequest struct {
	CoachID   uint   `json:"coach_id" binding:"required"`
	DayOfWeek string `json:"day_of_week" binding:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday Sunday"`
	StartTime string `json:"start_time" binding:"required"` // "HH:MM"
	EndTime   string `json:"end_time" binding:"required"`   // "HH:MM"
}

// ---- Booking ----

type CreateBookingRequest struct {
	UserID   uint      `json:"user_id" binding:"required"`
	CoachID  uint      `json:"coach_id" binding:"required"`
	DateTime time.Time `json:"datetime" binding:"required"`
}

// ---- Pagination ----

type PaginationQuery struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=10" binding:"min=1,max=100"`
}

func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}
