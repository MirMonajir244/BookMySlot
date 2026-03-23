package models

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// Booking represents a booked 30-minute session.
// The unique index on (coach_id, datetime) prevents double-booking at the DB level.
type Booking struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"not null;index" json:"user_id"`
	CoachID        uint           `gorm:"not null;uniqueIndex:idx_coach_datetime" json:"coach_id"`
	DateTime       time.Time      `gorm:"not null;uniqueIndex:idx_coach_datetime" json:"datetime"`
	Status         BookingStatus  `gorm:"size:20;not null;default:'confirmed'" json:"status"`
	IdempotencyKey string         `gorm:"size:255;uniqueIndex" json:"-"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	User           User           `gorm:"foreignKey:UserID" json:"-"`
	Coach          Coach          `gorm:"foreignKey:CoachID" json:"-"`
}
