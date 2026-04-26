package hotels

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type HotelsModule struct {
	*h.API
}

func (m HotelsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(s, "/", h.ListModel(api.Db, models.Hotel{}))
	fuego.Post(s, "/", h.CreateModel(api.Db, models.Hotel{}))
	fuego.Get(s, "/{id}", h.GetModel(api.Db, models.Hotel{}))
	fuego.Put(s, "/{id}", h.UpdateModel(api.Db, models.Hotel{}))
	fuego.Delete(s, "/{id}", h.DeleteModel(api.Db, models.Hotel{}))
}
