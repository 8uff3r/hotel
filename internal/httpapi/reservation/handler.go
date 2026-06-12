package reservation

import (
	"fmt"
	"time"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type ReservationModule struct {
	*h.API
}

func (m ReservationModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	re := ReservationModule{api}

	fuego.Get(
		s,
		"/",
		h.ListModel[models.Reservation](
			api.Db,
			h.WithPreload("Rooms", "Guest", "Payment", "Payment.Status"),
			h.WithTranslation[models.Translation](),
			h.WithAllowedFilters("status", "payment_status", "entry_date", "departure_date"),
		),
	)
	fuego.Get(
		s,
		"/{id}/detailed",
		re.getReservationDetails,
	)

	fuego.Post(s, "/", re.reservationsCreate)
	fuego.Get(s, "/{id}", h.GetModel[models.Reservation](api.Db))
	fuego.Put(s, "/{id}", h.UpdateModel[models.Reservation](api.Db))

	fuego.Post(s, "/{id}/accept", re.reservationsAccept)
	fuego.Post(s, "/{id}/check-in", re.reservationsCheckIn)
	fuego.Post(s, "/{id}/check-out", re.reservationsCheckOut)

	fuego.Get(s, "/check-availability", re.checkReservationAvailability)
	fuego.Get(s, "/statuses", h.ListModel[models.ReservationStatus](api.Db, h.WithTranslation[models.ReservationStatus]()))
}

type okResponse struct{ ok bool }

type getReservationDetailsResponse struct {
	models.Reservation
	Guest models.Guest `json:"guest"`
}

func (re *ReservationModule) getReservationDetails(c fuego.ContextNoBody) (models.Reservation, error) {
	var zero models.Reservation
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, err
	}
	var entity models.Reservation
	if err := re.Db.WithContext(c).
		Model(new(models.Reservation)).
		Preload("Payment").
		Preload("Payment.Status").
		Preload("Guest").
		Preload("Rooms").
		Where("id = ?", id).
		First(&entity).
		Error; err != nil {
		return zero, err
	}

	lang := c.Header("Accept-Language")
	if lang == "" {
		lang = "fa"
	}
	out := []models.Reservation{entity}
	rooms := entity.Rooms
	guests := []models.Guest{entity.Guest}
	models.ApplyTranslations(&out, lang)
	models.ApplyFieldTranslations(&out, lang)
	models.ApplyTranslations(&rooms, lang)
	models.ApplyFieldTranslations(&rooms, lang)
	models.ApplyTranslations(&guests, lang)
	models.ApplyFieldTranslations(&guests, lang)

	resp := out[0]
	resp.Rooms = rooms
	resp.Guest = guests[0]

	return resp, nil
}

type createReservationResponse struct {
	models.Reservation
	Warnings []string `json:"warnings"`
}

func (re *ReservationModule) reservationsCreate(c fuego.ContextWithBody[models.Reservation]) (createReservationResponse, error) {
	var zero createReservationResponse
	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	var warnings []string

	// Validate room availability for each room
	for _, room := range body.Rooms {
		occupied, err := re.isRoomOccupied(room.ID)
		if err != nil {
			return zero, fuego.InternalServerError{Title: "availability_check_failed"}
		}
		if occupied {
			return zero, fuego.BadRequestError{Title: "room_occupied"}
		}
	}

	if err := re.Db.WithContext(c).Create(&body).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "create_failed"}
	}

	// Mark rooms as reserved
	var reservedStatus models.RoomStatus
	if err := re.Db.Where("slug = ?", string(models.RoomStatusReserved)).First(&reservedStatus).Error; err == nil {
		for _, room := range body.Rooms {
			re.Db.Model(&models.Room{}).Where("id = ?", room.ID).Update("status_id", reservedStatus.ID)
		}
	}

	c.SetStatus(201)
	return createReservationResponse{Reservation: body, Warnings: warnings}, nil
}

func (re *ReservationModule) isRoomOccupied(roomID uint) (bool, error) {
	var room models.Room
	if err := re.Db.Preload("Status").First(&room, roomID).Error; err != nil {
		return false, err
	}
	return room.Status.Slug == string(models.RoomStatusOccupied) || room.Status.Slug == string(models.RoomStatusReserved), nil
}

func (re *ReservationModule) reservationsAccept(c fuego.ContextNoBody) (models.Stay, error) {
	var zero models.Stay
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	var reservation models.Reservation
	if err := re.Db.WithContext(c).Preload("Guest").Preload("Rooms").First(&reservation, id).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}

	// Get accepted status
	var acceptedStatus models.ReservationStatus
	if err := re.Db.Where("slug = ?", "accepted").First(&acceptedStatus).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "status_not_found"}
	}
	reservation.StatusID = &acceptedStatus.ID
	re.Db.Save(&reservation)

	// Create stay from reservation
	return re.createStayFromReservation(c, &reservation)
}

func (re *ReservationModule) reservationsCheckIn(c fuego.ContextNoBody) (models.Stay, error) {
	var zero models.Stay
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	var reservation models.Reservation
	if err := re.Db.WithContext(c).Preload("Guest").Preload("Rooms").First(&reservation, id).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}

	// Check early check-in
	var hotel models.Hotel
	re.Db.First(&hotel, reservation.HotelID)
	if hotel.Setting != nil && hotel.Setting.StandardCheckInTime != "" {
		checkInTime, parseErr := time.Parse("15:04", hotel.Setting.StandardCheckInTime)
		if parseErr == nil {
			currentTime := time.Now()
			standardCheckIn := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), checkInTime.Hour(), checkInTime.Minute(), 0, 0, currentTime.Location())
			if currentTime.Before(standardCheckIn) {
				// Early check-in logic: could return a warning, but for now proceed
			}
		}
	}

	stay, err := re.createStayFromReservation(c, &reservation)
	if err != nil {
		return zero, err
	}

	// Update reservation to accepted
	var acceptedStatus models.ReservationStatus
	if err := re.Db.Where("slug = ?", "accepted").First(&acceptedStatus).Error; err == nil {
		re.Db.Model(&reservation).Update("status_id", acceptedStatus.ID)
	}

	// Update guest status
	re.Db.Model(&models.Guest{}).Where("id = ?", reservation.GuestID).Update("status", string(models.GuestStatusResident))

	return stay, nil
}

func (re *ReservationModule) createStayFromReservation(ctx fuego.ContextNoBody, reservation *models.Reservation) (models.Stay, error) {
	var zero models.Stay

	// Get resident status
	var residentStatus models.StayStatus
	if err := re.Db.Where("slug = ?", "resident").First(&residentStatus).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "stay_status_not_found"}
	}

	roomID := uint(0)
	if len(reservation.Rooms) > 0 {
		roomID = reservation.Rooms[0].ID
	}
	if roomID == 0 {
		return zero, fuego.BadRequestError{Title: "no_room_assigned"}
	}

	stay := models.Stay{
		HotelID:                *reservation.HotelID,
		GuestID:                reservation.GuestID,
		RoomID:                 roomID,
		ReservationID:          &reservation.ID,
		AcceptanceID:           fmt.Sprintf("STY-%d", time.Now().Unix()),
		StayType:               string(models.StayTypeNormal),
		EntryDate:              reservation.EntryDate,
		DepartureDate:          reservation.DepartureDate,
		ScheduledEntryDate:     reservation.EntryDate,
		ScheduledDepartureDate: reservation.DepartureDate,
		ActualCheckIn:          func() *time.Time { t := time.Now().UTC(); return &t }(),
		DurationOfStay:         reservation.DurationOfStay,
		NumberOfPeople:         reservation.NumberOfPeople,
		Origin:                 reservation.Origin,
		Destination:            reservation.Destination,
		PurposeOfTravel:        reservation.PurposeOfTravel,
		RoomPrice:              reservation.RoomPrice,
		Breakfast:              reservation.Breakfast,
		HalfBoard:              reservation.HalfBoard,
		FullBoard:              reservation.FullBoard,
		Parking:                reservation.Parking,
		Notes:                  reservation.Notes,
		StatusID:               residentStatus.ID,
	}

	if err := re.Db.WithContext(ctx).Create(&stay).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "stay_create_failed"}
	}

	// Generate invoice
	if err := re.generateInvoiceForStay(&stay); err != nil {
		// Log but don't fail
	}

	// Update room status to occupied
	var occupiedStatus models.RoomStatus
	if err := re.Db.Where("slug = ?", string(models.RoomStatusOccupied)).First(&occupiedStatus).Error; err == nil {
		re.Db.Model(&models.Room{}).Where("id = ?", roomID).Update("status_id", occupiedStatus.ID)
	}

	return stay, nil
}

func (re *ReservationModule) generateInvoiceForStay(stay *models.Stay) error {
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
	if err := re.Db.Create(&invoice).Error; err != nil {
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

	// Early check-in fee
	if stay.EarlyCheckInFee > 0 {
		items = append(items, models.InvoiceItem{
			InvoiceID:       invoice.ID,
			StayID:          stay.ID,
			ItemType:        string(models.InvoiceItemTypeOther),
			Quantity:        1,
			UnitPrice:       stay.EarlyCheckInFee,
			TotalPrice:      stay.EarlyCheckInFee,
			Description:     "Early check-in fee",
			RemainingAmount: stay.EarlyCheckInFee,
			PaymentStatus:   string(models.PaymentStatusUnpaid),
		})
		invoice.TotalAmount += stay.EarlyCheckInFee
	}

	// Half day fee
	if stay.HalfDayFee > 0 {
		items = append(items, models.InvoiceItem{
			InvoiceID:       invoice.ID,
			StayID:          stay.ID,
			ItemType:        string(models.InvoiceItemTypeOther),
			Quantity:        1,
			UnitPrice:       stay.HalfDayFee,
			TotalPrice:      stay.HalfDayFee,
			Description:     "Half day fee",
			RemainingAmount: stay.HalfDayFee,
			PaymentStatus:   string(models.PaymentStatusUnpaid),
		})
		invoice.TotalAmount += stay.HalfDayFee
	}

	if len(items) > 0 {
		if err := re.Db.CreateInBatches(items, len(items)).Error; err != nil {
			return err
		}
	}

	invoice.RemainingAmount = invoice.TotalAmount
	return re.Db.Save(&invoice).Error
}

func (re *ReservationModule) reservationsCheckOut(c fuego.ContextNoBody) (okResponse, error) {
	var zero okResponse
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}
	res := re.Db.WithContext(c).Model(&models.Reservation{}).Where("id = ?", id).Updates(map[string]any{"status": "checked_out"})
	if res.Error != nil {
		return zero, fuego.InternalServerError{Title: "update_failed"}
	}
	if res.RowsAffected == 0 {
		return zero, fuego.NotFoundError{}
	}
	return okResponse{ok: true}, nil
}

type availabilityResponse struct {
	Available bool   `json:"available"`
	Message   string `json:"message,omitempty"`
}

func (re *ReservationModule) checkReservationAvailability(c fuego.ContextNoBody) (availabilityResponse, error) {
	roomIDStr := c.QueryParam("roomId")

	roomID, err := h.ParseID(roomIDStr)
	if err != nil {
		return availabilityResponse{}, fuego.BadRequestError{Title: "invalid_room_id"}
	}

	occupied, err := re.isRoomOccupied(roomID)
	if err != nil {
		return availabilityResponse{}, fuego.InternalServerError{Title: "availability_check_failed"}
	}

	if occupied {
		return availabilityResponse{Available: false, Message: "Room is occupied for the selected date range"}, nil
	}
	return availabilityResponse{Available: true, Message: "Room is available"}, nil
}
