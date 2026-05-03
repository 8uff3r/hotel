package seed

import (
	"hotel/internal/models"

	"gorm.io/gorm"
)

func seedRestaurantReferenceData(db *gorm.DB) {
	seedInventoryItemCategories(db)
	seedInventoryItemUnits(db)
	seedInventoryItemStatuses(db)
	seedRestaurantBillStatuses(db)
}

func seedInventoryItemCategories(db *gorm.DB) {
	t := Translations["inventory-item-category"]

	categories := []models.InventoryItemCategory{
		{TranslateBase: models.TranslateBase{Slug: "food", Translation: t["food"]}},
		{TranslateBase: models.TranslateBase{Slug: "beverage", Translation: t["beverage"]}},
		{TranslateBase: models.TranslateBase{Slug: "dessert", Translation: t["dessert"]}},
		{TranslateBase: models.TranslateBase{Slug: "other", Translation: t["other"]}},
	}

	seed(db, categories)
}

func seedInventoryItemUnits(db *gorm.DB) {
	t := Translations["inventory-item-unit"]

	units := []models.InventoryItemUnit{
		{TranslateBase: models.TranslateBase{Slug: "piece", Translation: t["piece"]}},
		{TranslateBase: models.TranslateBase{Slug: "gram", Translation: t["gram"]}},
		{TranslateBase: models.TranslateBase{Slug: "kilogram", Translation: t["kilogram"]}},
		{TranslateBase: models.TranslateBase{Slug: "liter", Translation: t["liter"]}},
		{TranslateBase: models.TranslateBase{Slug: "pack", Translation: t["pack"]}},
	}

	seed(db, units)
}

func seedInventoryItemStatuses(db *gorm.DB) {
	t := Translations["inventory-item-status"]

	statuses := []models.InventoryItemStatus{
		{TranslateBase: models.TranslateBase{Slug: "active", Translation: t["active"]}, ColorHex: "2ECC71"},
		{TranslateBase: models.TranslateBase{Slug: "inactive", Translation: t["inactive"]}, ColorHex: "E74C3C"},
		{TranslateBase: models.TranslateBase{Slug: "discontinued", Translation: t["discontinued"]}, ColorHex: "95A5A6"},
	}

	seed(db, statuses)
}

func seedRestaurantBillStatuses(db *gorm.DB) {
	t := Translations["restaurant-bill-status"]

	statuses := []models.RestaurantBillStatus{
		{TranslateBase: models.TranslateBase{Slug: "open", Translation: t["open"]}, ColorHex: "3498DB"},
		{TranslateBase: models.TranslateBase{Slug: "settled", Translation: t["settled"]}, ColorHex: "2ECC71"},
		{TranslateBase: models.TranslateBase{Slug: "cancelled", Translation: t["cancelled"]}, ColorHex: "E74C3C"},
	}

	seed(db, statuses)
}
