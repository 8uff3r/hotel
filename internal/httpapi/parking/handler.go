package parking

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"time"

	"github.com/go-fuego/fuego"
)

type ParkingModule struct {
	*h.API
}

func (m ParkingModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	p := ParkingModule{api}

	fuego.Get(s, "/stats", p.parkingStats)

	lotsGroup := fuego.Group(s, "/lots")
	fuego.Get(lotsGroup, "/", h.ListModel(api.Db, models.ParkingLot{}, nil))
	fuego.Post(lotsGroup, "/", h.CreateModel(api.Db, models.ParkingLot{}))
	fuego.Get(lotsGroup, "/{id}", h.GetModel(api.Db, models.ParkingLot{}, nil))
	fuego.Put(lotsGroup, "/{id}", h.UpdateModel(api.Db, models.ParkingLot{}))
	fuego.Delete(lotsGroup, "/{id}", h.DeleteModel(api.Db, models.ParkingLot{}))

	spotsGroup := fuego.Group(s, "/spots")
	fuego.Get(spotsGroup, "/", h.ListModel(api.Db, models.ParkingSpot{}, nil))
	fuego.Post(spotsGroup, "/", h.CreateModel(api.Db, models.ParkingSpot{}))
	fuego.Get(spotsGroup, "/{id}", h.GetModel(api.Db, models.ParkingSpot{}, nil))
	fuego.Put(spotsGroup, "/{id}", h.UpdateModel(api.Db, models.ParkingSpot{}))
	fuego.Delete(spotsGroup, "/{id}", h.DeleteModel(api.Db, models.ParkingSpot{}))
	fuego.Get(spotsGroup, "/statuses", h.ListModel(api.Db, models.ParkingSpotStatus{}, nil))
	fuego.Get(spotsGroup, "/types", h.ListModel(api.Db, models.ParkingSpotType{}, nil))

	vehiclesGroup := fuego.Group(s, "/vehicles")
	fuego.Get(vehiclesGroup, "/", h.ListModel(api.Db, models.Vehicle{}, nil))
	fuego.Post(vehiclesGroup, "/", h.CreateModel(api.Db, models.Vehicle{}))
	fuego.Get(vehiclesGroup, "/{id}", h.GetModel(api.Db, models.Vehicle{}, nil))
	fuego.Put(vehiclesGroup, "/{id}", h.UpdateModel(api.Db, models.Vehicle{}))
	fuego.Delete(vehiclesGroup, "/{id}", h.DeleteModel(api.Db, models.Vehicle{}))

	transactionsGroup := fuego.Group(s, "/transactions")
	fuego.Get(transactionsGroup, "/", h.ListModel(api.Db, models.ParkingTransaction{}, nil))
	fuego.Post(transactionsGroup, "/", h.CreateModel(api.Db, models.ParkingTransaction{}))
	fuego.Get(transactionsGroup, "/{id}", h.GetModel(api.Db, models.ParkingTransaction{}, nil))
	fuego.Post(transactionsGroup, "/{id}/check-out", p.transactionsCheckOut)
}

func (a *ParkingModule) transactionsCheckOut(c fuego.ContextNoBody) (map[string]bool, error) {
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_id"}
	}
	now := time.Now().UTC()
	res := a.Db.WithContext(c).Model(&models.ParkingTransaction{}).Where("id = ?", id).Updates(map[string]any{"status": "completed", "exit_time": now})
	if res.Error != nil {
		return nil, fuego.InternalServerError{Title: "update_failed"}
	}
	if res.RowsAffected == 0 {
		return nil, fuego.NotFoundError{}
	}
	return map[string]bool{"ok": true}, nil
}

func (a *ParkingModule) parkingStats(c fuego.ContextNoBody) (models.ParkingStats, error) {
	db := a.Db.WithContext(c)
	var totalLots int64
	var zero models.ParkingStats
	if err := db.Model(&models.ParkingLot{}).Count(&totalLots).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "failed"}
	}

	var totalSpots int64
	if err := db.Model(&models.ParkingSpot{}).Count(&totalSpots).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "failed"}
	}

	var availableSpots int64
	if err := db.Model(&models.ParkingSpot{}).Where("status = ?", "available").Count(&availableSpots).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "failed"}
	}

	return models.ParkingStats{
		Lots:           totalLots,
		Spots:          totalSpots,
		AvailableSpots: availableSpots,
	}, nil
}
