package db

import (
	"encoding/json"
	"fmt"
	"hotel/internal/models"
	"log"

	_ "embed"

	"gorm.io/gorm"
)

//go:embed translations.json
var translationsFile []byte

var Translations map[string]map[string]models.Translation

func init() {
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
	t := Translations["Amenity"]

	amenities := []models.Amenity{
		{Name: "WiFi", Translation: t["WiFi"]},
		{Name: "TV", Translation: t["TV"]},
		{Name: "Air Conditioning", Translation: t["Air Conditioning"]},
		{Name: "Mini Bar", Translation: t["Mini Bar"]},
		{Name: "Safe", Translation: t["Safe"]},
		{Name: "Ocean View", Translation: t["Ocean View"]},
		{Name: "City View", Translation: t["City View"]},
		{Name: "Balcony", Translation: t["Balcony"]},
		{Name: "Jacuzzi", Translation: t["Jacuzzi"]},
		{Name: "Room Service", Translation: t["Room Service"]},
	}

	seed(db, amenities)
}

func seedParkingSpotTypes(db *gorm.DB) {
	t := Translations["ParkingSpotType"]

	types := []models.ParkingSpotType{
		{Name: "Standard", Translation: t["Standard"]},
		{Name: "Handicap", Translation: t["Handicap"]},
		{Name: "Electric", Translation: t["Electric"]},
		{Name: "Compact", Translation: t["Compact"]},
		{Name: "Large", Translation: t["Large"]},
	}

	seed(db, types)
}

func seedParkingSpotStatuses(db *gorm.DB) {
	t := Translations["ParkingSpotStatus"]

	statuses := []models.ParkingSpotStatus{
		{Name: "Available", Translation: t["Available"]},
		{Name: "Occupied", Translation: t["Occupied"]},
		{Name: "Reserved", Translation: t["Reserved"]},
		{Name: "Maintenance", Translation: t["Maintenance"]},
	}

	seed(db, statuses)
}
