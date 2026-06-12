package rooms

import (
	"hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type RackRoom struct {
	ID          uint               `json:"id"`
	RoomNumber  string             `json:"roomNumber"`
	Floor       int                `json:"floor"`
	Capacity    int                `json:"capacity"`
	BasePrice   float64            `json:"basePrice"`
	Description string             `json:"description"`
	Amenities   []models.Amenity   `json:"amenities,omitempty"`
	TypeID      uint               `json:"roomTypeId"`
	Type        *models.RoomType   `json:"roomType,omitempty"`
	StatusID    uint               `json:"statusId"`
	Status      *models.RoomStatus `json:"status,omitempty"`
}

func (m RoomsModule) rackHandler(api *httpapi.API) httpapi.FuegoHandler[httpapi.PaginatedResponse[RackRoom], any, httpapi.Params] {
	return func(c fuego.ContextWithParams[httpapi.Params]) (httpapi.PaginatedResponse[RackRoom], error) {
		var rooms []models.Room
		if err := api.Db.WithContext(c).
			Preload("Amenities").
			Preload("Type").
			Preload("Status").
			Preload("Floor").
			Find(&rooms).Error; err != nil {
			return httpapi.PaginatedResponse[RackRoom]{}, fuego.InternalServerError{Err: err}
		}

		result := make([]RackRoom, 0, len(rooms))
		for i := range rooms {
			room := rooms[i]
			floorNum := 0
			if room.Floor != nil {
				floorNum = room.Floor.Number
			}
			result = append(result, RackRoom{
				ID:          room.ID,
				RoomNumber:  room.RoomNumber,
				Floor:       floorNum,
				Capacity:    room.Capacity,
				BasePrice:   room.BasePrice,
				Description: room.Description,
				Amenities:   room.Amenities,
				TypeID:      room.TypeID,
				Type:        &room.Type,
				StatusID:    room.StatusID,
				Status:      &room.Status,
			})
		}

		lang := c.Header("Accept-Language")
		if lang == "" {
			lang = "fa"
		}
		applyRackTranslations(&result, lang)

		return httpapi.PaginatedResponse[RackRoom]{
			Data:  result,
			Total: int64(len(result)),
		}, nil
	}
}

func applyRackTranslations(rooms *[]RackRoom, lang string) {
	for i := range *rooms {
		r := &(*rooms)[i]
		if r.Type != nil {
			models.ApplyTranslationOnTranslatable(r.Type, lang)
		}
		if r.Status != nil {
			models.ApplyTranslationOnTranslatable(r.Status, lang)
		}
		for j := range r.Amenities {
			models.ApplyTranslationOnTranslatable(&r.Amenities[j], lang)
		}
	}
}
