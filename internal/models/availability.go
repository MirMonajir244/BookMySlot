package models

import (
	"time"

	"gorm.io/gorm"
)

// Availability represents a recurring weekly time block for a coach.
// Times are stored as strings in "HH:MM" format and represent UTC.
type Availability struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CoachID   uint           `gorm:"not null;uniqueIndex:idx_coach_day_start" json:"coach_id"`
	DayOfWeek string         `gorm:"size:10;not null;uniqueIndex:idx_coach_day_start" json:"day_of_week"` // e.g. "Monday"
	StartTime string         `gorm:"size:5;not null;uniqueIndex:idx_coach_day_start" json:"start_time"`   // "HH:MM"
	EndTime   string         `gorm:"size:5;not null" json:"end_time"`                                     // "HH:MM"
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Coach     Coach          `gorm:"foreignKey:CoachID" json:"-"`
}
