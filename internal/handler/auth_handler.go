package handler

import (
	"net/http"

	"github.com/MirMonajir244/BookMySlot/internal/dto"
	"github.com/MirMonajir244/BookMySlot/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user or coach
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "Registration details"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if req.Role == "coach" {
		coach, err := h.authService.RegisterCoach(req.Name, req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "registration_failed",
				Message: err.Error(),
			})
			return
		}
		// Auto-login after registration
		_, token, _ := h.authService.LoginCoach(req.Email, req.Password)
		c.JSON(http.StatusCreated, dto.AuthResponse{
			Token: token,
			ID:    coach.ID,
			Name:  coach.Name,
			Email: coach.Email,
			Role:  "coach",
		})
		return
	}

	user, err := h.authService.RegisterUser(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "registration_failed",
			Message: err.Error(),
		})
		return
	}
	_, token, _ := h.authService.LoginUser(req.Email, req.Password)
	c.JSON(http.StatusCreated, dto.AuthResponse{
		Token: token,
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  "user",
	})
}

// Login godoc
// @Summary Login as user or coach
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if req.Role == "coach" {
		coach, token, err := h.authService.LoginCoach(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "login_failed",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, dto.AuthResponse{
			Token: token,
			ID:    coach.ID,
			Name:  coach.Name,
			Email: coach.Email,
			Role:  "coach",
		})
		return
	}

	user, token, err := h.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "login_failed",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.AuthResponse{
		Token: token,
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  "user",
	})
}
