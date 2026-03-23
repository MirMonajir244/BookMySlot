package repository

import (
	"github.com/MirMonajir244/BookMySlot/internal/models"
	"gorm.io/gorm"
)

type AvailabilityRepository struct {
	db *gorm.DB
}

func NewAvailabilityRepository(db *gorm.DB) *AvailabilityRepository {
	return &AvailabilityRepository{db: db}
}

func (r *AvailabilityRepository) Create(a *models.Availability) error {
	return r.db.Create(a).Error
}

func (r *AvailabilityRepository) FindByCoachAndDay(coachID uint, dayOfWeek string) ([]models.Availability, error) {
	var availabilities []models.Availability
	err := r.db.Where("coach_id = ? AND day_of_week = ?", coachID, dayOfWeek).
		Order("start_time ASC").
		Find(&availabilities).Error
	return availabilities, err
}

func (r *AvailabilityRepository) FindByCoach(coachID uint) ([]models.Availability, error) {
	var availabilities []models.Availability
	err := r.db.Where("coach_id = ?", coachID).
		Order("day_of_week ASC, start_time ASC").
		Find(&availabilities).Error
	return availabilities, err
}

func (r *AvailabilityRepository) Delete(id uint, coachID uint) error {
	return r.db.Where("id = ? AND coach_id = ?", id, coachID).Delete(&models.Availability{}).Error
}
