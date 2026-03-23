package router

import (
	"github.com/MirMonajir244/BookMySlot/internal/handler"
	"github.com/MirMonajir244/BookMySlot/internal/middleware"
	"github.com/MirMonajir244/BookMySlot/internal/service"
	"github.com/gin-gonic/gin"
)

func Setup(
	authHandler *handler.AuthHandler,
	coachHandler *handler.CoachHandler,
	userHandler *handler.UserHandler,
	authService *service.AuthService,
) *gin.Engine {
	r := gin.Default()

	// Global middleware
	r.Use(middleware.ErrorHandler())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	// ---- Public routes (no auth required) ----
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// ---- Protected routes ----
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authService))

	// Coach routes
	coaches := protected.Group("/coaches")
	coaches.Use(middleware.RoleMiddleware("coach"))
	{
		coaches.POST("/availability", coachHandler.SetAvailability)
		coaches.GET("/availability", coachHandler.GetAvailability)
	}

	// User routes
	users := protected.Group("/users")
	users.Use(middleware.RoleMiddleware("user"))
	{
		users.GET("/slots", userHandler.GetAvailableSlots)
		users.POST("/bookings", userHandler.CreateBooking)
		users.GET("/bookings", userHandler.GetBookings)
		users.DELETE("/bookings/:id", userHandler.CancelBooking)
	}

	return r
}
