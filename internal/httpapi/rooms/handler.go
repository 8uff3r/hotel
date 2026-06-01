package rooms

import (
	"hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type RoomsModule struct{}

func (m RoomsModule) RegisterRoutes(api *httpapi.API, s *fuego.Server) {
	fuego.Get(
		s,
		"/",
		httpapi.ListModel[models.Room](
			api.Db,
			httpapi.WithPreload("Amenities", "Type", "Status"),
			httpapi.WithTranslation[models.Room](),
		),
	)
	fuego.Post(s, "/", httpapi.CreateModel[models.Room](api.Db))

	fuego.Get(s, "/rack", m.rackHandler(api))

	fuego.Get(
		s,
		"/{id}",
		httpapi.GetModel[models.Room](
			api.Db,
			"Amenities", "Type", "Status",
		),
	)
	fuego.Put(s, "/{id}", httpapi.UpdateModel[models.Room](api.Db))
	fuego.Delete(s, "/{id}", httpapi.DeleteModel[models.Room](api.Db))

	fuego.Get(
		s,
		"/amenities",
		httpapi.ListModel[models.Amenity](
			api.Db,
			httpapi.WithTranslation[models.Amenity](),
		),
	)
	fuego.Get(
		s,
		"/types",
		httpapi.ListModel[models.RoomType](
			api.Db,
			httpapi.WithTranslation[models.RoomType](),
		),
	)
	fuego.Get(
		s,
		"/statuses",
		httpapi.ListModel[models.RoomStatus](
			api.Db,
			httpapi.WithTranslation[models.RoomStatus](),
		),
	)
}
