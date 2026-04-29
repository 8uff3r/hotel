package seed

import (
	"hotel/internal/models"

	"gorm.io/gorm"
)

func seedPermissions(db *gorm.DB) {
	t := Translations["Permission"]

	permissions := []models.Permission{
		// Dashboard
		{TranslateBase: models.TranslateBase{Name: "View Dashboard", Translation: t["View Dashboard"]}, Page: "index", Action: models.PermissionActionRead, Category: "Dashboard"},

		// Guests
		{TranslateBase: models.TranslateBase{Name: "View Guests", Translation: t["View Guests"]}, Page: "guests", Action: models.PermissionActionRead, Category: "Guests"},
		{TranslateBase: models.TranslateBase{Name: "Create Guest", Translation: t["Create Guest"]}, Page: "guests", Action: models.PermissionActionCreate, Category: "Guests"},
		{TranslateBase: models.TranslateBase{Name: "Update Guest", Translation: t["Update Guest"]}, Page: "guests", Action: models.PermissionActionUpdate, Category: "Guests"},
		{TranslateBase: models.TranslateBase{Name: "Delete Guest", Translation: t["Delete Guest"]}, Page: "guests", Action: models.PermissionActionDelete, Category: "Guests"},
		{TranslateBase: models.TranslateBase{Name: "Export Guests", Translation: t["Export Guests"]}, Page: "guests", Action: models.PermissionActionExport, Category: "Guests"},
		{TranslateBase: models.TranslateBase{Name: "View Guest Settlement", Translation: t["View Guest Settlement"]}, Page: "guests/settle", Action: models.PermissionActionRead, Category: "Guests"},

		// Rooms
		{TranslateBase: models.TranslateBase{Name: "View Rooms", Translation: t["View Rooms"]}, Page: "rooms", Action: models.PermissionActionRead, Category: "Rooms"},
		{TranslateBase: models.TranslateBase{Name: "Create Room", Translation: t["Create Room"]}, Page: "rooms", Action: models.PermissionActionCreate, Category: "Rooms"},
		{TranslateBase: models.TranslateBase{Name: "Update Room", Translation: t["Update Room"]}, Page: "rooms", Action: models.PermissionActionUpdate, Category: "Rooms"},
		{TranslateBase: models.TranslateBase{Name: "Delete Room", Translation: t["Delete Room"]}, Page: "rooms", Action: models.PermissionActionDelete, Category: "Rooms"},
		{TranslateBase: models.TranslateBase{Name: "View Room Rack", Translation: t["View Room Rack"]}, Page: "rooms/rack", Action: models.PermissionActionRead, Category: "Rooms"},

		// Reservations
		{TranslateBase: models.TranslateBase{Name: "View Reservations", Translation: t["View Reservations"]}, Page: "reservations", Action: models.PermissionActionRead, Category: "Reservations"},
		{TranslateBase: models.TranslateBase{Name: "Create Reservation", Translation: t["Create Reservation"]}, Page: "reservations", Action: models.PermissionActionCreate, Category: "Reservations"},
		{TranslateBase: models.TranslateBase{Name: "Update Reservation", Translation: t["Update Reservation"]}, Page: "reservations", Action: models.PermissionActionUpdate, Category: "Reservations"},
		{TranslateBase: models.TranslateBase{Name: "Delete Reservation", Translation: t["Delete Reservation"]}, Page: "reservations", Action: models.PermissionActionDelete, Category: "Reservations"},

		// Users
		{TranslateBase: models.TranslateBase{Name: "View Users", Translation: t["View Users"]}, Page: "users", Action: models.PermissionActionRead, Category: "Users"},
		{TranslateBase: models.TranslateBase{Name: "Create User", Translation: t["Create User"]}, Page: "users", Action: models.PermissionActionCreate, Category: "Users"},
		{TranslateBase: models.TranslateBase{Name: "Update User", Translation: t["Update User"]}, Page: "users", Action: models.PermissionActionUpdate, Category: "Users"},
		{TranslateBase: models.TranslateBase{Name: "Delete User", Translation: t["Delete User"]}, Page: "users", Action: models.PermissionActionDelete, Category: "Users"},

		// Accounting
		{TranslateBase: models.TranslateBase{Name: "View Accounting", Translation: t["View Accounting"]}, Page: "accounting", Action: models.PermissionActionRead, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "View Accounts", Translation: t["View Accounts"]}, Page: "accounting/accounts", Action: models.PermissionActionRead, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Create Account", Translation: t["Create Account"]}, Page: "accounting/accounts", Action: models.PermissionActionCreate, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Update Account", Translation: t["Update Account"]}, Page: "accounting/accounts", Action: models.PermissionActionUpdate, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Delete Account", Translation: t["Delete Account"]}, Page: "accounting/accounts", Action: models.PermissionActionDelete, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "View Expenses", Translation: t["View Expenses"]}, Page: "accounting/expenses", Action: models.PermissionActionRead, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Create Expense", Translation: t["Create Expense"]}, Page: "accounting/expenses", Action: models.PermissionActionCreate, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Update Expense", Translation: t["Update Expense"]}, Page: "accounting/expenses", Action: models.PermissionActionUpdate, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Delete Expense", Translation: t["Delete Expense"]}, Page: "accounting/expenses", Action: models.PermissionActionDelete, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Export Expenses", Translation: t["Export Expenses"]}, Page: "accounting/expenses", Action: models.PermissionActionExport, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "View Income", Translation: t["View Income"]}, Page: "accounting/income", Action: models.PermissionActionRead, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Create Income", Translation: t["Create Income"]}, Page: "accounting/income", Action: models.PermissionActionCreate, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Update Income", Translation: t["Update Income"]}, Page: "accounting/income", Action: models.PermissionActionUpdate, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Delete Income", Translation: t["Delete Income"]}, Page: "accounting/income", Action: models.PermissionActionDelete, Category: "Accounting"},
		{TranslateBase: models.TranslateBase{Name: "Export Income", Translation: t["Export Income"]}, Page: "accounting/income", Action: models.PermissionActionExport, Category: "Accounting"},

		// Parking
		{TranslateBase: models.TranslateBase{Name: "View Parking", Translation: t["View Parking"]}, Page: "parking", Action: models.PermissionActionRead, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "View Parking Lots", Translation: t["View Parking Lots"]}, Page: "parking/lots", Action: models.PermissionActionRead, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Create Parking Lot", Translation: t["Create Parking Lot"]}, Page: "parking/lots", Action: models.PermissionActionCreate, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Update Parking Lot", Translation: t["Update Parking Lot"]}, Page: "parking/lots", Action: models.PermissionActionUpdate, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Delete Parking Lot", Translation: t["Delete Parking Lot"]}, Page: "parking/lots", Action: models.PermissionActionDelete, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "View Parking Spots", Translation: t["View Parking Spots"]}, Page: "parking/spots", Action: models.PermissionActionRead, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Create Parking Spot", Translation: t["Create Parking Spot"]}, Page: "parking/spots", Action: models.PermissionActionCreate, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Update Parking Spot", Translation: t["Update Parking Spot"]}, Page: "parking/spots", Action: models.PermissionActionUpdate, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Delete Parking Spot", Translation: t["Delete Parking Spot"]}, Page: "parking/spots", Action: models.PermissionActionDelete, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "View Vehicles", Translation: t["View Vehicles"]}, Page: "parking/vehicles", Action: models.PermissionActionRead, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Create Vehicle", Translation: t["Create Vehicle"]}, Page: "parking/vehicles", Action: models.PermissionActionCreate, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Update Vehicle", Translation: t["Update Vehicle"]}, Page: "parking/vehicles", Action: models.PermissionActionUpdate, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Delete Vehicle", Translation: t["Delete Vehicle"]}, Page: "parking/vehicles", Action: models.PermissionActionDelete, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "View Parking Transactions", Translation: t["View Parking Transactions"]}, Page: "parking/transactions", Action: models.PermissionActionRead, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Create Parking Transaction", Translation: t["Create Parking Transaction"]}, Page: "parking/transactions", Action: models.PermissionActionCreate, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Update Parking Transaction", Translation: t["Update Parking Transaction"]}, Page: "parking/transactions", Action: models.PermissionActionUpdate, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Export Parking Transactions", Translation: t["Export Parking Transactions"]}, Page: "parking/transactions", Action: models.PermissionActionExport, Category: "Parking"},
		{TranslateBase: models.TranslateBase{Name: "Parking Check-In", Translation: t["Parking Check-In"]}, Page: "parking/transactions/check-in", Action: models.PermissionActionRead, Category: "Parking"},
	}

	seed(db, permissions)
}

func seedPermissionTemplates(db *gorm.DB) {
	t := Translations["PermissionTemplate"]

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
			TranslateBase: models.TranslateBase{Name: "admin", Translation: t["admin"]},
			Description:   "Full administrative access to all features",
		}
		db.Create(&adminTemplate)
	}
	db.Model(&adminTemplate).Association("Permissions").Append(&allPerms)

	var managerTemplate models.PermissionTemplate
	if err := db.Where("name = ?", "manager").First(&managerTemplate).Error; err == gorm.ErrRecordNotFound {
		managerTemplate = models.PermissionTemplate{
			TranslateBase: models.TranslateBase{Name: "manager", Translation: t["manager"]},
			Description:   "Management access to all operational features",
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
			TranslateBase: models.TranslateBase{Name: "receptionist", Translation: t["receptionist"]},
			Description:   "Front desk operations: check-in, reservations, guest management, parking",
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
			TranslateBase: models.TranslateBase{Name: "staff", Translation: t["staff"]},
			Description:   "Basic staff access - read only",
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
			TranslateBase: models.TranslateBase{Name: "housekeeper", Translation: t["housekeeper"]},
			Description:   "Housekeeping access to rooms",
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
