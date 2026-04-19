package guests

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type GuestsModule struct{}

func (m GuestsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(s, "/", h.ListModel(api.Db, models.Guest{}))
	fuego.Post(s, "/", h.CreateModel(api.Db, models.Guest{}))
	fuego.Get(s, "/{id}", h.GetModel(api.Db, models.Guest{}))
	fuego.Put(s, "/{id}", h.UpdateModel(api.Db, models.Guest{}))
}
