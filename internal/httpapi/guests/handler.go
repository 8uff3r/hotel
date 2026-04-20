package guests

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"time"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type GuestsModule struct {
	*h.API
}

type GuestWithReservationRequest struct {
	Guest       models.Guest       `json:"guest"`
	Reservation ReservationRequest `json:"reservation"`
	Payment     PaymentRequest     `json:"payment"`
}

type ReservationRequest struct {
	ReservationCode string        `json:"reservationCode"`
	EntryDate       time.Time     `json:"entryDate"`
	DepartureDate   time.Time     `json:"departureDate"`
	DurationOfStay  int           `json:"durationOfStay"`
	NumberOfPeople  int           `json:"numberOfPeople"`
	Origin          string        `json:"origin"`
	Destination     string        `json:"destination"`
	PurposeOfTravel string        `json:"purposeOfTravel"`
	Breakfast       bool          `json:"breakfast"`
	Guide           bool          `json:"guide"`
	RoomPrice       float64       `json:"roomPrice"`
	Notes           string        `json:"notes"`
	Rooms           []models.Room `json:"rooms"`
}

type PaymentRequest struct {
	IsCash       bool   `json:"isCash"`
	Agency       bool   `json:"agency"`
	Referrer     string `json:"referrer"`
	ContractType string `json:"contractType"`
}

func (m GuestsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	gm := GuestsModule{api}

	fuego.Get(s, "/", h.ListModel(api.Db, models.Guest{}))
	fuego.Post(s, "/", h.CreateModel(api.Db, models.Guest{}))
	fuego.Post(s, "/with-reservation", gm.createGuestWithReservation)
	fuego.Get(s, "/{id}", h.GetModel(api.Db, models.Guest{}))
	fuego.Put(s, "/{id}", h.UpdateModel(api.Db, models.Guest{}))
}

type GuestWithReservationResponse struct {
	Guest       models.Guest       `json:"guest"`
	Reservation models.Reservation `json:"reservation"`
	Payment     models.Payment     `json:"payment"`
}

func (gm *GuestsModule) createGuestWithReservation(c fuego.ContextWithBody[GuestWithReservationRequest]) (GuestWithReservationResponse, error) {
	var zero GuestWithReservationResponse
	body, err := c.Body()
	if err != nil {
		return zero, err
	}

	var result GuestWithReservationResponse

	err = gm.Db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&body.Guest).Error; err != nil {
			return err
		}

		reservation := models.Reservation{
			GuestID:         body.Guest.ID,
			ReservationCode: body.Reservation.ReservationCode,
			EntryDate:       body.Reservation.EntryDate,
			DepartureDate:   body.Reservation.DepartureDate,
			DurationOfStay:  body.Reservation.DurationOfStay,
			NumberOfPeople:  body.Reservation.NumberOfPeople,
			Origin:          body.Reservation.Origin,
			Destination:     body.Reservation.Destination,
			PurposeOfTravel: body.Reservation.PurposeOfTravel,
			Breakfast:       body.Reservation.Breakfast,
			Guide:           body.Reservation.Guide,
			RoomPrice:       body.Reservation.RoomPrice,
			Notes:           body.Reservation.Notes,
			Rooms:           body.Reservation.Rooms,
		}
		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}

		payment := models.Payment{
			ReservationID: reservation.ID,
			IsCash:        body.Payment.IsCash,
			Agency:        body.Payment.Agency,
			Referrer:      body.Payment.Referrer,
			ContractType:  body.Payment.ContractType,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		result.Guest = body.Guest
		result.Reservation = reservation
		result.Payment = payment

		return nil
	})

	if err != nil {
		return zero, fuego.BadRequestError{Title: "create_failed"}
	}

	return result, nil
}
