package guests

import (
	h "hotel/backend/internal/httpapi"
	"hotel/backend/internal/models"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(api *h.API, r chi.Router) {
	r.Route("/guests", func(r chi.Router) {
		r.Get("/", api.ListModel(&models.Guest{}, nil))
		r.Post("/", api.CreateModel(&models.Guest{}))
		r.Get("/{id}", api.GetModel(&models.Guest{}, nil))
		r.Put("/{id}", api.UpdateModel(&models.Guest{}))
	})
}
