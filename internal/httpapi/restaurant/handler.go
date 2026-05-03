package restaurant

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"time"

	"github.com/go-fuego/fuego"
)

type RestaurantModule struct {
	*h.API
}

func (m RestaurantModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	p := RestaurantModule{api}

	fuego.Get(s, "/stats", p.restaurantStats)

	inventoryGroup := fuego.Group(s, "/inventory")
	fuego.Get(inventoryGroup, "/", h.ListModel(api.Db, models.InventoryItem{}))
	fuego.Post(inventoryGroup, "/", h.CreateModel(api.Db, models.InventoryItem{}))
	fuego.Get(inventoryGroup, "/{id}", h.GetModel(api.Db, models.InventoryItem{}))
	fuego.Put(inventoryGroup, "/{id}", h.UpdateModel(api.Db, models.InventoryItem{}))
	fuego.Delete(inventoryGroup, "/{id}", h.DeleteModel(api.Db, models.InventoryItem{}))

	billsGroup := fuego.Group(s, "/bills")
	fuego.Get(billsGroup, "/", h.ListModel(api.Db, models.RestaurantBill{}))
	fuego.Post(billsGroup, "/", h.CreateModel(api.Db, models.RestaurantBill{}))
	fuego.Get(billsGroup, "/{id}", h.GetModel(api.Db, models.RestaurantBill{}))
	fuego.Put(billsGroup, "/{id}", h.UpdateModel(api.Db, models.RestaurantBill{}))
	fuego.Delete(billsGroup, "/{id}", h.DeleteModel(api.Db, models.RestaurantBill{}))
	fuego.Post(billsGroup, "/{id}/settle", p.settleBill)

	transactionsGroup := fuego.Group(s, "/transactions")
	fuego.Get(transactionsGroup, "/", h.ListModel(api.Db, models.MealTransaction{}))
	fuego.Post(transactionsGroup, "/", h.CreateModel(api.Db, models.MealTransaction{}))
	fuego.Get(transactionsGroup, "/{id}", h.GetModel(api.Db, models.MealTransaction{}))
	fuego.Put(transactionsGroup, "/{id}", h.UpdateModel(api.Db, models.MealTransaction{}))
	fuego.Delete(transactionsGroup, "/{id}", h.DeleteModel(api.Db, models.MealTransaction{}))
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
