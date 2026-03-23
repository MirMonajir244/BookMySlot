package handler

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/MirMonajir244/BookMySlot/internal/dto"
	"github.com/MirMonajir244/BookMySlot/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	slotService    *service.SlotService
	bookingService *service.BookingService
}

func NewUserHandler(slotService *service.SlotService, bookingService *service.BookingService) *UserHandler {
	return &UserHandler{
		slotService:    slotService,
		bookingService: bookingService,
	}
}

// GetAvailableSlots godoc
// @Summary Get available slots for a coach on a date
// @Tags users
// @Produce json
// @Param coach_id query int true "Coach ID"
// @Param date query string true "Date in YYYY-MM-DD format"
// @Success 200 {array} dto.SlotResponse
// @Security BearerAuth
// @Router /api/v1/users/slots [get]
func (h *UserHandler) GetAvailableSlots(c *gin.Context) {
	coachIDStr := c.Query("coach_id")
	dateStr := c.Query("date")

	if coachIDStr == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "coach_id and date are required query parameters",
		})
		return
	}

	coachID, err := strconv.ParseUint(coachIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "invalid coach_id",
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "invalid date format, use YYYY-MM-DD",
		})
		return
	}

	slots, err := h.slotService.GetAvailableSlots(uint(coachID), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "server_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coach_id": coachID,
		"date":     dateStr,
		"slots":    slots,
		"count":    len(slots),
	})
}

// CreateBooking godoc
// @Summary Book an appointment slot
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.CreateBookingRequest true "Booking details"
// @Param Idempotency-Key header string false "Idempotency key"
// @Success 201 {object} dto.BookingResponse
// @Failure 400,409 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/bookings [post]
func (h *UserHandler) CreateBooking(c *gin.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Ensure the authenticated user matches the request
	userID, _ := c.Get("user_id")
	if req.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "forbidden",
			Message: "you can only create bookings for yourself",
		})
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")

	booking, err := h.bookingService.CreateBooking(req.UserID, req.CoachID, req.DateTime, idempotencyKey)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "slot already booked" {
			status = http.StatusConflict
		}
		c.JSON(status, dto.ErrorResponse{
			Error:   "booking_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.BookingResponse{
		ID:        booking.ID,
		UserID:    booking.UserID,
		CoachID:   booking.CoachID,
		DateTime:  booking.DateTime,
		Status:    string(booking.Status),
		CreatedAt: booking.CreatedAt,
	})
}

// GetBookings godoc
// @Summary Get user's bookings (paginated)
// @Tags users
// @Produce json
// @Param user_id query int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Security BearerAuth
// @Router /api/v1/users/bookings [get]
func (h *UserHandler) GetBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination.Page = 1
		pagination.PageSize = 10
	}

	bookings, total, err := h.bookingService.GetUserBookings(userID.(uint), pagination.Offset(), pagination.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "server_error",
			Message: err.Error(),
		})
		return
	}

	var response []dto.BookingResponse
	for _, b := range bookings {
		response = append(response, dto.BookingResponse{
			ID:        b.ID,
			UserID:    b.UserID,
			CoachID:   b.CoachID,
			DateTime:  b.DateTime,
			Status:    string(b.Status),
			CreatedAt: b.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:       response,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(pagination.PageSize))),
	})
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Tags users
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/users/bookings/{id} [delete]
func (h *UserHandler) CancelBooking(c *gin.Context) {
	bookingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "invalid booking id",
		})
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.bookingService.CancelBooking(uint(bookingID), userID.(uint)); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "booking not found" {
			status = http.StatusNotFound
		} else if err.Error() == "unauthorized: you can only cancel your own bookings" {
			status = http.StatusForbidden
		}
		c.JSON(status, dto.ErrorResponse{
			Error:   "cancellation_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "booking cancelled successfully",
	})
}
