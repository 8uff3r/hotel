package rooms

import (
	"time"

	"hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type RackRoom struct {
	ID                 uint               `json:"id"`
	RoomNumber         string             `json:"roomNumber"`
	Floor              int                `json:"floor"`
	Capacity           int                `json:"capacity"`
	BasePrice          float64            `json:"basePrice"`
	Description        string             `json:"description"`
	Amenities          []models.Amenity   `json:"amenities,omitempty"`
	TypeID             uint               `json:"roomTypeId"`
	Type               *models.RoomType   `json:"roomType,omitempty"`
	StatusID           uint               `json:"statusId"`
	Status             *models.RoomStatus `json:"status,omitempty"`
	CurrentReservation *ReservationBrief  `json:"currentReservation,omitempty"`
}

type ReservationBrief struct {
	ID              uint       `json:"id"`
	EntryDate       time.Time  `json:"entryDate"`
	DepartureDate   time.Time  `json:"departureDate"`
	DurationOfStay  int        `json:"durationOfStay"`
	NumberOfPeople  int        `json:"numberOfPeople"`
	Origin          string     `json:"origin"`
	PurposeOfTravel string     `json:"purposeOfTravel"`
	Breakfast       bool       `json:"breakfast"`
	Parking         bool       `json:"parking"`
	FullBoard       bool       `json:"fullBoard"`
	RoomPrice       float64    `json:"roomPrice"`
	Notes           string     `json:"notes"`
	Guest           GuestBrief `json:"guest"`
}

type GuestBrief struct {
	ID           uint             `json:"id"`
	FirstName    string           `json:"firstName"`
	LastName     string           `json:"lastName"`
	FatherName   string           `json:"fatherName"`
	NationalID   string           `json:"nationalId"`
	IDNumber     string           `json:"idNumber"`
	Gender       string           `json:"gender"`
	DateOfBirth  time.Time        `json:"dateOfBirth"`
	PlaceOfBirth string           `json:"placeOfBirth"`
	Phone        string           `json:"phone"`
	Address      string           `json:"address"`
	Occupation   string           `json:"occupation"`
	Email        string           `json:"email"`
	Landline     string           `json:"landline"`
	Nationality  *models.Country  `json:"nationality,omitempty"`
	Companions   []CompanionBrief `json:"companions,omitempty"`
}

type CompanionBrief struct {
	ID        uint                       `json:"id"`
	FirstName string                     `json:"firstName"`
	LastName  string                     `json:"lastName"`
	Relation  *models.FamilyRelationship `json:"relation,omitempty"`
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

		var checkedInStatusID uint
		api.Db.Model(&models.ReservationStatus{}).
			Where("slug = ?", "checked_in").
			Select("id").
			Scan(&checkedInStatusID)

		var reservations []models.Reservation
		if err := api.Db.WithContext(c).
			Preload("Guest").
			Preload("Guest.Companions").
			Preload("Guest.Companions.Relation").
			Preload("Guest.Nationality").
			Preload("Rooms").
			Where("status_id = ?", checkedInStatusID).
			Find(&reservations).Error; err != nil {
			return httpapi.PaginatedResponse[RackRoom]{}, fuego.InternalServerError{Err: err}
		}

		resMap := make(map[uint]models.Reservation)
		for _, res := range reservations {
			for _, room := range res.Rooms {
				resMap[room.ID] = res
			}
		}

		result := make([]RackRoom, 0, len(rooms))
		for i := range rooms {
			room := rooms[i]
			floorNum := 0
			if room.Floor != nil {
				floorNum = room.Floor.Number
			}
			rr := RackRoom{
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
			}
			if res, ok := resMap[room.ID]; ok {
				companions := make([]CompanionBrief, len(res.Guest.Companions))
				for j := range res.Guest.Companions {
					companion := res.Guest.Companions[j]
					companions[j] = CompanionBrief{
						ID:        companion.ID,
						FirstName: companion.FirstName,
						LastName:  companion.LastName,
						Relation:  &companion.Relation,
					}
				}
				rr.CurrentReservation = &ReservationBrief{
					ID:              res.ID,
					EntryDate:       res.EntryDate,
					DepartureDate:   res.DepartureDate,
					DurationOfStay:  res.DurationOfStay,
					NumberOfPeople:  res.NumberOfPeople,
					Origin:          res.Origin,
					PurposeOfTravel: res.PurposeOfTravel,
					Breakfast:       res.Breakfast,
					Parking:         res.Parking,
					FullBoard:       res.FullBoard,
					RoomPrice:       res.RoomPrice,
					Notes:           res.Notes,
					Guest: GuestBrief{
						ID:           res.Guest.ID,
						FirstName:    res.Guest.FirstName,
						LastName:     res.Guest.LastName,
						FatherName:   res.Guest.FatherName,
						NationalID:   res.Guest.NationalID,
						IDNumber:     res.Guest.IDNumber,
						Gender:       res.Guest.Gender,
						DateOfBirth:  res.Guest.DateOfBirth,
						PlaceOfBirth: res.Guest.PlaceOfBirth,
						Phone:        res.Guest.Phone,
						Address:      res.Guest.Address,
						Occupation:   res.Guest.Occupation,
						Email:        res.Guest.Email,
						Landline:     res.Guest.Landline,
						Nationality:  &res.Guest.Nationality,
						Companions:   companions,
					},
				}
			}
			result = append(result, rr)
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
		if r.CurrentReservation != nil {
			g := &r.CurrentReservation.Guest
			if g.Nationality != nil {
				models.ApplyTranslationOnTranslatable(g.Nationality, lang)
			}
			for j := range g.Companions {
				if g.Companions[j].Relation != nil {
					models.ApplyTranslationOnTranslatable(g.Companions[j].Relation, lang)
				}
			}
		}
	}
}
