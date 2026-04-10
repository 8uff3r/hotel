package system

import (
	h "hotel/backend/internal/httpapi"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(api *h.API, r chi.Router) {
	api.Mux.Get("/healthz", a.health)
	api.Mux.Get("/readyz", a.ready)
}
