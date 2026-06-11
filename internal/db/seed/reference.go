package seed

import (
	"hotel/internal/models"

	"gorm.io/gorm"
)

func seedAllReferenceData(db *gorm.DB) {
	// seedGuestCompanionRelations(db)
	seedParkingLotStatuses(db)
	seedExpenseCategories(db)
	seedIncomeCategories(db)
	seedPaymentStatuses(db)
	seedPaymentMethods(db)
	seedVehicleTypes(db)
}

// TODO: this must be removed and the seed function in ./sana.go must be used.
//
// func seedGuestCompanionRelations(db *gorm.DB) {
// 	t := Translations["guest-companion-relations"]
// 	relations := []models.FamilyRelationship{
// 		{TranslateBase: models.TranslateBase{Slug: "spouse", Translation: t["spouse"]}, SanaID: uuid.NewString()},
// 		{TranslateBase: models.TranslateBase{Slug: "child", Translation: t["child"]}, SanaID: uuid.NewString()},
// 		{TranslateBase: models.TranslateBase{Slug: "parent", Translation: t["parent"]}, SanaID: uuid.NewString()},
// 		{TranslateBase: models.TranslateBase{Slug: "sibling", Translation: t["sibling"]}, SanaID: uuid.NewString()},
// 		{TranslateBase: models.TranslateBase{Slug: "relative", Translation: t["relative"]}, SanaID: uuid.NewString()},
// 		{TranslateBase: models.TranslateBase{Slug: "friend", Translation: t["friend"]}, SanaID: uuid.NewString()},
// 		{TranslateBase: models.TranslateBase{Slug: "colleague", Translation: t["colleague"]}, SanaID: uuid.NewString()},
// 		{TranslateBase: models.TranslateBase{Slug: "other", Translation: t["other"]}, SanaID: uuid.NewString()},
// 	}
//
// 	_ = seed(db, relations)
// }

func seedParkingLotStatuses(db *gorm.DB) {
	t := Translations["parking-lot-status"]

	statuses := []models.ParkingLotStatus{
		{TranslateBase: models.TranslateBase{Slug: "active", Translation: t["active"]}},
		{TranslateBase: models.TranslateBase{Slug: "full", Translation: t["full"]}},
		{TranslateBase: models.TranslateBase{Slug: "closed", Translation: t["closed"]}},
	}

	_ = seed(db, statuses)
}

func seedExpenseCategories(db *gorm.DB) {
	t := Translations["expense-category"]

	categories := []models.ExpenseCategory{
		{TranslateBase: models.TranslateBase{Slug: "rooms", Translation: t["rooms"]}},
		{TranslateBase: models.TranslateBase{Slug: "food", Translation: t["food"]}},
		{TranslateBase: models.TranslateBase{Slug: "maintenance", Translation: t["maintenance"]}},
		{TranslateBase: models.TranslateBase{Slug: "utilities", Translation: t["utilities"]}},
		{TranslateBase: models.TranslateBase{Slug: "salaries", Translation: t["salaries"]}},
		{TranslateBase: models.TranslateBase{Slug: "marketing", Translation: t["marketing"]}},
		{TranslateBase: models.TranslateBase{Slug: "other", Translation: t["other"]}},
	}

	_ = seed(db, categories)
}

func seedIncomeCategories(db *gorm.DB) {
	t := Translations["income-category"]

	categories := []models.IncomeCategory{
		{TranslateBase: models.TranslateBase{Slug: "room_revenue", Translation: t["room_revenue"]}},
		{TranslateBase: models.TranslateBase{Slug: "parking", Translation: t["parking"]}},
		{TranslateBase: models.TranslateBase{Slug: "restaurant", Translation: t["restaurant"]}},
		{TranslateBase: models.TranslateBase{Slug: "other", Translation: t["other"]}},
	}

	_ = seed(db, categories)
}

func seedPaymentStatuses(db *gorm.DB) {
	t := Translations["payment-status"]

	statuses := []models.PaymentStatus{
		{TranslateBase: models.TranslateBase{Slug: "pending", Translation: t["pending"]}, ColorHex: "F39C12"},
		{TranslateBase: models.TranslateBase{Slug: "received", Translation: t["received"]}, ColorHex: "2ECC71"},
		{TranslateBase: models.TranslateBase{Slug: "cancelled", Translation: t["cancelled"]}, ColorHex: "E74C3C"},
		{TranslateBase: models.TranslateBase{Slug: "refunded", Translation: t["refunded"]}, ColorHex: "9B59B6"},
	}

	_ = seed(db, statuses)
}

func seedPaymentMethods(db *gorm.DB) {
	t := Translations["payment-method"]

	methods := []models.PaymentMethod{
		{TranslateBase: models.TranslateBase{Slug: "cash", Translation: t["cash"]}},
		{TranslateBase: models.TranslateBase{Slug: "card", Translation: t["card"]}},
		{TranslateBase: models.TranslateBase{Slug: "bank_transfer", Translation: t["bank_transfer"]}},
		{TranslateBase: models.TranslateBase{Slug: "cheque", Translation: t["cheque"]}},
		{TranslateBase: models.TranslateBase{Slug: "credit", Translation: t["credit"]}},
		{TranslateBase: models.TranslateBase{Slug: "online", Translation: t["online"]}},
		{TranslateBase: models.TranslateBase{Slug: "contracting_party", Translation: t["contracting_party"]}},
		{TranslateBase: models.TranslateBase{Slug: "agency", Translation: t["agency"]}},
		{TranslateBase: models.TranslateBase{Slug: "other", Translation: t["other"]}},
	}

	_ = seed(db, methods)
}

func seedVehicleTypes(db *gorm.DB) {
	t := Translations["vehicle-type"]

	types := []models.VehicleType{
		{TranslateBase: models.TranslateBase{Slug: "car", Translation: t["car"]}},
		{TranslateBase: models.TranslateBase{Slug: "motorcycle", Translation: t["motorcycle"]}},
		{TranslateBase: models.TranslateBase{Slug: "truck", Translation: t["truck"]}},
		{TranslateBase: models.TranslateBase{Slug: "van", Translation: t["van"]}},
		{TranslateBase: models.TranslateBase{Slug: "other", Translation: t["other"]}},
	}

	_ = seed(db, types)
}
