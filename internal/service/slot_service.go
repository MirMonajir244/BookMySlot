package service

import (
	"time"

	"github.com/MirMonajir244/BookMySlot/internal/dto"
	"github.com/MirMonajir244/BookMySlot/internal/models"
	"github.com/MirMonajir244/BookMySlot/internal/repository"
)

const SlotDuration = 30 * time.Minute

type SlotService struct {
	availRepo   *repository.AvailabilityRepository
	bookingRepo *repository.BookingRepository
}

func NewSlotService(availRepo *repository.AvailabilityRepository, bookingRepo *repository.BookingRepository) *SlotService {
	return &SlotService{
		availRepo:   availRepo,
		bookingRepo: bookingRepo,
	}
}

// GetAvailableSlots generates all available 30-minute slots for a coach on a specific date.
// It dynamically generates slots from the coach's weekly availability and subtracts already booked slots.
func (s *SlotService) GetAvailableSlots(coachID uint, date time.Time) ([]dto.SlotResponse, error) {
	// Determine the day of week from the date
	dayOfWeek := date.Weekday().String()

	// Get coach availability for that day
	availabilities, err := s.availRepo.FindByCoachAndDay(coachID, dayOfWeek)
	if err != nil {
		return nil, err
	}

	if len(availabilities) == 0 {
		return []dto.SlotResponse{}, nil
	}

	// Get existing bookings for the date range (start of day to end of day)
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	bookings, err := s.bookingRepo.FindByCoachAndTimeRange(coachID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	// Build a set of booked times for O(1) lookup
	bookedTimes := make(map[time.Time]bool)
	for _, b := range bookings {
		bookedTimes[b.DateTime] = true
	}

	// Generate slots
	var slots []dto.SlotResponse
	for _, avail := range availabilities {
		generated := GenerateSlots(avail, date, bookedTimes)
		slots = append(slots, generated...)
	}

	return slots, nil
}

// GenerateSlots creates 30-min slots from an availability block, excluding booked ones.
func GenerateSlots(avail models.Availability, date time.Time, bookedTimes map[time.Time]bool) []dto.SlotResponse {
	startTime, _ := time.Parse("15:04", avail.StartTime)
	endTime, _ := time.Parse("15:04", avail.EndTime)

	var slots []dto.SlotResponse

	current := time.Date(date.Year(), date.Month(), date.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
	end := time.Date(date.Year(), date.Month(), date.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

	for current.Add(SlotDuration).Before(end) || current.Add(SlotDuration).Equal(end) {
		slotEnd := current.Add(SlotDuration)
		if !bookedTimes[current] {
			slots = append(slots, dto.SlotResponse{
				StartTime: current,
				EndTime:   slotEnd,
			})
		}
		current = slotEnd
	}

	return slots
}
