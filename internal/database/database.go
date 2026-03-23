package database

import (
	"log"

	"github.com/MirMonajir244/BookMySlot/internal/config"
	"github.com/MirMonajir244/BookMySlot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Println("✅ Connected to PostgreSQL")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Coach{},
		&models.Availability{},
		&models.Booking{},
	)
	if err != nil {
		return err
	}

	log.Println("✅ Database migration completed")
	return nil
}
