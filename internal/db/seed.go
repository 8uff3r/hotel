package db

import (
	"fmt"
	"hotel/backend/internal/models"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {

	if err := db.AutoMigrate(models.AllPtr()...); err != nil {
		panic(fmt.Sprintf("auto migrate: %s", err))
	}
	seed(db, []models.Amenity{
		{Name: "WiFi"},
		{Name: "TV"},
		{Name: "Air Conditioning"},
		{Name: "Mini Bar"},
		{Name: "Safe"},
		{Name: "Ocean View"},
		{Name: "City View"},
		{Name: "Balcony"},
		{Name: "Jacuzzi"},
		{Name: "Room Service"},
	})

	seed(db, []models.ParkingSpotType{
		{Name: "Standard"},
		{Name: "Handicap"},
		{Name: "Electric"},
		{Name: "Compact"},
		{Name: "Large"},
	})

	seed(db, []models.ParkingSpotStatus{
		{Name: "Available"},
		{Name: "Occupied"},
		{Name: "Reserved"},
		{Name: "Maintenance"},
	})
}

func seed[T any](db *gorm.DB, defaultValues []T) error {
	var count int64
	var model T

	if len(defaultValues) == 0 {
		return nil
	}

	if err := db.Model(&model).Count(&count).Error; err != nil {
		return fmt.Errorf("error while seeding the db: %w", err)
	}
	if count > 0 {
		return nil
	}
	db.CreateInBatches(defaultValues, 100)

	return nil
}
