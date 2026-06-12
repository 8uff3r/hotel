package rooms

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoomsModule struct{}

func (m RoomsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(
		s,
		"/",
		h.ListModel[models.Room](
			api.Db,
			h.WithPreload("Amenities", "Type", "Status", "Floor"),
			h.WithTranslation[models.Room](),
		),
	)
	fuego.Post(s, "/", h.CreateModel[models.Room](api.Db))

	fuego.Get(s, "/rack", m.rackHandler(api))
	fuego.Get(s, "/{id}/status", roomStatus(api.Db))

	fuego.Get(
		s,
		"/{id}",
		h.GetModel[models.Room](
			api.Db,
			"Amenities", "Type", "Status", "Floor", "Pictures",
		),
	)
	fuego.Put(s, "/{id}", roomUpdate(api.Db))
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
	fuego.Get(
		s,
		"/floors",
		h.ListModel[models.Floor](api.Db),
	)
	fuego.Post(s, "/floors", createFloor(api.Db))

	fuego.Get(s, "/{id}/pictures", roomPicturesGet(api.Db))
	fuego.Post(s, "/{id}/pictures", roomPicturesAdd(api.Db))
	fuego.Delete(s, "/{id}/pictures/{pictureId}", roomPicturesDelete(api.Db))
}

func roomUpdate(db *gorm.DB) func(c fuego.ContextWithBody[models.Room]) (models.Room, error) {
	return func(c fuego.ContextWithBody[models.Room]) (models.Room, error) {
		var zero models.Room
		id, err := h.ParseID(c.PathParam("id"))
		if err != nil {
			return zero, fuego.BadRequestError{Title: "invalid_id"}
		}
		body, err := c.Body()
		if err != nil {
			return zero, err
		}
		body.ID = id

		// Resolve status slug if provided without ID
		if body.Status.Slug != "" && body.StatusID == 0 {
			var status models.RoomStatus
			if err := db.Where("slug = ?", body.Status.Slug).First(&status).Error; err != nil {
				return zero, fuego.BadRequestError{Title: "invalid_status_slug"}
			}
			body.StatusID = status.ID
			body.Status = status
		}

		res := db.Model(&models.Room{}).Where("id = ?", id).Clauses(clause.Returning{}).Updates(body).Scan(&zero)
		if res.Error != nil {
			return zero, fuego.BadRequestError{Title: "update_failed"}
		}
		if res.RowsAffected == 0 {
			return zero, fuego.NotFoundError{}
		}
		return zero, nil
	}
}

func roomStatus(db *gorm.DB) func(c fuego.ContextNoBody) (models.RoomStatus, error) {
	return func(c fuego.ContextNoBody) (models.RoomStatus, error) {
		id, err := h.ParseID(c.PathParam("id"))
		if err != nil {
			return models.RoomStatus{}, fuego.BadRequestError{Title: "invalid_id"}
		}
		var room models.Room
		if err := db.Preload("Status").First(&room, id).Error; err != nil {
			return models.RoomStatus{}, fuego.NotFoundError{}
		}
		return room.Status, nil
	}
}

type roomPicturesResponse struct {
	Data []models.RoomPicture `json:"data"`
}

func roomPicturesGet(db *gorm.DB) func(c fuego.ContextNoBody) (roomPicturesResponse, error) {
	return func(c fuego.ContextNoBody) (roomPicturesResponse, error) {
		id, err := h.ParseID(c.PathParam("id"))
		if err != nil {
			return roomPicturesResponse{}, fuego.BadRequestError{Title: "invalid_id"}
		}
		var pictures []models.RoomPicture
		if err := db.Where("room_id = ?", id).Find(&pictures).Error; err != nil {
			return roomPicturesResponse{}, fuego.InternalServerError{Title: "query_failed"}
		}
		return roomPicturesResponse{Data: pictures}, nil
	}
}

type roomPictureDto struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

func roomPicturesAdd(db *gorm.DB) func(c fuego.ContextWithBody[roomPictureDto]) (models.RoomPicture, error) {
	return func(c fuego.ContextWithBody[roomPictureDto]) (models.RoomPicture, error) {
		var zero models.RoomPicture
		id, err := h.ParseID(c.PathParam("id"))
		if err != nil {
			return zero, fuego.BadRequestError{Title: "invalid_id"}
		}
		body, err := c.Body()
		if err != nil {
			return zero, fuego.BadRequestError{}
		}
		if body.URL == "" {
			return zero, fuego.BadRequestError{Title: "url_required"}
		}
		picture := models.RoomPicture{
			RoomID:      id,
			URL:         body.URL,
			Description: body.Description,
		}
		if err := db.Create(&picture).Error; err != nil {
			return zero, fuego.BadRequestError{Title: "create_failed"}
		}
		return picture, nil
	}
}

type roomDeleteResponse struct {
	Ok bool `json:"ok"`
}

func roomPicturesDelete(db *gorm.DB) func(c fuego.ContextNoBody) (roomDeleteResponse, error) {
	return func(c fuego.ContextNoBody) (roomDeleteResponse, error) {
		var zero roomDeleteResponse
		pictureID := c.PathParam("pictureId")
		if pictureID == "" {
			return zero, fuego.BadRequestError{Title: "invalid_id"}
		}
		res := db.Delete(&models.RoomPicture{}, pictureID)
		if res.Error != nil {
			return zero, fuego.InternalServerError{Title: "delete_failed"}
		}
		if res.RowsAffected == 0 {
			return zero, fuego.NotFoundError{}
		}
		return roomDeleteResponse{Ok: true}, nil
	}
}

type floorCreateDto struct {
	Number      int    `json:"number" validate:"required"`
	Description string `json:"description"`
}

func createFloor(db *gorm.DB) func(c fuego.ContextWithBody[floorCreateDto]) (models.Floor, error) {
	return func(c fuego.ContextWithBody[floorCreateDto]) (models.Floor, error) {
		var zero models.Floor
		body, err := c.Body()
		if err != nil {
			return zero, fuego.BadRequestError{}
		}
		hotelID := h.GetHotelIDFromContext(c.Context())
		floor := models.Floor{
			HotelID:     hotelID,
			Number:      body.Number,
			Description: body.Description,
		}
		if err := db.Create(&floor).Error; err != nil {
			return zero, fuego.BadRequestError{Title: "create_failed"}
		}
		return floor, nil
	}
}
