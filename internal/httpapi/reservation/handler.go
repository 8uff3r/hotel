package reservation

import (
	"time"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type ReservationModule struct {
	*h.API
}

func (m ReservationModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	re := ReservationModule{api}

	fuego.Get(s, "/", h.ListModel[models.Reservation](api.Db))
	fuego.Post(s, "/", h.CreateModel[models.Reservation](api.Db))
	fuego.Get(s, "/{id}", h.GetModel[models.Reservation](api.Db))
	fuego.Put(s, "/{id}", h.UpdateModel[models.Reservation](api.Db))

	fuego.Post(s, "/{id}/check-in", re.reservationsCheckIn)
	fuego.Post(s, "/{id}/check-out", re.reservationsCheckOut)
}

type okResponse struct{ ok bool }

func (re *ReservationModule) reservationsCheckIn(c fuego.ContextNoBody) (okResponse, error) {
	var zero okResponse
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}
	res := re.Db.WithContext(c).Model(&models.Reservation{}).Where("id = ?", id).Updates(map[string]any{"status": "checked_in", "actual_check_in": time.Now().UTC()})
	if res.Error != nil {
		return zero, fuego.InternalServerError{Title: "update_failed"}
	}
	if res.RowsAffected == 0 {
		return zero, fuego.NotFoundError{}
	}
	return okResponse{ok: true}, nil
}

func (re *ReservationModule) reservationsCheckOut(c fuego.ContextNoBody) (okResponse, error) {
	var zero okResponse
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}
	res := re.Db.WithContext(c).Model(&models.Reservation{}).Where("id = ?", id).Updates(map[string]any{"status": "checked_out", "actual_check_out": time.Now().UTC()})
	if res.Error != nil {
		return zero, fuego.InternalServerError{Title: "update_failed"}
	}
	if res.RowsAffected == 0 {
		return zero, fuego.NotFoundError{}
	}
	return okResponse{ok: true}, nil
}
