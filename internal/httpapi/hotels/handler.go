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
	fuego.Get(s, "/", h.ListModel[models.Hotel](api.Db))
	fuego.Post(s, "/", h.CreateModel[models.Hotel](api.Db))
	fuego.Get(s, "/{id}", h.GetModel[models.Hotel](api.Db))
	fuego.Put(s, "/{id}", h.UpdateModel[models.Hotel](api.Db))
	fuego.Delete(s, "/{id}", h.DeleteModel[models.Hotel](api.Db))
}
