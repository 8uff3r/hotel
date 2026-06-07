package rooms

import (
	"time"

	"hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type RoomsModule struct{}

func (m RoomsModule) RegisterRoutes(api *httpapi.API, s *fuego.Server) {
	fuego.Get(
		s,
		"/",
		httpapi.ListModel[models.Room](
			api.Db,
			httpapi.WithPreload("Amenities", "Type", "Status", "Floor"),
			httpapi.WithTranslation[models.Room](),
		),
	)
	fuego.Post(s, "/", httpapi.CreateModel[models.Room](api.Db))

	fuego.Get(s, "/rack", m.rackHandler(api))
	fuego.Get(s, "/{id}/status", roomDynamicStatus(api.Db))

	fuego.Get(
		s,
		"/{id}",
		httpapi.GetModel[models.Room](
			api.Db,
			"Amenities", "Type", "Status", "Floor", "Pictures",
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
	fuego.Get(
		s,
		"/floors",
		httpapi.ListModel[models.Floor](api.Db),
	)
}

type RoomDynamicStatus struct {
	Status    string `json:"status"`
	StatusSlug string `json:"statusSlug"`
	StayID    *uint  `json:"stayId,omitempty"`
	ReservationID *uint `json:"reservationId,omitempty"`
}

func roomDynamicStatus(db *gorm.DB) func(c fuego.ContextNoBody) (RoomDynamicStatus, error) {
	return func(c fuego.ContextNoBody) (RoomDynamicStatus, error) {
		id, err := httpapi.ParseID(c.PathParam("id"))
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
