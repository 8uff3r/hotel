package common

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type CommonModule struct{}

func (m CommonModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(s, "/countries", h.ListModel(api.Db, models.Country{}, h.WithTranslation()))
}
