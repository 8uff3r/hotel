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

	fuego.Get(
		s,
		"/",
		h.ListModel[models.Reservation](
			api.Db,
			h.WithPreload("Rooms", "Guest", "Payment", "Payment.PaymentStatus"),
			h.WithTranslation[models.Translation](),
			h.WithAllowedFilters("status", "payment_status", "entry_date", "departure_date"),
		),
	)
	fuego.Get(
		s,
		"/{id}/detailed",
		re.getReservationDetails,
	)

	fuego.Post(s, "/", h.CreateModel[models.Reservation](api.Db))
	fuego.Get(s, "/{id}", h.GetModel[models.Reservation](api.Db))
	fuego.Put(s, "/{id}", h.UpdateModel[models.Reservation](api.Db))

	fuego.Post(s, "/{id}/check-in", re.reservationsCheckIn)
	fuego.Post(s, "/{id}/check-out", re.reservationsCheckOut)
}

type okResponse struct{ ok bool }

type getReservationDetailsResponse struct {
	models.Reservation
	Guest models.Guest `json:"guest"`
}

func (re *ReservationModule) getReservationDetails(c fuego.ContextNoBody) (models.Reservation, error) {
	var zero models.Reservation
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, err
	}
	var entity models.Reservation
	if err := re.Db.WithContext(c).Model(new(models.Reservation)).Preload("Payment", "Payment.PaymentStatus", "Guest", "Rooms").Where("id = ?", id).First(&entity).Error; err != nil {
		return zero, err
	}

	lang := c.Header("Accept-Language")
	if lang == "" {
		lang = "fa"
	}
	out := []models.Reservation{entity}
	rooms := entity.Rooms
	guests := []models.Guest{entity.Guest}
	models.ApplyTranslations(&out, lang)
	models.ApplyFieldTranslations(&out, lang)
	models.ApplyTranslations(&rooms, lang)
	models.ApplyFieldTranslations(&rooms, lang)
	models.ApplyTranslations(&guests, lang)
	models.ApplyFieldTranslations(&guests, lang)

	resp := out[0]
	resp.Rooms = rooms
	resp.Guest = guests[0]

	return resp, nil
}

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
