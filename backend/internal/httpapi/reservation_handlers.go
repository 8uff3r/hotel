package httpapi

import (
	"net/http"

	"hotel/backend/internal/service"
)

func (a *API) reservationsCheckIn(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeErr(w, 400, "invalid_id")
		return
	}
	if err := a.services.Reservation.CheckIn(r.Context(), id); err != nil {
		if service.IsNotFound(err) {
			writeErr(w, 404, "not_found")
			return
		}
		writeErr(w, 500, "update_failed")
		return
	}
	writeJSON(w, 200, map[string]bool{"ok": true})
}

func (a *API) reservationsCheckOut(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeErr(w, 400, "invalid_id")
		return
	}
	if err := a.services.Reservation.CheckOut(r.Context(), id); err != nil {
		if service.IsNotFound(err) {
			writeErr(w, 404, "not_found")
			return
		}
		writeErr(w, 500, "update_failed")
		return
	}
	writeJSON(w, 200, map[string]bool{"ok": true})
}
