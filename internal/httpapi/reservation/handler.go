package reservation

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"hotel/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ReservationModule struct {
	*h.API
}

func (m ReservationModule) RegisterRoutes(api *h.API, r chi.Router) {
	re := ReservationModule{api}

	r.Get("/", api.ListModel(&models.Reservation{}, nil))
	r.Post("/", api.CreateModel(&models.Reservation{}))
	r.Get("/{id}", api.GetModel(&models.Reservation{}, nil))
	r.Put("/{id}", api.UpdateModel(&models.Reservation{}))

	r.Post("/{id}/check-in", re.reservationsCheckIn)
	r.Post("/{id}/check-out", re.reservationsCheckOut)

}

func (re *ReservationModule) reservationsCheckIn(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseID(r.PathValue("id"))
	if err != nil {
		h.WriteErr(w, 400, "invalid_id")
		return
	}
	if err := re.Services.Reservation.CheckIn(r.Context(), id); err != nil {
		if service.IsNotFound(err) {
			h.WriteErr(w, 404, "not_found")
			return
		}
		h.WriteErr(w, 500, "update_failed")
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
	if err := re.Services.Reservation.CheckOut(r.Context(), id); err != nil {
		if service.IsNotFound(err) {
			h.WriteErr(w, 404, "not_found")
			return
		}
		h.WriteErr(w, 500, "update_failed")
		return
	}
	h.WriteJSON(w, 200, map[string]bool{"ok": true})
}
