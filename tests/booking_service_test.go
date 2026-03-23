package tests

import (
	"testing"

	"github.com/MirMonajir244/BookMySlot/internal/models"
)

func TestBookingStatus_Constants(t *testing.T) {
	if models.BookingStatusConfirmed != "confirmed" {
		t.Errorf("expected 'confirmed', got '%s'", models.BookingStatusConfirmed)
	}
	if models.BookingStatusCancelled != "cancelled" {
		t.Errorf("expected 'cancelled', got '%s'", models.BookingStatusCancelled)
	}
}

func TestBookingModel_Fields(t *testing.T) {
	booking := models.Booking{
		UserID:         1,
		CoachID:        2,
		Status:         models.BookingStatusConfirmed,
		IdempotencyKey: "test-key-123",
	}

	if booking.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", booking.UserID)
	}
	if booking.CoachID != 2 {
		t.Errorf("expected CoachID 2, got %d", booking.CoachID)
	}
	if booking.Status != models.BookingStatusConfirmed {
		t.Errorf("expected confirmed status, got %s", booking.Status)
	}
	if booking.IdempotencyKey != "test-key-123" {
		t.Errorf("expected idempotency key 'test-key-123', got '%s'", booking.IdempotencyKey)
	}
}

func TestBookingStatus_Cancellation(t *testing.T) {
	booking := models.Booking{
		Status: models.BookingStatusConfirmed,
	}

	// Simulate cancellation
	booking.Status = models.BookingStatusCancelled

	if booking.Status != models.BookingStatusCancelled {
		t.Errorf("expected cancelled status after cancellation, got %s", booking.Status)
	}
}
