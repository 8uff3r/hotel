package reservation

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type ReservationModule struct {
	*h.API
}

func (m ReservationModule) RegisterRoutes(api *h.API, r chi.Router) {
	re := ReservationModule{api}

	r.Get("/", api.ListModel(models.Reservation{}, nil))
	r.Post("/", api.CreateModel(models.Reservation{}))
	r.Get("/{id}", api.GetModel(models.Reservation{}, nil))
	r.Put("/{id}", api.UpdateModel(models.Reservation{}))

	r.Post("/{id}/check-in", re.reservationsCheckIn)
	r.Post("/{id}/check-out", re.reservationsCheckOut)

}

func (re *ReservationModule) reservationsCheckIn(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseID(r.PathValue("id"))
	if err != nil {
		h.WriteErr(w, 400, "invalid_id")
		return
	}
	res := re.Db.WithContext(r.Context()).Model(&models.Reservation{}).Where("id = ?", id).Updates(map[string]any{"status": "checked_in", "actual_check_in": time.Now().UTC()})
	if res.Error != nil {
		h.WriteErr(w, 500, "update_failed")
		return
	}
	if res.RowsAffected == 0 {
		h.WriteErr(w, 404, "not_found")
		return
	}
	h.WriteJSON(w, 200, map[string]bool{"ok": true})
}

func (re *ReservationModule) reservationsCheckOut(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseID(r.PathValue("id"))
	if err != nil {
		h.WriteErr(w, 400, "invalid_id")
		return
	}
	res := re.Db.WithContext(r.Context()).Model(&models.Reservation{}).Where("id = ?", id).Updates(map[string]any{"status": "checked_out", "actual_check_out": time.Now().UTC()})
	if res.Error != nil {
		h.WriteErr(w, 500, "update_failed")
		return
	}
	if res.RowsAffected == 0 {
		h.WriteErr(w, 404, "not_found")
		return
	}
	h.WriteJSON(w, 200, map[string]bool{"ok": true})
}
