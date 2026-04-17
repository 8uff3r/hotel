package rooms

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type RoomsModule struct{}

func (m RoomsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(s, "/", h.ListModel(api.Db, models.Room{}, []string{"Amenities"}))
	fuego.Post(s, "/", h.CreateModel(api.Db, models.Room{}))
	fuego.Get(s, "/{id}", h.GetModel(api.Db, models.Room{}, []string{"Amenities"}))
	fuego.Put(s, "/{id}", h.UpdateModel(api.Db, models.Room{}))
	fuego.Delete(s, "/{id}", h.DeleteModel(api.Db, models.Room{}))

	fuego.Get(s, "/amenities", h.ListModel(api.Db, models.Amenity{}, nil))
}
