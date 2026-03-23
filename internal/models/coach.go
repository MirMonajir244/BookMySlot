package models

import (
	"time"

	"gorm.io/gorm"
)

type Coach struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	Email          string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash   string         `gorm:"size:255;not null" json:"-"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Availabilities []Availability `gorm:"foreignKey:CoachID" json:"availabilities,omitempty"`
	Bookings       []Booking      `gorm:"foreignKey:CoachID" json:"bookings,omitempty"`
}
