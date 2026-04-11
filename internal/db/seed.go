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
		{Name: "WiFi", NameFa: "وای‌فای"},
		{Name: "TV", NameFa: "تلویزیون"},
		{Name: "Air Conditioning", NameFa: "تهویه مطبوع"},
		{Name: "Mini Bar", NameFa: "مینی‌بار"},
		{Name: "Safe", NameFa: "صندوق امانات"},
		{Name: "Ocean View", NameFa: "نمای دریا"},
		{Name: "City View", NameFa: "نمای شهر"},
		{Name: "Balcony", NameFa: "بالکن"},
		{Name: "Jacuzzi", NameFa: "جکوزی"},
		{Name: "Room Service", NameFa: "سرویس اتاق"},
	})

	seed(db, []models.ParkingSpotType{
		{Name: "Standard", NameFa: "استاندارد"},
		{Name: "Handicap", NameFa: "مخصوص معلولین"},
		{Name: "Electric", NameFa: "شارژ خودرو برقی"},
		{Name: "Compact", NameFa: "خودرو کوچک"},
		{Name: "Large", NameFa: "خودرو بزرگ"},
	})

	seed(db, []models.ParkingSpotStatus{
		{Name: "Available", NameFa: "در دسترس"},
		{Name: "Occupied", NameFa: "اشغال شده"},
		{Name: "Reserved", NameFa: "رزرو شده"},
		{Name: "Maintenance", NameFa: "در حال تعمیر"},
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
