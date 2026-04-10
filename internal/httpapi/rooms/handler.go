package rooms

import (
	h "hotel/backend/internal/httpapi"
	"hotel/backend/internal/models"
	"hotel/backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(api *h.API, r chi.Router) {
	r.Route("/rooms", func(r chi.Router) {
		r.Get("/", api.ListModel(&models.Room{}, &repository.ListOptions{
			Preload: []string{"Amenities"},
		}))
		r.Post("/", api.CreateModel(&models.Room{}))
		r.Get("/{id}", api.GetModel(&models.Room{}, &repository.GetOptions{Preload: []string{"Amenities"}}))
		r.Put("/{id}", api.UpdateModel(&models.Room{}))
		r.Delete("/{id}", api.DeleteModel(&models.Room{}))

		r.Get("/amenities", api.ListModel(&models.Amenity{}, nil))
	})
}
