package seed

import (
	"encoding/json"
	"fmt"
	"log"

	"hotel/internal/config"
	"hotel/internal/models"

	_ "embed"

	"gorm.io/gorm"
)

//go:embed translations.json
var translationsFile []byte

//go:embed countries.json
var countriesFile []byte

var Translations map[string]map[string]models.Translation

func init() {
	err := json.Unmarshal(translationsFile, &Translations)
	if err != nil {
		log.Fatalf("failed to load translations: %v", err)
	}
}

func Seed(db *gorm.DB, cfg config.Config) {
	if err := db.AutoMigrate(models.AllForDB()...); err != nil {
		panic(fmt.Sprintf("auto migrate: %s", err))
	}

	seedSanaReferenceData(db, cfg)

	seedAmenities(db)
	seedParkingSpotStatuses(db)
	seedParkingSpotTypes(db)
	seedRoomStatuses(db)
	seedRoomTypes(db)
	seedParkingLotStatuses(db)
	seedExpenseCategories(db)
	seedIncomeCategories(db)
	seedPaymentStatuses(db)
	seedPaymentMethods(db)
	seedVehicleTypes(db)
	seedPermissions(db)
	seedPermissionTemplates(db)
	seedCountries(db)
	seedRestaurantReferenceData(db)
	seedAllReferenceData(db)
	seedReservationStatuses(db)
	seedStayStatuses(db)
	seedInvoiceItemTypes(db)
	seedServices(db)
}

type Seedable interface {
	UniqueCondition() any
}

func seed[T Seedable](db *gorm.DB, values []T) error {
	if len(values) == 0 {
		return nil
	}

	var model T
	var count int64
	if err := db.Model(&model).Count(&count).Error; err != nil {
		return fmt.Errorf("error counting %T: %w", model, err)
	}

	// Table is empty — bulk insert is much faster
	if count == 0 {
		return db.CreateInBatches(values, 100).Error
	}

	// Table already has data — use FirstOrCreate to fill any gaps
	for i := range values {
		result := db.Where(values[i].UniqueCondition()).FirstOrCreate(&values[i])
		if result.Error != nil {
			return fmt.Errorf("error seeding %T: %w", values[i], result.Error)
		}
	}
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
		{
			TranslateBase: models.TranslateBase{
				Slug:        string(models.RoomStatusAvailable),
				Translation: t[string(models.RoomStatusAvailable)],
			},
			ColorHex: "2ECC71", // green
		},
		{
			TranslateBase: models.TranslateBase{
				Slug:        string(models.RoomStatusOccupied),
				Translation: t[string(models.RoomStatusOccupied)],
			},
			ColorHex: "E74C3C", // red
		},
		{
			TranslateBase: models.TranslateBase{
				Slug:        string(models.RoomStatusReserved),
				Translation: t[string(models.RoomStatusReserved)],
			},
			ColorHex: "F39C12", // orange
		},
		{
			TranslateBase: models.TranslateBase{
				Slug:        string(models.RoomStatusUnderRepair),
				Translation: t[string(models.RoomStatusUnderRepair)],
			},
			ColorHex: "95A5A6", // gray
		},
		{
			TranslateBase: models.TranslateBase{
				Slug:        string(models.RoomStatusCleaning),
				Translation: t[string(models.RoomStatusCleaning)],
			},
			ColorHex: "4a412a", // dark brown
		},
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

type jsonCountry struct {
	Slug        string             `json:"slug"`
	Translation models.Translation `json:"translation"`
}

func seedCountries(db *gorm.DB) {
	var jsonCountries []jsonCountry
	if err := json.Unmarshal(countriesFile, &jsonCountries); err != nil {
		fmt.Printf("Failed to parse countries: %v\n", err)
		return
	}

	var countries []models.Country
	for _, jc := range jsonCountries {
		countries = append(countries, models.Country{
			TranslateBase: models.TranslateBase{
				Slug:        jc.Slug,
				Translation: jc.Translation,
			},
		})
	}

	seed(db, countries)
}

func seedReservationStatuses(db *gorm.DB) {
	t := Translations["reservation-status"]

	statuses := []models.ReservationStatus{
		{TranslateBase: models.TranslateBase{Slug: "awaiting_payment", Translation: t["awaiting_payment"]}, ColorHex: "FFA500"}, // Amber
		{TranslateBase: models.TranslateBase{Slug: "verified", Translation: t["verified"]}, ColorHex: "2E8B57"},                 // Forest Green
		{TranslateBase: models.TranslateBase{Slug: "accepted", Translation: t["accepted"]}, ColorHex: "1E90FF"},                 // Dodger Blue
		{TranslateBase: models.TranslateBase{Slug: "cancelled", Translation: t["cancelled"]}, ColorHex: "DC143C"},               // Crimson
		{TranslateBase: models.TranslateBase{Slug: "absence", Translation: t["absence"]}, ColorHex: "8B0000"},                   // Dark Red
		{TranslateBase: models.TranslateBase{Slug: "expired", Translation: t["expired"]}, ColorHex: "708090"},                   // Slate Gray
	}

	seed(db, statuses)
}

func seedStayStatuses(db *gorm.DB) {
	t := Translations["stay-status"]

	statuses := []models.StayStatus{
		{TranslateBase: models.TranslateBase{Slug: "waiting", Translation: t["waiting"]}, ColorHex: "FFA500"},
		{TranslateBase: models.TranslateBase{Slug: "resident", Translation: t["resident"]}, ColorHex: "2ECC71"},
		{TranslateBase: models.TranslateBase{Slug: "checked_out", Translation: t["checked_out"]}, ColorHex: "708090"},
		{TranslateBase: models.TranslateBase{Slug: "cancelled", Translation: t["cancelled"]}, ColorHex: "DC143C"},
		{TranslateBase: models.TranslateBase{Slug: "no_show", Translation: t["no_show"]}, ColorHex: "8B0000"},
	}

	seed(db, statuses)
}

func seedInvoiceItemTypes(db *gorm.DB) {
	t := Translations["invoice-item-type"]

	types := []models.TranslateBase{
		{Slug: "room_charge", Translation: t["room_charge"]},
		{Slug: "breakfast", Translation: t["breakfast"]},
		{Slug: "half_board", Translation: t["half_board"]},
		{Slug: "full_board", Translation: t["full_board"]},
		{Slug: "parking", Translation: t["parking"]},
		{Slug: "room_service", Translation: t["room_service"]},
		{Slug: "other", Translation: t["other"]},
	}

	// These are not stored in DB as a separate table; they are string enums
	// But we seed them as reference data for consistency if needed in future
	_ = types
}

func seedServices(db *gorm.DB) {
	// No-op for now; services are hotel-specific and created via UI
}
