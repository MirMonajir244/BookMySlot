package tests

import (
	"github.com/MirMonajir244/BookMySlot/internal/config"
	"github.com/MirMonajir244/BookMySlot/internal/handler"
	"github.com/MirMonajir244/BookMySlot/internal/models"
	"github.com/MirMonajir244/BookMySlot/internal/repository"
	"github.com/MirMonajir244/BookMySlot/internal/router"
	"github.com/MirMonajir244/BookMySlot/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestRouter initializes a full test environment with an in-memory SQLite database
func SetupTestRouter() (*gin.Engine, *gorm.DB, *config.Config) {
	gin.SetMode(gin.TestMode)

	// 1. Setup in-memory SQLite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 2. Run migrations
	db.AutoMigrate(&models.User{}, &models.Coach{}, &models.Availability{}, &models.Booking{})

	// 3. Setup Config
	cfg := &config.Config{
		JWTSecret: "test-secret",
		JWTExpiry: 24,
	}

	// 4. Setup Repositories
	userRepo := repository.NewUserRepository(db)
	coachRepo := repository.NewCoachRepository(db)
	availRepo := repository.NewAvailabilityRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// 5. Setup Services
	authService := service.NewAuthService(userRepo, coachRepo, cfg)
	availService := service.NewAvailabilityService(availRepo, coachRepo)
	slotService := service.NewSlotService(availRepo, bookingRepo)
	bookingService := service.NewBookingService(bookingRepo, availRepo, coachRepo, userRepo)

	// 6. Setup Handlers
	authHandler := handler.NewAuthHandler(authService)
	coachHandler := handler.NewCoachHandler(availService)
	userHandler := handler.NewUserHandler(slotService, bookingService)

	// 7. Setup Router
	r := router.Setup(authHandler, coachHandler, userHandler, authService)

	return r, db, cfg
}
