package rooms

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-chi/chi/v5"
)

type RoomsModule struct{}

func (m RoomsModule) RegisterRoutes(api *h.API, r chi.Router) {
	r.Get("/", api.ListModel(models.Room{}, []string{"Amenities"}))
	r.Post("/", api.CreateModel(models.Room{}))
	r.Get("/{id}", api.GetModel(models.Room{}, []string{"Amenities"}))
	r.Put("/{id}", api.UpdateModel(models.Room{}))
	r.Delete("/{id}", api.DeleteModel(models.Room{}))

	r.Get("/amenities", api.ListModel(models.Amenity{}, nil))
}
