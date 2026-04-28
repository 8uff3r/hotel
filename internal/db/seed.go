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
	seedPermissions(db)
	seedPermissionTemplates(db)
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
		{TranslateBase: models.TranslateBase{Name: "Available", Translation: t["Available"]}, ColorHex: "2ECC71"},     // green
		{TranslateBase: models.TranslateBase{Name: "Occupied", Translation: t["Occupied"]}, ColorHex: "E74C3C"},       // red
		{TranslateBase: models.TranslateBase{Name: "Reserved", Translation: t["Reserved"]}, ColorHex: "F39C12"},       // orange
		{TranslateBase: models.TranslateBase{Name: "Maintenance", Translation: t["Maintenance"]}, ColorHex: "95A5A6"}, // gray
	}

	seed(db, statuses)
}

func seedRoomTypes(db *gorm.DB) {
	t := Translations["RoomType"]

	statuses := []models.RoomType{
		{TranslateBase: models.TranslateBase{Name: "Single", Translation: t["Single"]}, ColorHex: "87CEEB"},
		{TranslateBase: models.TranslateBase{Name: "Double", Translation: t["Double"]}, ColorHex: "98FB98"},
		{TranslateBase: models.TranslateBase{Name: "Suite", Translation: t["Suite"]}, ColorHex: "DDA0DD"},
		{TranslateBase: models.TranslateBase{Name: "Deluxe", Translation: t["Deluxe"]}, ColorHex: "FFB6C1"},
	}

	seed(db, statuses)
}

func seedPermissions(db *gorm.DB) {
	permissions := []models.Permission{
		// Dashboard
		{Page: "index", Action: models.PermissionActionRead, Label: "View Dashboard", Category: "Dashboard"},

		// Guests
		{Page: "guests", Action: models.PermissionActionRead, Label: "View Guests", Category: "Guests"},
		{Page: "guests", Action: models.PermissionActionCreate, Label: "Create Guest", Category: "Guests"},
		{Page: "guests", Action: models.PermissionActionUpdate, Label: "Update Guest", Category: "Guests"},
		{Page: "guests", Action: models.PermissionActionDelete, Label: "Delete Guest", Category: "Guests"},
		{Page: "guests", Action: models.PermissionActionExport, Label: "Export Guests", Category: "Guests"},
		{Page: "guests/settle", Action: models.PermissionActionRead, Label: "View Guest Settlement", Category: "Guests"},

		// Rooms
		{Page: "rooms", Action: models.PermissionActionRead, Label: "View Rooms", Category: "Rooms"},
		{Page: "rooms", Action: models.PermissionActionCreate, Label: "Create Room", Category: "Rooms"},
		{Page: "rooms", Action: models.PermissionActionUpdate, Label: "Update Room", Category: "Rooms"},
		{Page: "rooms", Action: models.PermissionActionDelete, Label: "Delete Room", Category: "Rooms"},
		{Page: "rooms/rack", Action: models.PermissionActionRead, Label: "View Room Rack", Category: "Rooms"},

		// Reservations
		{Page: "reservations", Action: models.PermissionActionRead, Label: "View Reservations", Category: "Reservations"},
		{Page: "reservations", Action: models.PermissionActionCreate, Label: "Create Reservation", Category: "Reservations"},
		{Page: "reservations", Action: models.PermissionActionUpdate, Label: "Update Reservation", Category: "Reservations"},
		{Page: "reservations", Action: models.PermissionActionDelete, Label: "Delete Reservation", Category: "Reservations"},

		// Users
		{Page: "users", Action: models.PermissionActionRead, Label: "View Users", Category: "Users"},
		{Page: "users", Action: models.PermissionActionCreate, Label: "Create User", Category: "Users"},
		{Page: "users", Action: models.PermissionActionUpdate, Label: "Update User", Category: "Users"},
		{Page: "users", Action: models.PermissionActionDelete, Label: "Delete User", Category: "Users"},

		// Accounting
		{Page: "accounting", Action: models.PermissionActionRead, Label: "View Accounting", Category: "Accounting"},
		{Page: "accounting/accounts", Action: models.PermissionActionRead, Label: "View Accounts", Category: "Accounting"},
		{Page: "accounting/accounts", Action: models.PermissionActionCreate, Label: "Create Account", Category: "Accounting"},
		{Page: "accounting/accounts", Action: models.PermissionActionUpdate, Label: "Update Account", Category: "Accounting"},
		{Page: "accounting/accounts", Action: models.PermissionActionDelete, Label: "Delete Account", Category: "Accounting"},
		{Page: "accounting/expenses", Action: models.PermissionActionRead, Label: "View Expenses", Category: "Accounting"},
		{Page: "accounting/expenses", Action: models.PermissionActionCreate, Label: "Create Expense", Category: "Accounting"},
		{Page: "accounting/expenses", Action: models.PermissionActionUpdate, Label: "Update Expense", Category: "Accounting"},
		{Page: "accounting/expenses", Action: models.PermissionActionDelete, Label: "Delete Expense", Category: "Accounting"},
		{Page: "accounting/expenses", Action: models.PermissionActionExport, Label: "Export Expenses", Category: "Accounting"},
		{Page: "accounting/income", Action: models.PermissionActionRead, Label: "View Income", Category: "Accounting"},
		{Page: "accounting/income", Action: models.PermissionActionCreate, Label: "Create Income", Category: "Accounting"},
		{Page: "accounting/income", Action: models.PermissionActionUpdate, Label: "Update Income", Category: "Accounting"},
		{Page: "accounting/income", Action: models.PermissionActionDelete, Label: "Delete Income", Category: "Accounting"},
		{Page: "accounting/income", Action: models.PermissionActionExport, Label: "Export Income", Category: "Accounting"},

		// Parking
		{Page: "parking", Action: models.PermissionActionRead, Label: "View Parking", Category: "Parking"},
		{Page: "parking/lots", Action: models.PermissionActionRead, Label: "View Parking Lots", Category: "Parking"},
		{Page: "parking/lots", Action: models.PermissionActionCreate, Label: "Create Parking Lot", Category: "Parking"},
		{Page: "parking/lots", Action: models.PermissionActionUpdate, Label: "Update Parking Lot", Category: "Parking"},
		{Page: "parking/lots", Action: models.PermissionActionDelete, Label: "Delete Parking Lot", Category: "Parking"},
		{Page: "parking/spots", Action: models.PermissionActionRead, Label: "View Parking Spots", Category: "Parking"},
		{Page: "parking/spots", Action: models.PermissionActionCreate, Label: "Create Parking Spot", Category: "Parking"},
		{Page: "parking/spots", Action: models.PermissionActionUpdate, Label: "Update Parking Spot", Category: "Parking"},
		{Page: "parking/spots", Action: models.PermissionActionDelete, Label: "Delete Parking Spot", Category: "Parking"},
		{Page: "parking/vehicles", Action: models.PermissionActionRead, Label: "View Vehicles", Category: "Parking"},
		{Page: "parking/vehicles", Action: models.PermissionActionCreate, Label: "Create Vehicle", Category: "Parking"},
		{Page: "parking/vehicles", Action: models.PermissionActionUpdate, Label: "Update Vehicle", Category: "Parking"},
		{Page: "parking/vehicles", Action: models.PermissionActionDelete, Label: "Delete Vehicle", Category: "Parking"},
		{Page: "parking/transactions", Action: models.PermissionActionRead, Label: "View Parking Transactions", Category: "Parking"},
		{Page: "parking/transactions", Action: models.PermissionActionCreate, Label: "Create Parking Transaction", Category: "Parking"},
		{Page: "parking/transactions", Action: models.PermissionActionUpdate, Label: "Update Parking Transaction", Category: "Parking"},
		{Page: "parking/transactions", Action: models.PermissionActionExport, Label: "Export Parking Transactions", Category: "Parking"},
		{Page: "parking/transactions/check-in", Action: models.PermissionActionRead, Label: "Parking Check-In", Category: "Parking"},
	}

	seed(db, permissions)
}

func seedPermissionTemplates(db *gorm.DB) {
	var allPerms []models.Permission
	db.Find(&allPerms)

	permMap := make(map[string]uint)
	for _, p := range allPerms {
		key := p.Page + ":" + string(p.Action)
		permMap[key] = p.ID
	}

	var adminTemplate models.PermissionTemplate
	if err := db.Where("name = ?", "admin").First(&adminTemplate).Error; err == gorm.ErrRecordNotFound {
		adminTemplate = models.PermissionTemplate{
			Name:        "admin",
			Description: "Full administrative access to all features",
		}
		db.Create(&adminTemplate)
	}
	db.Model(&adminTemplate).Association("Permissions").Append(&allPerms)

	var managerTemplate models.PermissionTemplate
	if err := db.Where("name = ?", "manager").First(&managerTemplate).Error; err == gorm.ErrRecordNotFound {
		managerTemplate = models.PermissionTemplate{
			Name:        "manager",
			Description: "Management access to all operational features",
		}
		db.Create(&managerTemplate)
	}
	var managerPerms []models.Permission
	for _, p := range allPerms {
		if p.Page != "users" {
			managerPerms = append(managerPerms, p)
		}
	}
	db.Model(&managerTemplate).Association("Permissions").Replace(&managerPerms)

	var receptionistTemplate models.PermissionTemplate
	if err := db.Where("name = ?", "receptionist").First(&receptionistTemplate).Error; err == gorm.ErrRecordNotFound {
		receptionistTemplate = models.PermissionTemplate{
			Name:        "receptionist",
			Description: "Front desk operations: check-in, reservations, guest management, parking",
		}
		db.Create(&receptionistTemplate)
	}
	var receptionistPerms []models.Permission
	for _, p := range allPerms {
		if p.Page == "index" ||
			p.Page == "guests" || p.Page == "guests/settle" ||
			p.Page == "rooms" || p.Page == "rooms/rack" ||
			p.Page == "reservations" ||
			p.Page == "parking" || p.Page == "parking/transactions" || p.Page == "parking/transactions/check-in" ||
			p.Page == "parking/vehicles" {
			receptionistPerms = append(receptionistPerms, p)
		}
	}
	db.Model(&receptionistTemplate).Association("Permissions").Replace(&receptionistPerms)

	var staffTemplate models.PermissionTemplate
	if err := db.Where("name = ?", "staff").First(&staffTemplate).Error; err == gorm.ErrRecordNotFound {
		staffTemplate = models.PermissionTemplate{
			Name:        "staff",
			Description: "Basic staff access - read only",
		}
		db.Create(&staffTemplate)
	}
	var staffPerms []models.Permission
	for _, p := range allPerms {
		if p.Action == models.PermissionActionRead {
			staffPerms = append(staffPerms, p)
		}
	}
	db.Model(&staffTemplate).Association("Permissions").Replace(&staffPerms)

	var housekeeperTemplate models.PermissionTemplate
	if err := db.Where("name = ?", "housekeeper").First(&housekeeperTemplate).Error; err == gorm.ErrRecordNotFound {
		housekeeperTemplate = models.PermissionTemplate{
			Name:        "housekeeper",
			Description: "Housekeeping access to rooms",
		}
		db.Create(&housekeeperTemplate)
	}
	var housekeeperPerms []models.Permission
	for _, p := range allPerms {
		if p.Page == "rooms" || p.Page == "rooms/rack" {
			housekeeperPerms = append(housekeeperPerms, p)
		}
	}
	db.Model(&housekeeperTemplate).Association("Permissions").Replace(&housekeeperPerms)
}
