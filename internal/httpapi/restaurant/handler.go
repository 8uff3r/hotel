package restaurant

import (
	"time"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type RestaurantModule struct {
	*h.API
}

func (m RestaurantModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	p := RestaurantModule{api}

	fuego.Get(s, "/stats", p.restaurantStats)

	inventoryGroup := fuego.Group(s, "/inventory")
	fuego.Get(inventoryGroup, "/", h.ListModel[models.InventoryItem](api.Db))
	fuego.Post(inventoryGroup, "/", h.CreateModel[models.InventoryItem](api.Db))
	fuego.Get(inventoryGroup, "/{id}", h.GetModel[models.InventoryItem](api.Db))
	fuego.Put(inventoryGroup, "/{id}", h.UpdateModel[models.InventoryItem](api.Db))
	fuego.Delete(inventoryGroup, "/{id}", h.DeleteModel[models.InventoryItem](api.Db))

	fuego.Get(
		inventoryGroup,
		"/categories",
		h.ListModel[models.InventoryItemCategory](api.Db, h.WithTranslation[models.InventoryItemCategory]()),
	)

	fuego.Get(
		inventoryGroup,
		"/units",
		h.ListModel[models.InventoryItemUnit](api.Db, h.WithTranslation[models.InventoryItemUnit]()),
	)

	fuego.Get(
		inventoryGroup,
		"/statuses",
		h.ListModel[models.InventoryItemStatus](api.Db, h.WithTranslation[models.InventoryItemStatus]()),
	)

	billsGroup := fuego.Group(s, "/bills")
	fuego.Get(billsGroup, "/", h.ListModel[models.RestaurantBill](api.Db))
	fuego.Get(billsGroup, "/statuses", h.ListModel[models.RestaurantBillStatus](api.Db))
	fuego.Post(billsGroup, "/", h.CreateModel[models.RestaurantBill](api.Db))
	fuego.Get(billsGroup, "/{id}", h.GetModel[models.RestaurantBill](api.Db))
	fuego.Put(billsGroup, "/{id}", h.UpdateModel[models.RestaurantBill](api.Db))
	fuego.Delete(billsGroup, "/{id}", h.DeleteModel[models.RestaurantBill](api.Db))
	fuego.Post(billsGroup, "/{id}/settle", p.settleBill)

	transactionsGroup := fuego.Group(s, "/transactions")
	fuego.Get(transactionsGroup, "/", h.ListModel[models.MealTransaction](api.Db))
	fuego.Post(transactionsGroup, "/", h.CreateModel[models.MealTransaction](api.Db))
	fuego.Get(transactionsGroup, "/{id}", h.GetModel[models.MealTransaction](api.Db))
	fuego.Put(transactionsGroup, "/{id}", h.UpdateModel[models.MealTransaction](api.Db))
	fuego.Delete(transactionsGroup, "/{id}", h.DeleteModel[models.MealTransaction](api.Db))
}

func (a *RestaurantModule) restaurantStats(c fuego.ContextNoBody) (models.RestaurantStats, error) {
	db := a.Db.WithContext(c)
	var zero models.RestaurantStats

	var totalBills int64
	if err := db.Model(&models.RestaurantBill{}).Count(&totalBills).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "failed"}
	}

	if err := db.Model(&models.RestaurantBill{}).Where("settled = ?", true).Count(&totalBills).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "failed"}
	}

	var totalRevenue float64
	if err := db.Model(&models.RestaurantBill{}).Where("settled = ?", true).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "failed"}
	}

	var internalRevenue float64
	if err := db.Model(&models.RestaurantBill{}).Where("settled = ? AND is_external = ?", true, false).Select("COALESCE(SUM(total_amount), 0)").Scan(&internalRevenue).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "failed"}
	}

	var externalRevenue float64
	if err := db.Model(&models.RestaurantBill{}).Where("settled = ? AND is_external = ?", true, true).Select("COALESCE(SUM(total_amount), 0)").Scan(&externalRevenue).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "failed"}
	}

	var totalMeals int64
	if err := db.Model(&models.MealTransaction{}).Count(&totalMeals).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "failed"}
	}

	return models.RestaurantStats{
		TotalBills:      totalBills,
		TotalRevenue:    totalRevenue,
		InternalRevenue: internalRevenue,
		ExternalRevenue: externalRevenue,
		TotalMeals:      totalMeals,
	}, nil
}

func (a *RestaurantModule) settleBill(c fuego.ContextNoBody) (map[string]bool, error) {
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_id"}
	}
	now := time.Now().UTC()
	res := a.Db.WithContext(c).Model(&models.RestaurantBill{}).Where("id = ?", id).Updates(map[string]any{
		"settled":    true,
		"settled_at": now,
		"status":     "settled",
	})
	if res.Error != nil {
		return nil, fuego.InternalServerError{Title: "update_failed"}
	}
	if res.RowsAffected == 0 {
		return nil, fuego.NotFoundError{}
	}
	return map[string]bool{"ok": true}, nil
}
