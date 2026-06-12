package guests

import (
	"fmt"
	"strconv"
	"time"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type GuestsModule struct {
	*h.API
}

type okResponse struct{ ok bool }

type GuestWithReservationRequest struct {
	Guest       models.Guest        `json:"guest"`
	Reservation *ReservationRequest `json:"reservation,omitzero"`
	Payment     *models.Payment     `json:"payment,omitzero"`
	Companions  *[]CompanionRequest `json:"companions,omitzero"`
}

type CompanionRequest struct {
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	NationalID  string    `json:"nationalId"`
	FatherName  string    `json:"fatherName"`
	IDNumber    string    `json:"idNumber"`
	Gender      string    `json:"gender"`
	Relation    uint      `json:"relation"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Phone       string    `json:"phone"`

	NationalityID uint           `json:"nationalityID"`
	Nationality   models.Country `gorm:"foreignKey:NationalityID" json:"nationality,omitzero"`
}

type ReservationRequest struct {
	ReservationCode string    `json:"reservationCode"`
	EntryDate       time.Time `json:"entryDate"`
	DepartureDate   time.Time `json:"departureDate"`
	DurationOfStay  int       `json:"durationOfStay"`
	NumberOfPeople  int       `json:"numberOfPeople"`
	Origin          string    `json:"origin"`
	Destination     string    `json:"destination"`
	PurposeOfTravel string    `json:"purposeOfTravel"`
	Breakfast       bool      `json:"breakfast"`
	Parking         bool      `json:"parking"`
	FullBoard       bool      `json:"fullBoard"`
	RoomPrice       float64   `json:"roomPrice"`
	Notes           string    `json:"notes"`
	Rooms           []uint    `json:"rooms"`
}

type StayRequest struct {
	EntryDate       time.Time `json:"entryDate"`
	DepartureDate   time.Time `json:"departureDate"`
	DurationOfStay  int       `json:"durationOfStay"`
	NumberOfPeople  int       `json:"numberOfPeople"`
	Origin          string    `json:"origin"`
	Destination     string    `json:"destination"`
	PurposeOfTravel string    `json:"purposeOfTravel"`
	Breakfast       bool      `json:"breakfast"`
	Parking         bool      `json:"parking"`
	FullBoard       bool      `json:"fullBoard"`
	RoomPrice       float64   `json:"roomPrice"`
	Notes           string    `json:"notes"`
	Rooms           []uint    `json:"rooms"`
}

type GuestWithStayRequest struct {
	Guest      models.Guest        `json:"guest"`
	Stay       *StayRequest        `json:"stay,omitzero"`
	Payment    *models.Payment     `json:"payment,omitzero"`
	Companions *[]CompanionRequest `json:"companions,omitzero"`
}

type GuestWithStayResponse struct {
	Guest   models.Guest   `json:"guest"`
	Stay    models.Stay    `json:"stay"`
	Payment models.Payment `json:"payment"`
}

func (gm GuestsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	gm = GuestsModule{api}

	fuego.Get(
		s,
		"/",
		h.ListModel[models.Guest](
			api.Db,
			h.WithAllowedFilters("first_name", "last_name", "phone", "national_id", "id_number"),
		),
	)

	fuego.Get(s, "/archived", gm.getArchivedGuestsHandler)
	fuego.Post(s, "/", h.CreateModel[models.Guest](api.Db))
	fuego.Post(s, "/with-reservation", gm.createGuestWithReservation)
	fuego.Post(s, "/with-stay", gm.createGuestWithStay)
	fuego.Get(s, "/{id}", h.GetModel[models.Guest](api.Db))
	fuego.Put(s, "/{id}", h.UpdateModel[models.Guest](api.Db))
	fuego.Delete(s, "/{id}", h.DeleteModel[models.Guest](api.Db))

	fuego.Get(s, "/{id}/settle", gm.getGuestSettlementHandler)
	fuego.Post(s, "/{id}/settle", gm.settleGuestAccount)

	fuego.Get(s, "/relations", h.ListModel[models.FamilyRelationship](api.Db, h.WithTranslation[models.FamilyRelationship]()))
}

type GuestWithReservationResponse struct {
	Guest       models.Guest       `json:"guest"`
	Reservation models.Reservation `json:"reservation"`
	Payment     models.Payment     `json:"payment"`
}

type GuestSettlementResponse struct {
	Reservations    []ReservationSettlement `json:"reservations"`
	ParkingTxns     []ParkingSettlement     `json:"parkingTransactions"`
	RestaurantBills []RestaurantSettlement  `json:"restaurantBills"`
	TotalRoom       float64                 `json:"totalRoom"`
	TotalParking    float64                 `json:"totalParking"`
	TotalRestaurant float64                 `json:"totalRestaurant"`
	TotalDue        float64                 `json:"totalDue"`
	TotalPaid       float64                 `json:"totalPaid"`
	Balance         float64                 `json:"balance"`
}

type ReservationSettlement struct {
	ID              uint    `json:"id"`
	ReservationCode string  `json:"reservationCode"`
	CheckInDate     string  `json:"checkInDate"`
	CheckOutDate    string  `json:"checkOutDate"`
	Status          string  `json:"status"`
	StatusLabel     string  `json:"statusLabel"`
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

type RestaurantSettlement struct {
	ID          uint    `json:"id"`
	BillDate    string  `json:"billDate"`
	TotalAmount float64 `json:"totalAmount"`
	Notes       string  `json:"notes"`
	IsExternal  bool    `json:"isExternal"`
}

func resolveRoomStatusID(db *gorm.DB, slug string) uint {
	var status models.RoomStatus
	if err := db.Where("slug = ?", slug).First(&status).Error; err != nil {
		return 0
	}
	return status.ID
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

		if body.Companions != nil {
			for _, comp := range *body.Companions {
				companion := models.GuestCompanion{
					GuestID:       body.Guest.ID,
					FirstName:     comp.FirstName,
					LastName:      comp.LastName,
					NationalID:    comp.NationalID,
					IDNumber:      comp.IDNumber,
					RelationID:    comp.Relation,
					FatherName:    comp.FatherName,
					Gender:        comp.Gender,
					DateOfBirth:   comp.DateOfBirth,
					NationalityID: comp.NationalityID,
					Phone:         comp.Phone,
				}
				if err := tx.Create(&companion).Error; err != nil {
					return err
				}
			}
		}

		var reservation models.Reservation
		if body.Reservation != nil {
			rooms := []models.Room{}
			for _, v := range body.Reservation.Rooms {
				rooms = append(rooms, models.Room{Base: models.Base{ID: v}})
			}

			reservation = models.Reservation{
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
				FullBoard:       body.Reservation.FullBoard,
				Parking:         body.Reservation.Parking,
				RoomPrice:       body.Reservation.RoomPrice,
				Notes:           body.Reservation.Notes,
				Rooms:           rooms,
			}
			if err := tx.Create(&reservation).Error; err != nil {
				return err
			}

			// Mark rooms as reserved
			reservedStatusID := resolveRoomStatusID(tx, string(models.RoomStatusReserved))
			if reservedStatusID > 0 {
				for _, room := range rooms {
					tx.Model(&models.Room{}).Where("id = ?", room.ID).Update("status_id", reservedStatusID)
				}
			}
		}

		var payment models.Payment
		if body.Payment != nil {
			payment = *body.Payment
			payment.ID = 0
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
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

func (gm *GuestsModule) createGuestWithStay(c fuego.ContextWithBody[GuestWithStayRequest]) (GuestWithStayResponse, error) {
	var zero GuestWithStayResponse
	body, err := c.Body()
	if err != nil {
		return zero, err
	}

	var result GuestWithStayResponse

	err = gm.Db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&body.Guest).Error; err != nil {
			return err
		}

		if body.Companions != nil {
			for _, comp := range *body.Companions {
				companion := models.GuestCompanion{
					GuestID:       body.Guest.ID,
					FirstName:     comp.FirstName,
					LastName:      comp.LastName,
					NationalID:    comp.NationalID,
					IDNumber:      comp.IDNumber,
					RelationID:    comp.Relation,
					FatherName:    comp.FatherName,
					Gender:        comp.Gender,
					DateOfBirth:   comp.DateOfBirth,
					NationalityID: comp.NationalityID,
					Phone:         comp.Phone,
				}
				if err := tx.Create(&companion).Error; err != nil {
					return err
				}
			}
		}

		var stay models.Stay
		if body.Stay != nil {
			roomID := uint(0)
			if len(body.Stay.Rooms) > 0 {
				roomID = body.Stay.Rooms[0]
			}
			if roomID == 0 {
				return fuego.BadRequestError{Title: "no_room_assigned"}
			}

			// Get resident status
			var residentStatus models.StayStatus
			if err := tx.Where("slug = ?", "resident").First(&residentStatus).Error; err != nil {
				return err
			}

			stay = models.Stay{
				GuestID:         body.Guest.ID,
				RoomID:          roomID,
				AcceptanceID:    fmt.Sprintf("STY-%d", time.Now().Unix()),
				StayType:        string(models.StayTypeNormal),
				EntryDate:       body.Stay.EntryDate,
				DepartureDate:   body.Stay.DepartureDate,
				DurationOfStay:  body.Stay.DurationOfStay,
				NumberOfPeople:  body.Stay.NumberOfPeople,
				Origin:          body.Stay.Origin,
				Destination:     body.Stay.Destination,
				PurposeOfTravel: body.Stay.PurposeOfTravel,
				RoomPrice:       body.Stay.RoomPrice,
				Breakfast:       body.Stay.Breakfast,
				HalfBoard:       body.Stay.FullBoard,
				FullBoard:       body.Stay.FullBoard,
				Parking:         body.Stay.Parking,
				Notes:           body.Stay.Notes,
				StatusID:        residentStatus.ID,
				ActualCheckIn:   func() *time.Time { t := time.Now().UTC(); return &t }(),
			}
			if err := tx.Create(&stay).Error; err != nil {
				return err
			}

			// Mark room as occupied
			occupiedStatusID := resolveRoomStatusID(tx, string(models.RoomStatusOccupied))
			if occupiedStatusID > 0 {
				tx.Model(&models.Room{}).Where("id = ?", roomID).Update("status_id", occupiedStatusID)
			}

			// Generate invoice
			_ = gm.generateInvoiceForStay(tx, &stay)
		}

		var payment models.Payment
		if body.Payment != nil {
			payment = *body.Payment
			payment.ID = 0
			if err := tx.Create(&payment).Error; err != nil {
				return err
			}
		}

		result.Guest = body.Guest
		result.Stay = stay
		result.Payment = payment

		return nil
	})
	if err != nil {
		return zero, fuego.BadRequestError{Title: "create_failed"}
	}

	return result, nil
}

func (gm *GuestsModule) generateInvoiceForStay(tx *gorm.DB, stay *models.Stay) error {
	if stay.ID == 0 {
		return fmt.Errorf("stay not saved")
	}

	nights := stay.DurationOfStay
	if nights <= 0 {
		nights = 1
	}

	invoice := models.Invoice{
		StayID:          stay.ID,
		HotelID:         stay.HotelID,
		TotalAmount:     0,
		PaidAmount:      0,
		RemainingAmount: 0,
		PaymentStatus:   string(models.PaymentStatusUnpaid),
	}
	if err := tx.Create(&invoice).Error; err != nil {
		return err
	}

	var items []models.InvoiceItem

	roomTotal := stay.RoomPrice * float64(nights)
	items = append(items, models.InvoiceItem{
		InvoiceID:       invoice.ID,
		StayID:          stay.ID,
		ItemType:        string(models.InvoiceItemTypeRoomCharge),
		Quantity:        nights,
		UnitPrice:       stay.RoomPrice,
		TotalPrice:      roomTotal,
		Description:     "Room charge",
		RemainingAmount: roomTotal,
		PaymentStatus:   string(models.PaymentStatusUnpaid),
	})
	invoice.TotalAmount += roomTotal

	if stay.Breakfast {
		breakfastTotal := 10.0 * float64(nights)
		items = append(items, models.InvoiceItem{
			InvoiceID:       invoice.ID,
			StayID:          stay.ID,
			ItemType:        string(models.InvoiceItemTypeBreakfast),
			Quantity:        nights,
			UnitPrice:       10.0,
			TotalPrice:      breakfastTotal,
			Description:     "Breakfast",
			RemainingAmount: breakfastTotal,
			PaymentStatus:   string(models.PaymentStatusUnpaid),
		})
		invoice.TotalAmount += breakfastTotal
	}

	if stay.HalfBoard {
		halfBoardTotal := 20.0 * float64(nights)
		items = append(items, models.InvoiceItem{
			InvoiceID:       invoice.ID,
			StayID:          stay.ID,
			ItemType:        string(models.InvoiceItemTypeHalfBoard),
			Quantity:        nights,
			UnitPrice:       20.0,
			TotalPrice:      halfBoardTotal,
			Description:     "Half board",
			RemainingAmount: halfBoardTotal,
			PaymentStatus:   string(models.PaymentStatusUnpaid),
		})
		invoice.TotalAmount += halfBoardTotal
	}

	if stay.FullBoard {
		fullBoardTotal := 35.0 * float64(nights)
		items = append(items, models.InvoiceItem{
			InvoiceID:       invoice.ID,
			StayID:          stay.ID,
			ItemType:        string(models.InvoiceItemTypeFullBoard),
			Quantity:        nights,
			UnitPrice:       35.0,
			TotalPrice:      fullBoardTotal,
			Description:     "Full board",
			RemainingAmount: fullBoardTotal,
			PaymentStatus:   string(models.PaymentStatusUnpaid),
		})
		invoice.TotalAmount += fullBoardTotal
	}

	if stay.Parking {
		parkingTotal := 5.0 * float64(nights)
		items = append(items, models.InvoiceItem{
			InvoiceID:       invoice.ID,
			StayID:          stay.ID,
			ItemType:        string(models.InvoiceItemTypeParking),
			Quantity:        nights,
			UnitPrice:       5.0,
			TotalPrice:      parkingTotal,
			Description:     "Parking",
			RemainingAmount: parkingTotal,
			PaymentStatus:   string(models.PaymentStatusUnpaid),
		})
		invoice.TotalAmount += parkingTotal
	}

	if len(items) > 0 {
		if err := tx.CreateInBatches(items, len(items)).Error; err != nil {
			return err
		}
	}

	invoice.RemainingAmount = invoice.TotalAmount
	return tx.Save(&invoice).Error
}

type SettleGuestRequest struct {
	ReservationIDs    []uint  `json:"reservationIds"`
	ParkingTxnIDs     []uint  `json:"parkingTxnIds"`
	RestaurantBillIDs []uint  `json:"restaurantBillIds"`
	Amount            float64 `json:"amount"`
	PaymentMethod     uint    `json:"paymentMethod"`
	Reference         string  `json:"reference"`
	Notes             string  `json:"notes"`
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

	receivedStatusID := resolvePaymentStatus(gm.Db, "received")

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
				IncomeDate:      time.Now().UTC(),
				Description:     "Room payment - Reservation " + res.ReservationCode,
				Amount:          res.RoomPrice,
				Source:          guest.FirstName + " " + guest.LastName,
				PaymentMethodID: body.PaymentMethod,
				Notes:           body.Notes,
			}
			income.CategoryID = resolveIncomeCategory(tx, "room_revenue")
			income.PaymentStatusID = receivedStatusID
			income.AccountID = resolveAccount(tx, "1000")
			income.ReservationID = &resID
			if err := tx.Create(&income).Error; err != nil {
				return err
			}

			var payment models.Payment
			if err := tx.Where("reservation_id = ?", resID).First(&payment).Error; err != nil {
				payment = models.Payment{
					ReservationID: resID,
					Amount:        res.RoomPrice,
					MethodID:      body.PaymentMethod,
					StatusID:      receivedStatusID,
				}
				if err := tx.Create(&payment).Error; err != nil {
					return err
				}
			} else {
				payment.Amount += res.RoomPrice
				payment.MethodID = body.PaymentMethod
				payment.StatusID = receivedStatusID
				if err := tx.Save(&payment).Error; err != nil {
					return err
				}
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
			pt.PaymentMethodID = &body.PaymentMethod
			if err := tx.Save(&pt).Error; err != nil {
				return err
			}
			income := models.Income{
				IncomeDate:      time.Now().UTC(),
				Description:     "Parking payment - " + pt.LicensePlate,
				Amount:          pt.AmountDue,
				Source:          guest.FirstName + " " + guest.LastName,
				PaymentMethodID: body.PaymentMethod,
				Notes:           body.Notes,
			}
			income.CategoryID = resolveIncomeCategory(tx, "parking")
			income.PaymentStatusID = receivedStatusID
			income.AccountID = resolveAccount(tx, "1000")
			if err := tx.Create(&income).Error; err != nil {
				return err
			}
		}

		for _, billID := range body.RestaurantBillIDs {
			var bill models.RestaurantBill
			if err := tx.First(&bill, billID).Error; err != nil {
				continue
			}
			if bill.GuestID == nil || *bill.GuestID != id {
				continue
			}

			now := time.Now().UTC()
			bill.Settled = true
			bill.SettledAt = &now
			if err := tx.Save(&bill).Error; err != nil {
				return err
			}

			income := models.Income{
				IncomeDate:      now,
				Description:     "Restaurant payment - Bill #" + formatID(bill.ID),
				Amount:          bill.TotalAmount,
				Source:          guest.FirstName + " " + guest.LastName,
				PaymentMethodID: body.PaymentMethod,
				Notes:           body.Notes,
			}
			income.CategoryID = resolveIncomeCategory(tx, "restaurant")
			income.PaymentStatusID = receivedStatusID
			income.AccountID = resolveAccount(tx, "1000")
			if bill.ReservationID != nil {
				income.ReservationID = bill.ReservationID
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

func formatID(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

func (gm *GuestsModule) getGuestSettlement(id uint) (GuestSettlementResponse, error) {
	var zero GuestSettlementResponse

	var reservations []models.Reservation
	if err := gm.Db.Preload("Status").Where("guest_id = ?", id).Find(&reservations).Error; err != nil {
		return zero, err
	}

	var parkingTxns []models.ParkingTransaction
	if err := gm.Db.Where("guest_id = ?", id).Find(&parkingTxns).Error; err != nil {
		return zero, err
	}

	var bills []models.RestaurantBill
	if err := gm.Db.Where("guest_id = ? AND settled = ?", id, false).Find(&bills).Error; err != nil {
		return zero, err
	}

	var resp GuestSettlementResponse
	resp.Reservations = make([]ReservationSettlement, 0, len(reservations))
	resp.ParkingTxns = make([]ParkingSettlement, 0, len(parkingTxns))
	resp.RestaurantBills = make([]RestaurantSettlement, 0, len(bills))

	var totalRoom, totalParking, totalRestaurant, totalPaid float64

	for _, r := range reservations {
		var paid float64
		gm.Db.Model(&models.Income{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("reservation_id = ?", r.ID).
			Scan(&paid)

		statusSlug := ""
		statusLabel := ""
		if r.Status.Slug != "" {
			statusSlug = r.Status.Slug
			statusLabel = r.Status.Label
		}

		s := ReservationSettlement{
			ID:              r.ID,
			ReservationCode: r.ReservationCode,
			CheckInDate:     r.EntryDate.Format(time.RFC3339),
			CheckOutDate:    r.DepartureDate.Format(time.RFC3339),
			Status:          statusSlug,
			StatusLabel:     statusLabel,
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

	for _, b := range bills {
		s := RestaurantSettlement{
			ID:          b.ID,
			BillDate:    b.BillDate.Format(time.RFC3339),
			TotalAmount: b.TotalAmount,
			Notes:       b.Notes,
			IsExternal:  b.IsExternal,
		}
		resp.RestaurantBills = append(resp.RestaurantBills, s)
		totalRestaurant += b.TotalAmount
	}

	resp.TotalRoom = totalRoom
	resp.TotalParking = totalParking
	resp.TotalRestaurant = totalRestaurant
	resp.TotalDue = totalRoom + totalParking + totalRestaurant
	resp.TotalPaid = totalPaid
	resp.Balance = resp.TotalDue - totalPaid

	return resp, nil
}

func (gm *GuestsModule) getGuestSettlementHandler(c fuego.ContextNoBody) (GuestSettlementResponse, error) {
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return GuestSettlementResponse{}, fuego.BadRequestError{Title: "invalid_id"}
	}

	var guest models.Guest
	if err := gm.Db.First(&guest, id).Error; err != nil {
		return GuestSettlementResponse{}, fuego.NotFoundError{}
	}

	return gm.getGuestSettlement(id)
}

func (gm *GuestsModule) getArchivedGuestsHandler(c fuego.ContextWithParams[h.Params]) (h.PaginatedResponse[models.Guest], error) {
	var zeroResponse h.PaginatedResponse[models.Guest]

	status := c.QueryParam("status")
	q := gm.Db.Model(&models.Guest{})

	if status != "" {
		q = q.Where("status = ?", status)
	} else {
		q = q.Where("status IN ?", []string{string(models.GuestStatusCheckedOut), string(models.GuestStatusCancelled), string(models.GuestStatusAbsence)})
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return zeroResponse, err
	}

	f := h.WithAllowedFilters("first_name", "last_name", "phone", "national_id", "id_number")
	f(c, q)

	return h.Paginate[models.Guest](c, q)
}

func resolveIncomeCategory(db *gorm.DB, slug string) uint {
	var cat models.IncomeCategory
	if err := db.Where("slug = ?", slug).First(&cat).Error; err != nil {
		return 1
	}
	return cat.ID
}

func resolvePaymentStatus(db *gorm.DB, slug string) uint {
	var status models.PaymentStatus
	if err := db.Where("slug = ?", slug).First(&status).Error; err != nil {
		return 1
	}
	return status.ID
}

func resolveAccount(db *gorm.DB, code string) *uint {
	var account models.Account
	if err := db.Where("account_code = ?", code).First(&account).Error; err != nil {
		return nil
	}
	return &account.ID
}
