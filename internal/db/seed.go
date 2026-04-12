package db

import (
	"encoding/json"
	"fmt"
	"hotel/backend/internal/models"
	"log"

	_ "embed"

	"gorm.io/gorm"
)

//go:embed translations.json
var translationsFile []byte

var Translations map[string]map[string]models.Translation

func init() {
	println(len(translationsFile))
	err := json.Unmarshal(translationsFile, &Translations)
	if err != nil {
		log.Fatalf("failed to load translations: %v", err)
	}
}
func Seed(db *gorm.DB) {
	if err := db.AutoMigrate(models.AllPtr()...); err != nil {
		panic(fmt.Sprintf("auto migrate: %s", err))
	}
	seedAmenities(db)
	seedParkingSpotStatuses(db)
	seedParkingSpotTypes(db)
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

func seedAmenities(db *gorm.DB) {
	amenityTranslations := Translations["Amenity"]

	amenities := []models.Amenity{
		{Name: "WiFi", Translation: amenityTranslations["WiFi"]},
		{Name: "TV", Translation: amenityTranslations["TV"]},
		{Name: "Air Conditioning", Translation: amenityTranslations["Air Conditioning"]},
		{Name: "Mini Bar", Translation: amenityTranslations["Mini Bar"]},
		{Name: "Safe", Translation: amenityTranslations["Safe"]},
		{Name: "Ocean View", Translation: amenityTranslations["Ocean View"]},
		{Name: "City View", Translation: amenityTranslations["City View"]},
		{Name: "Balcony", Translation: amenityTranslations["Balcony"]},
		{Name: "Jacuzzi", Translation: amenityTranslations["Jacuzzi"]},
		{Name: "Room Service", Translation: amenityTranslations["Room Service"]},
	}

	seed(db, amenities)
}

func seedParkingSpotTypes(db *gorm.DB) {
	parkingTypeTranslations := Translations["ParkingSpotType"]

	types := []models.ParkingSpotType{
		{Name: "Standard", Translation: parkingTypeTranslations["Standard"]},
		{Name: "Handicap", Translation: parkingTypeTranslations["Handicap"]},
		{Name: "Electric", Translation: parkingTypeTranslations["Electric"]},
		{Name: "Compact", Translation: parkingTypeTranslations["Compact"]},
		{Name: "Large", Translation: parkingTypeTranslations["Large"]},
	}

	seed(db, types)
}

func seedParkingSpotStatuses(db *gorm.DB) {
	statusTranslations := Translations["ParkingSpotStatus"]

	statuses := []models.ParkingSpotStatus{
		{Name: "Available", Translation: statusTranslations["Available"]},
		{Name: "Occupied", Translation: statusTranslations["Occupied"]},
		{Name: "Reserved", Translation: statusTranslations["Reserved"]},
		{Name: "Maintenance", Translation: statusTranslations["Maintenance"]},
	}

	seed(db, statuses)
}
