package seed

import (
	"encoding/json"
	"fmt"
	"hotel/internal/config"
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
func Seed(db *gorm.DB, cfg config.Config) {
	if err := db.AutoMigrate(models.AllForDb()...); err != nil {
		panic(fmt.Sprintf("auto migrate: %s", err))
	}
	seedAmenities(db)
	seedParkingSpotStatuses(db)
	seedParkingSpotTypes(db)
	seedRoomStatuses(db)
	seedRoomTypes(db)
	seedPermissions(db)
	seedPermissionTemplates(db)
	seedSanaReferenceData(db, cfg)
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
	t := Translations["amenity"]

	amenities := []models.Amenity{
		{TranslateBase: models.TranslateBase{Slug: "wifi", Translation: t["wifi"]}},
		{TranslateBase: models.TranslateBase{Slug: "tv", Translation: t["tv"]}},
		{TranslateBase: models.TranslateBase{Slug: "air-conditioning", Translation: t["air-conditioning"]}},
		{TranslateBase: models.TranslateBase{Slug: "mini-Bar", Translation: t["mini-bar"]}},
		{TranslateBase: models.TranslateBase{Slug: "safe", Translation: t["safe"]}},
		{TranslateBase: models.TranslateBase{Slug: "ocean-view", Translation: t["ocean-view"]}},
		{TranslateBase: models.TranslateBase{Slug: "city-view", Translation: t["city-view"]}},
		{TranslateBase: models.TranslateBase{Slug: "balcony", Translation: t["balcony"]}},
		{TranslateBase: models.TranslateBase{Slug: "jacuzzi", Translation: t["jacuzzi"]}},
		{TranslateBase: models.TranslateBase{Slug: "room-service", Translation: t["room-service"]}},
	}

	seed(db, amenities)
}

func seedParkingSpotTypes(db *gorm.DB) {
	t := Translations["parking-spot-type"]

	types := []models.ParkingSpotType{
		{TranslateBase: models.TranslateBase{Slug: "standard", Translation: t["standard"]}},
		{TranslateBase: models.TranslateBase{Slug: "handicap", Translation: t["handicap"]}},
		{TranslateBase: models.TranslateBase{Slug: "electric", Translation: t["electric"]}},
		{TranslateBase: models.TranslateBase{Slug: "compact", Translation: t["compact"]}},
		{TranslateBase: models.TranslateBase{Slug: "large", Translation: t["large"]}},
	}

	seed(db, types)
}

func seedParkingSpotStatuses(db *gorm.DB) {
	t := Translations["parking-spot-status"]

	statuses := []models.ParkingSpotStatus{
		{TranslateBase: models.TranslateBase{Slug: "available", Translation: t["available"]}},
		{TranslateBase: models.TranslateBase{Slug: "occupied", Translation: t["occupied"]}},
		{TranslateBase: models.TranslateBase{Slug: "reserved", Translation: t["reserved"]}},
		{TranslateBase: models.TranslateBase{Slug: "maintenance", Translation: t["maintenance"]}},
	}

	seed(db, statuses)
}

func seedRoomStatuses(db *gorm.DB) {
	t := Translations["room-status"]

	statuses := []models.RoomStatus{
		{TranslateBase: models.TranslateBase{Slug: "available", Translation: t["available"]}, ColorHex: "2ECC71"},     // green
		{TranslateBase: models.TranslateBase{Slug: "occupied", Translation: t["occupied"]}, ColorHex: "E74C3C"},       // red
		{TranslateBase: models.TranslateBase{Slug: "reserved", Translation: t["reserved"]}, ColorHex: "F39C12"},       // orange
		{TranslateBase: models.TranslateBase{Slug: "maintenance", Translation: t["maintenance"]}, ColorHex: "95A5A6"}, // gray
	}

	seed(db, statuses)
}

func seedRoomTypes(db *gorm.DB) {
	t := Translations["room-type"]

	statuses := []models.RoomType{
		{TranslateBase: models.TranslateBase{Slug: "single", Translation: t["single"]}, ColorHex: "87CEEB"},
		{TranslateBase: models.TranslateBase{Slug: "double", Translation: t["double"]}, ColorHex: "98FB98"},
		{TranslateBase: models.TranslateBase{Slug: "suite", Translation: t["suite"]}, ColorHex: "DDA0DD"},
		{TranslateBase: models.TranslateBase{Slug: "deluxe", Translation: t["deluxe"]}, ColorHex: "FFB6C1"},
	}

	seed(db, statuses)
}
