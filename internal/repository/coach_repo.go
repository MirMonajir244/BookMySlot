package repository

import (
	"github.com/MirMonajir244/BookMySlot/internal/models"
	"gorm.io/gorm"
)

type CoachRepository struct {
	db *gorm.DB
}

func NewCoachRepository(db *gorm.DB) *CoachRepository {
	return &CoachRepository{db: db}
}

func (r *CoachRepository) Create(coach *models.Coach) error {
	return r.db.Create(coach).Error
}

func (r *CoachRepository) FindByEmail(email string) (*models.Coach, error) {
	var coach models.Coach
	err := r.db.Where("email = ?", email).First(&coach).Error
	if err != nil {
		return nil, err
	}
	return &coach, nil
}

func (r *CoachRepository) FindByID(id uint) (*models.Coach, error) {
	var coach models.Coach
	err := r.db.First(&coach, id).Error
	if err != nil {
		return nil, err
	}
	return &coach, nil
}
