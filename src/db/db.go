package db

import (
	"fmt"
	"github.com/legenhand/recipe-test/src/config"
	"github.com/legenhand/recipe-test/src/model"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// SetupGormDB initializes the PostgreSQL database using GORM.
func SetupGormDB() *gorm.DB {
	cfg := config.Cfg

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Unit{}, &model.UnitConversion{}, &model.Inventory{},
		&model.Recipe{}, &model.RecipeIngredient{}); err != nil {
		return fmt.Errorf("failed to auto-migrate models: %w", err)
	}
	log.Println("Database migration completed successfully using GORM.")
	return nil
}
