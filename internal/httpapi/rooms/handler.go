package rooms

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type RoomsModule struct{}

func (m RoomsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(
		s,
		"/",
		h.ListModel[models.Room](
			api.Db,
			h.WithPreload("Amenities", "Type", "Status"),
			h.WithTranslation[models.Room](),
		),
	)
	fuego.Post(s, "/", h.CreateModel[models.Room](api.Db))
	fuego.Get(
		s,
		"/{id}",
		h.GetModel[models.Room](
			api.Db,
			h.WithPreload("Amenities", "Type", "Status"),
			h.WithTranslation[models.Room](),
		),
	)
	fuego.Put(s, "/{id}", h.UpdateModel[models.Room](api.Db))
	fuego.Delete(s, "/{id}", h.DeleteModel[models.Room](api.Db))

	fuego.Get(
		s,
		"/amenities",
		h.ListModel[models.Amenity](
			api.Db,
			h.WithTranslation[models.Amenity](),
		),
	)
	fuego.Get(
		s,
		"/types",
		h.ListModel[models.RoomType](
			api.Db,
			h.WithTranslation[models.RoomType](),
		),
	)
	fuego.Get(
		s,
		"/statuses",
		h.ListModel[models.RoomStatus](
			api.Db,
			h.WithTranslation[models.RoomStatus](),
		),
	)
}
