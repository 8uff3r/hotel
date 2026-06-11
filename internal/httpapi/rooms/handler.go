package rooms

import (
	"time"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
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
		fuego.Get(s, "/{id}/status", roomDynamicStatus(api.Db))

		fuego.Get(
			s,
			"/{id}",
			h.GetModel[models.Room](
				api.Db,
				"Amenities", "Type", "Status", "Floor", "Pictures",
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

type RoomDynamicStatus struct {
	Status    string `json:"status"`
	StatusSlug string `json:"statusSlug"`
	StayID    *uint  `json:"stayId,omitempty"`
	ReservationID *uint `json:"reservationId,omitempty"`
}

func roomDynamicStatus(db *gorm.DB) func(c fuego.ContextNoBody) (RoomDynamicStatus, error) {
	return func(c fuego.ContextNoBody) (RoomDynamicStatus, error) {
		id, err := h.ParseID(c.PathParam("id"))
		if err != nil {
			return RoomDynamicStatus{}, fuego.BadRequestError{Title: "invalid_id"}
		}

		dateStr := c.QueryParam("date")
		var date time.Time
		if dateStr != "" {
			date, _ = time.Parse("2006-01-02", dateStr)
		}
		if date.IsZero() {
			date = time.Now()
		}

		var room models.Room
		if err := db.First(&room, id).Error; err != nil {
			return RoomDynamicStatus{}, fuego.NotFoundError{}
		}

		// Manual overrides (cleaning, under_repair) take precedence
		var roomStatus models.RoomStatus
		if err := db.First(&roomStatus, room.StatusID).Error; err == nil {
			if roomStatus.Slug == "cleaning" || roomStatus.Slug == "under_repair" {
				return RoomDynamicStatus{Status: roomStatus.Label, StatusSlug: roomStatus.Slug}, nil
			}
		}

		// Check active stays for this room on the date
		var stay models.Stay
		if err := db.Where("room_id = ? AND entry_date <= ? AND (departure_date IS NULL OR departure_date >= ?) AND status_id IN (SELECT id FROM stay_statuses WHERE slug IN ('waiting', 'resident'))", id, date, date).
			Order("id DESC").
			First(&stay).Error; err == nil {
			return RoomDynamicStatus{Status: "Occupied", StatusSlug: "occupied", StayID: &stay.ID}, nil
		}

		// Check reservations for this room on the date
		var reservation models.Reservation
		if err := db.Joins("JOIN reservation_rooms ON reservation_rooms.reservation_id = reservations.id").
			Where("reservation_rooms.room_id = ? AND entry_date <= ? AND (departure_date IS NULL OR departure_date >= ?) AND status_id IN (SELECT id FROM reservation_statuses WHERE slug IN ('awaiting_payment', 'verified', 'accepted'))", id, date, date).
			Order("reservations.id DESC").
			First(&reservation).Error; err == nil {
			return RoomDynamicStatus{Status: "Reserved", StatusSlug: "reserved", ReservationID: &reservation.ID}, nil
		}

		return RoomDynamicStatus{Status: "Available", StatusSlug: "available"}, nil
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
