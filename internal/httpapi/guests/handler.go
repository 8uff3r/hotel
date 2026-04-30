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

type okResponse struct{ ok bool }

type GuestWithReservationRequest struct {
	Guest       models.Guest       `json:"guest"`
	Reservation ReservationRequest `json:"reservation"`
	Payment     PaymentRequest     `json:"payment"`
	Companions  []CompanionRequest `json:"companions"`
}

type CompanionRequest struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	NationalID string `json:"nationalId"`
	IDNumber   string `json:"idNumber"`
	Relation   string `json:"relation"`
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

	fuego.Get(s, "/{id}/settle", gm.getGuestSettlementHandler)
	fuego.Post(s, "/{id}/settle", gm.settleGuestAccount)
}

type GuestWithReservationResponse struct {
	Guest       models.Guest       `json:"guest"`
	Reservation models.Reservation `json:"reservation"`
	Payment     models.Payment     `json:"payment"`
}

type GuestSettlementResponse struct {
	Reservations []ReservationSettlement `json:"reservations"`
	ParkingTxns  []ParkingSettlement     `json:"parkingTransactions"`
	TotalRoom    float64                 `json:"totalRoom"`
	TotalParking float64                 `json:"totalParking"`
	TotalDue     float64                 `json:"totalDue"`
	TotalPaid    float64                 `json:"totalPaid"`
	Balance      float64                 `json:"balance"`
}

type ReservationSettlement struct {
	ID              uint    `json:"id"`
	ReservationCode string  `json:"reservationCode"`
	CheckInDate     string  `json:"checkInDate"`
	CheckOutDate    string  `json:"checkOutDate"`
	Status          string  `json:"status"`
	RoomPrice       float64 `json:"roomPrice"`
	PaidAmount      float64 `json:"paidAmount"`
}

type ParkingSettlement struct {
	ID           uint    `json:"id"`
	LicensePlate string  `json:"licensePlate"`
	EntryTime    string  `json:"entryTime"`
	ExitTime     string  `json:"exitTime"`
	HoursParked  float64 `json:"hoursParked"`
	RateApplied  float64 `json:"rateApplied"`
	AmountDue    float64 `json:"amountDue"`
	AmountPaid   float64 `json:"amountPaid"`
	Status       string  `json:"status"`
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

		for _, comp := range body.Companions {
			companion := models.GuestCompanion{
				GuestID:    body.Guest.ID,
				FirstName:  comp.FirstName,
				LastName:   comp.LastName,
				NationalID: comp.NationalID,
				IDNumber:   comp.IDNumber,
				Relation:   comp.Relation,
			}
			if err := tx.Create(&companion).Error; err != nil {
				return err
			}
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

type SettleGuestRequest struct {
	ReservationIDs []uint  `json:"reservationIds"`
	ParkingTxnIDs  []uint  `json:"parkingTxnIds"`
	Amount         float64 `json:"amount"`
	PaymentMethod  string  `json:"paymentMethod"`
	Reference      string  `json:"reference"`
	Notes          string  `json:"notes"`
}

func (gm *GuestsModule) settleGuestAccount(c fuego.ContextWithBody[SettleGuestRequest]) (okResponse, error) {
	var zero okResponse
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	body, err := c.Body()
	if err != nil {
		return zero, err
	}

	if body.Amount <= 0 {
		return zero, fuego.BadRequestError{Title: "invalid_amount"}
	}

	var guest models.Guest
	if err := gm.Db.First(&guest, id).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}

	err = gm.Db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		for _, resID := range body.ReservationIDs {
			var res models.Reservation
			if err := tx.First(&res, resID).Error; err != nil {
				continue
			}
			if res.GuestID != id {
				continue
			}
			income := models.Income{
				IncomeDate:    time.Now().UTC(),
				Description:   "Room payment - Reservation " + res.ReservationCode,
				Amount:        body.Amount * (res.RoomPrice / (res.RoomPrice + 1)),
				Category:      "room_revenue",
				Source:        guest.FirstName + " " + guest.LastName,
				PaymentMethod: body.PaymentMethod,
				PaymentStatus: "received",
				ReservationID: &resID,
				Notes:         body.Notes,
			}
			if err := tx.Create(&income).Error; err != nil {
				return err
			}
		}

		for _, pID := range body.ParkingTxnIDs {
			var pt models.ParkingTransaction
			if err := tx.First(&pt, pID).Error; err != nil {
				continue
			}
			if pt.GuestID == nil || *pt.GuestID != id {
				continue
			}
			pt.AmountPaid = pt.AmountDue
			pt.PaymentStatus = "paid"
			pt.PaymentMethod = body.PaymentMethod
			if err := tx.Save(&pt).Error; err != nil {
				return err
			}
			income := models.Income{
				IncomeDate:    time.Now().UTC(),
				Description:   "Parking payment - " + pt.LicensePlate,
				Amount:        pt.AmountDue,
				Category:      "parking",
				Source:        guest.FirstName + " " + guest.LastName,
				PaymentMethod: body.PaymentMethod,
				PaymentStatus: "received",
				Notes:         body.Notes,
			}
			if err := tx.Create(&income).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return zero, fuego.BadRequestError{Title: "settlement_failed"}
	}

	return okResponse{ok: true}, nil
}

func (gm *GuestsModule) getGuestSettlement(id uint) (GuestSettlementResponse, error) {
	var zero GuestSettlementResponse

	var reservations []models.Reservation
	if err := gm.Db.Where("guest_id = ?", id).Find(&reservations).Error; err != nil {
		return zero, err
	}

	var parkingTxns []models.ParkingTransaction
	if err := gm.Db.Where("guest_id = ?", id).Find(&parkingTxns).Error; err != nil {
		return zero, err
	}

	var resp GuestSettlementResponse
	resp.Reservations = make([]ReservationSettlement, 0, len(reservations))
	resp.ParkingTxns = make([]ParkingSettlement, 0, len(parkingTxns))

	var totalRoom, totalParking, totalPaid float64

	for _, r := range reservations {
		var paid float64
		var income models.Income
		gm.Db.Where("reservation_id = ?", r.ID).First(&income)
		paid = income.Amount

		s := ReservationSettlement{
			ID:              r.ID,
			ReservationCode: r.ReservationCode,
			CheckInDate:     r.EntryDate.Format(time.RFC3339),
			CheckOutDate:    r.DepartureDate.Format(time.RFC3339),
			Status:          r.UserCheckOut,
			RoomPrice:       r.RoomPrice,
			PaidAmount:      paid,
		}
		resp.Reservations = append(resp.Reservations, s)
		totalRoom += r.RoomPrice
		totalPaid += paid
	}

	for _, p := range parkingTxns {
		s := ParkingSettlement{
			ID:           p.ID,
			LicensePlate: p.LicensePlate,
			EntryTime:    p.EntryTime.Format(time.RFC3339),
			ExitTime:     "",
			HoursParked:  0,
			RateApplied:  0,
			AmountDue:    p.AmountDue,
			AmountPaid:   p.AmountPaid,
			Status:       p.Status,
		}
		if p.ExitTime != nil {
			s.ExitTime = p.ExitTime.Format(time.RFC3339)
		}
		if p.HoursParked != nil {
			s.HoursParked = *p.HoursParked
		}
		if p.RateApplied != nil {
			s.RateApplied = *p.RateApplied
		}
		resp.ParkingTxns = append(resp.ParkingTxns, s)
		totalParking += p.AmountDue
		totalPaid += p.AmountPaid
	}

	resp.TotalRoom = totalRoom
	resp.TotalParking = totalParking
	resp.TotalDue = totalRoom + totalParking
	resp.TotalPaid = totalPaid
	resp.Balance = (totalRoom + totalParking) - totalPaid

	return resp, nil
}

func (gm *GuestsModule) getGuestSettlementHandler(c fuego.ContextNoBody) (GuestSettlementResponse, error) {
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return GuestSettlementResponse{}, fuego.BadRequestError{Title: "invalid_id"}
	}
	return gm.getGuestSettlement(id)
}
