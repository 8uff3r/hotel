package guests

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-chi/chi/v5"
)

type GuestsModule struct{}

func (m GuestsModule) RegisterRoutes(api *h.API, r chi.Router) {
	r.Get("/", api.ListModel(&models.Guest{}, nil))
	r.Post("/", api.CreateModel(&models.Guest{}))
	r.Get("/{id}", api.GetModel(&models.Guest{}, nil))
	r.Put("/{id}", api.UpdateModel(&models.Guest{}))
}
