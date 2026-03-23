package main

import (
	"fmt"
	"log"

	"github.com/MirMonajir244/BookMySlot/internal/config"
	"github.com/MirMonajir244/BookMySlot/internal/database"
	"github.com/MirMonajir244/BookMySlot/internal/handler"
	"github.com/MirMonajir244/BookMySlot/internal/repository"
	"github.com/MirMonajir244/BookMySlot/internal/router"
	"github.com/MirMonajir244/BookMySlot/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	coachRepo := repository.NewCoachRepository(db)
	availRepo := repository.NewAvailabilityRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, coachRepo, cfg)
	availService := service.NewAvailabilityService(availRepo, coachRepo)
	slotService := service.NewSlotService(availRepo, bookingRepo)
	bookingService := service.NewBookingService(bookingRepo, availRepo, coachRepo, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	coachHandler := handler.NewCoachHandler(availService)
	userHandler := handler.NewUserHandler(slotService, bookingService)

	// Setup router
	r := router.Setup(authHandler, coachHandler, userHandler, authService)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 BookMySlot server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
