package service

import (
	"errors"
	"time"

	"github.com/MirMonajir244/BookMySlot/internal/models"
	"github.com/MirMonajir244/BookMySlot/internal/repository"
	"gorm.io/gorm"
)

type BookingService struct {
	bookingRepo *repository.BookingRepository
	availRepo   *repository.AvailabilityRepository
	coachRepo   *repository.CoachRepository
	userRepo    *repository.UserRepository
}

func NewBookingService(
	bookingRepo *repository.BookingRepository,
	availRepo *repository.AvailabilityRepository,
	coachRepo *repository.CoachRepository,
	userRepo *repository.UserRepository,
) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		availRepo:   availRepo,
		coachRepo:   coachRepo,
		userRepo:    userRepo,
	}
}

// CreateBooking creates a new booking with concurrency protection.
// It validates the slot exists in the coach's availability and uses
// transaction + row-level locking to prevent double booking.
func (s *BookingService) CreateBooking(userID, coachID uint, dateTime time.Time, idempotencyKey string) (*models.Booking, error) {
	// Check idempotency — if a booking with this key already exists, return it
	if idempotencyKey != "" {
		existing, err := s.bookingRepo.FindByIdempotencyKey(idempotencyKey)
		if err == nil {
			return existing, nil // idempotent: return existing booking
		}
	}

	// Validate user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate coach exists
	_, err = s.coachRepo.FindByID(coachID)
	if err != nil {
		return nil, errors.New("coach not found")
	}

	// Ensure datetime is in UTC
	dateTime = dateTime.UTC()

	// Validate the slot exists within coach's availability
	if err := s.validateSlotAvailability(coachID, dateTime); err != nil {
		return nil, err
	}

	booking := &models.Booking{
		UserID:         userID,
		CoachID:        coachID,
		DateTime:       dateTime,
		Status:         models.BookingStatusConfirmed,
		IdempotencyKey: idempotencyKey,
	}

	// Create with transaction + row lock to prevent double booking
	if err := s.bookingRepo.CreateWithLock(booking); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("slot already booked")
		}
		return nil, errors.New("failed to create booking: " + err.Error())
	}

	return booking, nil
}

// validateSlotAvailability checks that the requested datetime falls within
// the coach's availability for that weekday.
func (s *BookingService) validateSlotAvailability(coachID uint, dateTime time.Time) error {
	dayOfWeek := dateTime.Weekday().String()
	availabilities, err := s.availRepo.FindByCoachAndDay(coachID, dayOfWeek)
	if err != nil {
		return err
	}

	requestedHour := dateTime.Hour()
	requestedMinute := dateTime.Minute()

	for _, avail := range availabilities {
		startTime, _ := time.Parse("15:04", avail.StartTime)
		endTime, _ := time.Parse("15:04", avail.EndTime)

		slotStart := requestedHour*60 + requestedMinute
		availStart := startTime.Hour()*60 + startTime.Minute()
		availEnd := endTime.Hour()*60 + endTime.Minute()

		// The slot start must be >= availability start
		// The slot end (start + 30min) must be <= availability end
		if slotStart >= availStart && (slotStart+30) <= availEnd {
			return nil // valid slot!
		}
	}

	return errors.New("requested time slot is not within coach's availability")
}

// GetUserBookings returns paginated bookings for a user.
func (s *BookingService) GetUserBookings(userID uint, offset, limit int) ([]models.Booking, int64, error) {
	return s.bookingRepo.FindByUser(userID, offset, limit)
}

// CancelBooking cancels a booking if it belongs to the requesting user.
func (s *BookingService) CancelBooking(bookingID, userID uint) error {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	if booking.UserID != userID {
		return errors.New("unauthorized: you can only cancel your own bookings")
	}

	if booking.Status == models.BookingStatusCancelled {
		return errors.New("booking is already cancelled")
	}

	return s.bookingRepo.Cancel(bookingID)
}
