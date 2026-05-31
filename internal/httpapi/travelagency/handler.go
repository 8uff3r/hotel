package travelagency

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type TravelAgencyModule struct{}

func (m TravelAgencyModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(s, "/", h.ListModel[models.TravelAgency](api.Db))
	fuego.Post(s, "/", h.CreateModel[models.TravelAgency](api.Db))
	fuego.Get(s, "/{id}", h.GetModel[models.TravelAgency](api.Db))
	fuego.Put(s, "/{id}", h.UpdateModel[models.TravelAgency](api.Db))
	fuego.Delete(s, "/{id}", h.DeleteModel[models.TravelAgency](api.Db))
}
