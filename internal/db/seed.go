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
	if err := db.AutoMigrate(models.AllForDb()...); err != nil {
		panic(fmt.Sprintf("auto migrate: %s", err))
	}
	seedAmenities(db)
	seedParkingSpotStatuses(db)
	seedParkingSpotTypes(db)
	seedRoomStatuses(db)
	seedRoomTypes(db)
	seedRoles(db)
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
		{TranslateBase: models.TranslateBase{Name: "WiFi", Translation: t["WiFi"]}},
		{TranslateBase: models.TranslateBase{Name: "TV", Translation: t["TV"]}},
		{TranslateBase: models.TranslateBase{Name: "Air Conditioning", Translation: t["Air Conditioning"]}},
		{TranslateBase: models.TranslateBase{Name: "Mini Bar", Translation: t["Mini Bar"]}},
		{TranslateBase: models.TranslateBase{Name: "Safe", Translation: t["Safe"]}},
		{TranslateBase: models.TranslateBase{Name: "Ocean View", Translation: t["Ocean View"]}},
		{TranslateBase: models.TranslateBase{Name: "City View", Translation: t["City View"]}},
		{TranslateBase: models.TranslateBase{Name: "Balcony", Translation: t["Balcony"]}},
		{TranslateBase: models.TranslateBase{Name: "Jacuzzi", Translation: t["Jacuzzi"]}},
		{TranslateBase: models.TranslateBase{Name: "Room Service", Translation: t["Room Service"]}},
	}

	seed(db, amenities)
}

func seedRoles(db *gorm.DB) {
	t := Translations["Role"]

	amenities := []models.Role{
		{TranslateBase: models.TranslateBase{Name: "admin", Translation: t["admin"]}},
		{TranslateBase: models.TranslateBase{Name: "staff", Translation: t["staff"]}},
		{TranslateBase: models.TranslateBase{Name: "receptionist", Translation: t["receptionist"]}},
		{TranslateBase: models.TranslateBase{Name: "housekeeper", Translation: t["housekeeper"]}},
	}

	seed(db, amenities)
}

func seedParkingSpotTypes(db *gorm.DB) {
	t := Translations["ParkingSpotType"]

	types := []models.ParkingSpotType{
		{TranslateBase: models.TranslateBase{Name: "Standard", Translation: t["Standard"]}},
		{TranslateBase: models.TranslateBase{Name: "Handicap", Translation: t["Handicap"]}},
		{TranslateBase: models.TranslateBase{Name: "Electric", Translation: t["Electric"]}},
		{TranslateBase: models.TranslateBase{Name: "Compact", Translation: t["Compact"]}},
		{TranslateBase: models.TranslateBase{Name: "Large", Translation: t["Large"]}},
	}

	seed(db, types)
}

func seedParkingSpotStatuses(db *gorm.DB) {
	t := Translations["ParkingSpotStatus"]

	statuses := []models.ParkingSpotStatus{
		{TranslateBase: models.TranslateBase{Name: "Available", Translation: t["Available"]}},
		{TranslateBase: models.TranslateBase{Name: "Occupied", Translation: t["Occupied"]}},
		{TranslateBase: models.TranslateBase{Name: "Reserved", Translation: t["Reserved"]}},
		{TranslateBase: models.TranslateBase{Name: "Maintenance", Translation: t["Maintenance"]}},
	}

	seed(db, statuses)
}

func seedRoomStatuses(db *gorm.DB) {
	t := Translations["RoomStatus"]

	statuses := []models.RoomStatus{
		{TranslateBase: models.TranslateBase{Name: "Available", Translation: t["Available"]}},
		{TranslateBase: models.TranslateBase{Name: "Occupied", Translation: t["Occupied"]}},
		{TranslateBase: models.TranslateBase{Name: "Reserved", Translation: t["Reserved"]}},
		{TranslateBase: models.TranslateBase{Name: "Maintenance", Translation: t["Maintenance"]}},
	}

	seed(db, statuses)
}

func seedRoomTypes(db *gorm.DB) {
	t := Translations["RoomType"]

	statuses := []models.RoomType{
		{TranslateBase: models.TranslateBase{}},
		{TranslateBase: models.TranslateBase{Name: "Single", Translation: t["Single"]}},
		{TranslateBase: models.TranslateBase{Name: "Double", Translation: t["Double"]}},
		{TranslateBase: models.TranslateBase{Name: "Suite", Translation: t["Suite"]}},
		{TranslateBase: models.TranslateBase{Name: "Deluxe", Translation: t["Deluxe"]}},
	}

	seed(db, statuses)
}
