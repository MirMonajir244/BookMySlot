package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/MirMonajir244/BookMySlot/internal/models"
	"github.com/MirMonajir244/BookMySlot/internal/repository"
)

type AvailabilityService struct {
	availRepo *repository.AvailabilityRepository
	coachRepo *repository.CoachRepository
}

func NewAvailabilityService(availRepo *repository.AvailabilityRepository, coachRepo *repository.CoachRepository) *AvailabilityService {
	return &AvailabilityService{
		availRepo: availRepo,
		coachRepo: coachRepo,
	}
}

func (s *AvailabilityService) SetAvailability(coachID uint, dayOfWeek, startTime, endTime string) (*models.Availability, error) {
	// Validate coach exists
	_, err := s.coachRepo.FindByID(coachID)
	if err != nil {
		return nil, errors.New("coach not found")
	}

	// Validate time format and logic
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format, use HH:MM: %v", err)
	}
	end, err := time.Parse("15:04", endTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time format, use HH:MM: %v", err)
	}
	if !end.After(start) {
		return nil, errors.New("end_time must be after start_time")
	}
	// Minimum 30-minute block
	if end.Sub(start) < 30*time.Minute {
		return nil, errors.New("availability block must be at least 30 minutes")
	}

	avail := &models.Availability{
		CoachID:   coachID,
		DayOfWeek: dayOfWeek,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err := s.availRepo.Create(avail); err != nil {
		return nil, errors.New("availability already exists for this time slot")
	}

	return avail, nil
}

func (s *AvailabilityService) GetCoachAvailability(coachID uint) ([]models.Availability, error) {
	return s.availRepo.FindByCoach(coachID)
}
