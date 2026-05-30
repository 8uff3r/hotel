package parking

import (
	"time"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type ParkingModule struct {
	*h.API
}

func (m ParkingModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	p := ParkingModule{api}

	fuego.Get(s, "/stats", p.parkingStats)

	lotsGroup := fuego.Group(s, "/lots")
	fuego.Get(lotsGroup, "/", h.ListModel[models.ParkingLot](api.Db))
	fuego.Post(lotsGroup, "/", h.CreateModel[models.ParkingLot](api.Db))
	fuego.Get(lotsGroup, "/{id}", h.GetModel[models.ParkingLot](api.Db))
	fuego.Put(lotsGroup, "/{id}", h.UpdateModel[models.ParkingLot](api.Db))
	fuego.Delete(lotsGroup, "/{id}", h.DeleteModel[models.ParkingLot](api.Db))

	spotsGroup := fuego.Group(s, "/spots")
	fuego.Get(spotsGroup, "/", h.ListModel[models.ParkingSpot](api.Db))
	fuego.Post(spotsGroup, "/", h.CreateModel[models.ParkingSpot](api.Db))
	fuego.Get(spotsGroup, "/{id}", h.GetModel[models.ParkingSpot](api.Db))
	fuego.Put(spotsGroup, "/{id}", h.UpdateModel[models.ParkingSpot](api.Db))
	fuego.Delete(spotsGroup, "/{id}", h.DeleteModel[models.ParkingSpot](api.Db))
	fuego.Get(
		spotsGroup,
		"/statuses",
		h.ListModel[models.ParkingSpotStatus](api.Db, h.WithTranslation[models.ParkingSpotStatus]()),
	)
	fuego.Get(
		spotsGroup,
		"/types",
		h.ListModel[models.ParkingSpotType](api.Db, h.WithTranslation[models.ParkingSpotType]()),
	)

	vehiclesGroup := fuego.Group(s, "/vehicles")
	fuego.Get(vehiclesGroup, "/", h.ListModel[models.Vehicle](api.Db))
	fuego.Post(vehiclesGroup, "/", h.CreateModel[models.Vehicle](api.Db))
	fuego.Get(vehiclesGroup, "/{id}", h.GetModel[models.Vehicle](api.Db))
	fuego.Put(vehiclesGroup, "/{id}", h.UpdateModel[models.Vehicle](api.Db))
	fuego.Delete(vehiclesGroup, "/{id}", h.DeleteModel[models.Vehicle](api.Db))

	transactionsGroup := fuego.Group(s, "/transactions")
	fuego.Get(transactionsGroup, "/", h.ListModel[models.ParkingTransaction](api.Db))
	fuego.Post(transactionsGroup, "/", h.CreateModel[models.ParkingTransaction](api.Db))
	fuego.Get(transactionsGroup, "/{id}", h.GetModel[models.ParkingTransaction](api.Db))
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
