package repository

import (
	"time"

	"github.com/MirMonajir244/BookMySlot/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// CreateWithLock creates a booking within a transaction using row-level locking
// to prevent race conditions on the same coach+datetime slot.
func (r *BookingRepository) CreateWithLock(booking *models.Booking) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Row-level lock: check if a confirmed booking already exists for this coach+datetime
		var existing models.Booking
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("coach_id = ? AND date_time = ? AND status = ?",
				booking.CoachID, booking.DateTime, models.BookingStatusConfirmed).
			First(&existing).Error

		if err == nil {
			// A confirmed booking already exists — double booking!
			return gorm.ErrDuplicatedKey
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}

		// No existing confirmed booking; create the new one
		return tx.Create(booking).Error
	})
}

// FindByIdempotencyKey returns a booking matching the given idempotency key.
func (r *BookingRepository) FindByIdempotencyKey(key string) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Where("idempotency_key = ?", key).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByCoachAndTimeRange returns confirmed bookings for a coach within a time range.
func (r *BookingRepository) FindByCoachAndTimeRange(coachID uint, start, end time.Time) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("coach_id = ? AND date_time >= ? AND date_time < ? AND status = ?",
		coachID, start, end, models.BookingStatusConfirmed).
		Find(&bookings).Error
	return bookings, err
}

// FindByUser returns bookings for a user with pagination.
func (r *BookingRepository) FindByUser(userID uint, offset, limit int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Where("user_id = ? AND status = ?", userID, models.BookingStatusConfirmed)

	err := query.Model(&models.Booking{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("date_time ASC").Offset(offset).Limit(limit).Find(&bookings).Error
	return bookings, total, err
}

// FindByID returns a booking by its ID.
func (r *BookingRepository) FindByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.First(&booking, id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// Cancel sets a booking's status to cancelled.
func (r *BookingRepository) Cancel(id uint) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).
		Update("status", models.BookingStatusCancelled).Error
}
