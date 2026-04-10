package parking

import (
	h "hotel/backend/internal/httpapi"
	"hotel/backend/internal/models"
	"hotel/backend/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ParkingModule struct {
	*h.API
}

func (m ParkingModule) RegisterRoutes(api *h.API, r chi.Router) {
	p := ParkingModule{api}

	r.Get("/stats", p.parkingStats)

	r.Route("/lots", func(r chi.Router) {
		r.Get("/", api.ListModel(&models.ParkingLot{}, nil))
		r.Post("/", api.CreateModel(&models.ParkingLot{}))
		r.Get("/{id}", api.GetModel(&models.ParkingLot{}, nil))
		r.Put("/{id}", api.UpdateModel(&models.ParkingLot{}))
		r.Delete("/{id}", api.DeleteModel(&models.ParkingLot{}))
	})

	r.Route("/spots", func(r chi.Router) {
		r.Get("/", api.ListModel(&models.ParkingSpot{}, nil))
		r.Post("/", api.CreateModel(&models.ParkingSpot{}))
		r.Get("/{id}", api.GetModel(&models.ParkingSpot{}, nil))
		r.Put("/{id}", api.UpdateModel(&models.ParkingSpot{}))
		r.Delete("/{id}", api.DeleteModel(&models.ParkingSpot{}))

		r.Get("/statuses", api.ListModel(&models.ParkingSpotStatus{}, nil))
		r.Get("/types", api.ListModel(&models.ParkingSpotType{}, nil))
	})

	r.Route("/vehicles", func(r chi.Router) {
		r.Get("/", api.ListModel(&models.Vehicle{}, nil))
		r.Post("/", api.CreateModel(&models.Vehicle{}))
		r.Get("/{id}", api.GetModel(&models.Vehicle{}, nil))
		r.Put("/{id}", api.UpdateModel(&models.Vehicle{}))
		r.Delete("/{id}", api.DeleteModel(&models.Vehicle{}))
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", api.ListModel(&models.ParkingTransaction{}, nil))
		r.Post("/", api.CreateModel(&models.ParkingTransaction{}))
		r.Get("/{id}", api.GetModel(&models.ParkingTransaction{}, nil))
		r.Post("/{id}/check-out", p.transactionsCheckOut)
	})
}

func (a *ParkingModule) transactionsCheckOut(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseID(r.PathValue("id"))
	if err != nil {
		h.WriteErr(w, 400, "invalid_id")
		return
	}
	if err := a.Services.ParkingTx.CheckOut(r.Context(), id); err != nil {
		if service.IsNotFound(err) {
			h.WriteErr(w, 404, "not_found")
			return
		}
		h.WriteErr(w, 500, "update_failed")
		return
	}
	h.WriteJSON(w, 200, map[string]bool{"ok": true})
}

func (a *ParkingModule) parkingStats(w http.ResponseWriter, r *http.Request) {
	stats, err := a.Services.Parking.Stats(r.Context())
	if err != nil {
		h.WriteErr(w, 500, "failed")
		return
	}
	h.WriteJSON(w, 201, stats)
}
