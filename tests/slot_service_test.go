package tests

import (
	"testing"
	"time"

	"github.com/MirMonajir244/BookMySlot/internal/dto"
	"github.com/MirMonajir244/BookMySlot/internal/models"
	"github.com/MirMonajir244/BookMySlot/internal/service"
)

func TestGenerateSlots_BasicBlock(t *testing.T) {
	avail := models.Availability{
		CoachID:   1,
		DayOfWeek: "Monday",
		StartTime: "10:00",
		EndTime:   "12:00",
	}
	date := time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC) // Monday
	booked := make(map[time.Time]bool)

	slots := service.GenerateSlots(avail, date, booked)

	// 10:00-12:00 = 4 slots of 30 min
	if len(slots) != 4 {
		t.Errorf("expected 4 slots, got %d", len(slots))
	}

	expected := []dto.SlotResponse{
		{
			StartTime: time.Date(2025, 10, 27, 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 10, 27, 10, 30, 0, 0, time.UTC),
		},
		{
			StartTime: time.Date(2025, 10, 27, 10, 30, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 10, 27, 11, 0, 0, 0, time.UTC),
		},
		{
			StartTime: time.Date(2025, 10, 27, 11, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 10, 27, 11, 30, 0, 0, time.UTC),
		},
		{
			StartTime: time.Date(2025, 10, 27, 11, 30, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 10, 27, 12, 0, 0, 0, time.UTC),
		},
	}

	for i, slot := range slots {
		if !slot.StartTime.Equal(expected[i].StartTime) || !slot.EndTime.Equal(expected[i].EndTime) {
			t.Errorf("slot %d: expected %v-%v, got %v-%v",
				i, expected[i].StartTime, expected[i].EndTime, slot.StartTime, slot.EndTime)
		}
	}
}

func TestGenerateSlots_WithBookedSlots(t *testing.T) {
	avail := models.Availability{
		CoachID:   1,
		DayOfWeek: "Monday",
		StartTime: "10:00",
		EndTime:   "12:00",
	}
	date := time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC)

	// Book the second slot (10:30)
	booked := map[time.Time]bool{
		time.Date(2025, 10, 27, 10, 30, 0, 0, time.UTC): true,
	}

	slots := service.GenerateSlots(avail, date, booked)

	// 4 total - 1 booked = 3 available
	if len(slots) != 3 {
		t.Errorf("expected 3 slots, got %d", len(slots))
	}

	// Verify the booked slot (10:30) is not in the results
	for _, slot := range slots {
		if slot.StartTime.Equal(time.Date(2025, 10, 27, 10, 30, 0, 0, time.UTC)) {
			t.Error("booked slot 10:30 should not appear in available slots")
		}
	}
}

func TestGenerateSlots_AllBooked(t *testing.T) {
	avail := models.Availability{
		CoachID:   1,
		DayOfWeek: "Monday",
		StartTime: "10:00",
		EndTime:   "11:00",
	}
	date := time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC)

	booked := map[time.Time]bool{
		time.Date(2025, 10, 27, 10, 0, 0, 0, time.UTC):  true,
		time.Date(2025, 10, 27, 10, 30, 0, 0, time.UTC): true,
	}

	slots := service.GenerateSlots(avail, date, booked)

	if len(slots) != 0 {
		t.Errorf("expected 0 slots when all booked, got %d", len(slots))
	}
}

func TestGenerateSlots_MinimumBlock(t *testing.T) {
	// Exactly 30 minutes = 1 slot
	avail := models.Availability{
		CoachID:   1,
		DayOfWeek: "Monday",
		StartTime: "10:00",
		EndTime:   "10:30",
	}
	date := time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC)
	booked := make(map[time.Time]bool)

	slots := service.GenerateSlots(avail, date, booked)

	if len(slots) != 1 {
		t.Errorf("expected 1 slot for 30-min block, got %d", len(slots))
	}
}

func TestGenerateSlots_LargeBlock(t *testing.T) {
	// 10:00-15:00 = 5 hours = 10 slots
	avail := models.Availability{
		CoachID:   1,
		DayOfWeek: "Monday",
		StartTime: "10:00",
		EndTime:   "15:00",
	}
	date := time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC)
	booked := make(map[time.Time]bool)

	slots := service.GenerateSlots(avail, date, booked)

	if len(slots) != 10 {
		t.Errorf("expected 10 slots for 5-hour block, got %d", len(slots))
	}
}

func TestGenerateSlots_NotDivisibleBy30(t *testing.T) {
	// 10:00-11:15 = 75 minutes = only 2 full 30-min slots
	avail := models.Availability{
		CoachID:   1,
		DayOfWeek: "Monday",
		StartTime: "10:00",
		EndTime:   "11:15",
	}
	date := time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC)
	booked := make(map[time.Time]bool)

	slots := service.GenerateSlots(avail, date, booked)

	if len(slots) != 2 {
		t.Errorf("expected 2 slots for 75-min block, got %d", len(slots))
	}
}
