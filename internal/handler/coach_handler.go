package handler

import (
	"net/http"

	"github.com/MirMonajir244/BookMySlot/internal/dto"
	"github.com/MirMonajir244/BookMySlot/internal/service"
	"github.com/gin-gonic/gin"
)

type CoachHandler struct {
	availService *service.AvailabilityService
}

func NewCoachHandler(availService *service.AvailabilityService) *CoachHandler {
	return &CoachHandler{availService: availService}
}

// SetAvailability godoc
// @Summary Set coach weekly availability
// @Tags coaches
// @Accept json
// @Produce json
// @Param body body dto.SetAvailabilityRequest true "Availability details"
// @Success 201 {object} dto.AvailabilityResponse
// @Failure 400 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/coaches/availability [post]
func (h *CoachHandler) SetAvailability(c *gin.Context) {
	var req dto.SetAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Ensure the authenticated coach matches the request
	coachID, _ := c.Get("user_id")
	if req.CoachID != coachID.(uint) {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "forbidden",
			Message: "you can only set your own availability",
		})
		return
	}

	avail, err := h.availService.SetAvailability(req.CoachID, req.DayOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "availability_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.AvailabilityResponse{
		ID:        avail.ID,
		CoachID:   avail.CoachID,
		DayOfWeek: avail.DayOfWeek,
		StartTime: avail.StartTime,
		EndTime:   avail.EndTime,
	})
}

// GetAvailability godoc
// @Summary Get coach availability
// @Tags coaches
// @Produce json
// @Param coach_id query uint true "Coach ID"
// @Success 200 {array} dto.AvailabilityResponse
// @Security BearerAuth
// @Router /api/v1/coaches/availability [get]
func (h *CoachHandler) GetAvailability(c *gin.Context) {
	coachID, _ := c.Get("user_id")

	availabilities, err := h.availService.GetCoachAvailability(coachID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "server_error",
			Message: err.Error(),
		})
		return
	}

	var response []dto.AvailabilityResponse
	for _, a := range availabilities {
		response = append(response, dto.AvailabilityResponse{
			ID:        a.ID,
			CoachID:   a.CoachID,
			DayOfWeek: a.DayOfWeek,
			StartTime: a.StartTime,
			EndTime:   a.EndTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
