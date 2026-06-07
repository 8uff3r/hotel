package stays

import (
	"fmt"
	"time"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type StaysModule struct {
	*h.API
}

func (m StaysModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	sm := StaysModule{api}

	fuego.Get(s, "/", sm.staysList)
	fuego.Post(s, "/", sm.staysCreate)
	fuego.Get(s, "/{id}", sm.stayView)
	fuego.Put(s, "/{id}", sm.stayUpdate)
	fuego.Post(s, "/{id}/check-in", sm.stayCheckIn)
	fuego.Post(s, "/{id}/check-out", sm.stayCheckOut)
	fuego.Post(s, "/{id}/change-room", sm.stayChangeRoom)
	fuego.Post(s, "/{id}/services", sm.stayAddService)
	fuego.Post(s, "/{id}/settle", sm.staySettle)
	fuego.Get(s, "/{id}/invoice", sm.stayGetInvoice)
	fuego.Post(s, "/{id}/invoice/items", sm.stayAddInvoiceItem)
	fuego.Post(s, "/{id}/invoice/pay", sm.stayPayInvoice)
}

func (sm *StaysModule) staysList(c fuego.ContextNoBody) (h.PaginatedResponse[models.Stay], error) {
	var rows []models.Stay
	page := max(c.QueryParamInt("page"), 1)
	limit := c.QueryParamInt("limit")
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	q := sm.Db.WithContext(c).Model(&models.Stay{}).Preload("Guest").Preload("Room").Preload("Status").Preload("TravelAgency")

	if status := c.QueryParam("status"); status != "" {
		q = q.Where("status_id = (SELECT id FROM stay_statuses WHERE slug = ?)", status)
	}
	if from := c.QueryParam("from"); from != "" {
		q = q.Where("entry_date >= ?", from)
	}
	if to := c.QueryParam("to"); to != "" {
		q = q.Where("departure_date <= ?", to)
	}

	var total int64
	q.Count(&total)

	if err := q.Order("id DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return h.PaginatedResponse[models.Stay]{}, fuego.InternalServerError{Title: "query_failed"}
	}
	return h.PaginatedResponse[models.Stay]{
		Data:       rows,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}, nil
}

func (sm *StaysModule) stayView(c fuego.ContextNoBody) (models.Stay, error) {
	var zero models.Stay
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}
	var stay models.Stay
	if err := sm.Db.WithContext(c).
		Preload("Guest").
		Preload("Guest.Companions").
		Preload("Guest.Companions.Relation").
		Preload("Guest.Nationality").
		Preload("Room").
		Preload("Room.Type").
		Preload("Room.Status").
		Preload("Status").
		Preload("TravelAgency").
		Preload("Invoice.Items").
		First(&stay, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return zero, fuego.NotFoundError{Title: "not_found"}
		}
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}
	return stay, nil
}

func (sm *StaysModule) staysCreate(c fuego.ContextWithBody[models.Stay]) (models.Stay, error) {
	var zero models.Stay
	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}
	body.AcceptanceID = generateAcceptanceID()
	if err := sm.Db.WithContext(c).Create(&body).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "create_failed"}
	}
	// Auto-generate invoice
	if err := sm.generateInvoiceForStay(&body); err != nil {
		// Log but don't fail
	}
	c.SetStatus(201)
	return body, nil
}

func (sm *StaysModule) stayUpdate(c fuego.ContextWithBody[models.Stay]) (models.Stay, error) {
	var zero models.Stay
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}
	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}
	body.ID = id
	res := sm.Db.WithContext(c).Model(&models.Stay{}).Where("id = ?", id).Updates(body)
	if res.Error != nil {
		return zero, fuego.BadRequestError{Title: "update_failed"}
	}
	if res.RowsAffected == 0 {
		return zero, fuego.NotFoundError{}
	}
	var updated models.Stay
	sm.Db.First(&updated, id)
	return updated, nil
}

func (sm *StaysModule) stayCheckIn(c fuego.ContextWithBody[map[string]any]) (models.Stay, error) {
	var zero models.Stay
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	var stay models.Stay
	if err := sm.Db.First(&stay, id).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}

	now := time.Now().UTC()
	stay.ActualCheckIn = &now

	// Check early check-in
	var hotel models.Hotel
	sm.Db.First(&hotel, stay.HotelID)
	if hotel.Setting != nil && hotel.Setting.StandardCheckInTime != "" {
		checkInTime, parseErr := time.Parse("15:04", hotel.Setting.StandardCheckInTime)
		if parseErr == nil {
			currentTime := time.Now()
			standardCheckIn := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), checkInTime.Hour(), checkInTime.Minute(), 0, 0, currentTime.Location())
			if currentTime.Before(standardCheckIn) {
				// Early check-in: could set stay_type or return a flag
				stay.StayType = string(models.StayTypeEarlyCheckIn)
			}
		}
	}

	// Update guest status to resident
	sm.Db.Model(&models.Guest{}).Where("id = ?", stay.GuestID).Update("status", string(models.GuestStatusResident))

	// Update room status to occupied
	var occupiedStatus models.RoomStatus
	sm.Db.Where("slug = ?", "occupied").First(&occupiedStatus)
	if occupiedStatus.ID > 0 {
		sm.Db.Model(&models.Room{}).Where("id = ?", stay.RoomID).Update("status_id", occupiedStatus.ID)
	}

	// Update stay status to resident
	var residentStatus models.StayStatus
	sm.Db.Where("slug = ?", "resident").First(&residentStatus)
	if residentStatus.ID > 0 {
		stay.StatusID = residentStatus.ID
	}

	if err := sm.Db.Save(&stay).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "check_in_failed"}
	}

	// Generate invoice if not exists
	var existingInvoice models.Invoice
	if err := sm.Db.Where("stay_id = ?", stay.ID).First(&existingInvoice).Error; err != nil {
		sm.generateInvoiceForStay(&stay)
	}

	return stay, nil
}

func (sm *StaysModule) stayCheckOut(c fuego.ContextWithBody[map[string]any]) (models.Stay, error) {
	var zero models.Stay
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	var stay models.Stay
	if err := sm.Db.Preload("Invoice").First(&stay, id).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}

	// Verify invoice is cleared
	if stay.Invoice != nil && stay.Invoice.RemainingAmount > 0 {
		return zero, fuego.BadRequestError{Title: "invoice_not_cleared"}
	}

	now := time.Now().UTC()
	stay.ActualCheckOut = &now

	// Update guest status
	sm.Db.Model(&models.Guest{}).Where("id = ?", stay.GuestID).Update("status", string(models.GuestStatusCheckedOut))

	// Update room status to cleaning
	var cleaningStatus models.RoomStatus
	sm.Db.Where("slug = ?", "cleaning").First(&cleaningStatus)
	if cleaningStatus.ID > 0 {
		sm.Db.Model(&models.Room{}).Where("id = ?", stay.RoomID).Update("status_id", cleaningStatus.ID)
	}

	// Update stay status
	var checkedOutStatus models.StayStatus
	sm.Db.Where("slug = ?", "checked_out").First(&checkedOutStatus)
	if checkedOutStatus.ID > 0 {
		stay.StatusID = checkedOutStatus.ID
	}

	if err := sm.Db.Save(&stay).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "check_out_failed"}
	}
	return stay, nil
}

type changeRoomDto struct {
	NewRoomID uint `json:"newRoomId"`
}

func (sm *StaysModule) stayChangeRoom(c fuego.ContextWithBody[changeRoomDto]) (models.Stay, error) {
	var zero models.Stay
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	var stay models.Stay
	if err := sm.Db.First(&stay, id).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}

	var newRoom models.Room
	if err := sm.Db.First(&newRoom, body.NewRoomID).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "room_not_found"}
	}

	if newRoom.Capacity < stay.NumberOfPeople {
		return zero, fuego.BadRequestError{Title: "room_capacity_insufficient"}
	}

	oldRoomID := stay.RoomID
	stay.RoomID = body.NewRoomID
	if err := sm.Db.Save(&stay).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "change_room_failed"}
	}

	// Update old room status if no other active stays
	var oldRoomStayCount int64
	sm.Db.Model(&models.Stay{}).Where("room_id = ? AND status_id IN (SELECT id FROM stay_statuses WHERE slug IN ('waiting', 'resident'))", oldRoomID).Count(&oldRoomStayCount)
	if oldRoomStayCount == 0 {
		var availableStatus models.RoomStatus
		sm.Db.Where("slug = ?", "available").First(&availableStatus)
		if availableStatus.ID > 0 {
			sm.Db.Model(&models.Room{}).Where("id = ?", oldRoomID).Update("status_id", availableStatus.ID)
		}
	}

	return stay, nil
}

type addServiceDto struct {
	ServiceID   uint    `json:"serviceId"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
}

func (sm *StaysModule) stayAddService(c fuego.ContextWithBody[addServiceDto]) (models.InvoiceItem, error) {
	var zero models.InvoiceItem
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	var stay models.Stay
	if err := sm.Db.First(&stay, id).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}

	var invoice models.Invoice
	if err := sm.Db.Where("stay_id = ?", stay.ID).First(&invoice).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "invoice_not_found"}
	}

	var service models.Service
	if err := sm.Db.First(&service, body.ServiceID).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "service_not_found"}
	}

	quantity := body.Quantity
	if quantity <= 0 {
		quantity = 1
	}
	totalPrice := service.BaseAmount * float64(quantity)

	item := models.InvoiceItem{
		InvoiceID:       invoice.ID,
		StayID:          stay.ID,
		ItemType:        string(models.InvoiceItemTypeRoomService),
		ServiceID:       &service.ID,
		Quantity:        quantity,
		UnitPrice:       service.BaseAmount,
		TotalPrice:      totalPrice,
		Description:     body.Description,
		RemainingAmount: totalPrice,
		PaymentStatus:   string(models.PaymentStatusUnpaid),
	}
	if err := sm.Db.Create(&item).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "create_failed"}
	}

	// Update invoice totals
	invoice.TotalAmount += totalPrice
	invoice.RemainingAmount += totalPrice
	sm.Db.Save(&invoice)

	return item, nil
}

type settleDto struct {
	Amount        float64 `json:"amount"`
	PaymentMethod uint    `json:"paymentMethod"`
}

func (sm *StaysModule) staySettle(c fuego.ContextWithBody[settleDto]) (models.Invoice, error) {
	var zero models.Invoice
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	var invoice models.Invoice
	if err := sm.Db.Where("stay_id = ?", id).First(&invoice).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}

	if body.Amount <= 0 {
		return zero, fuego.BadRequestError{Title: "invalid_amount"}
	}

	invoice.PaidAmount += body.Amount
	invoice.RemainingAmount = invoice.TotalAmount - invoice.PaidAmount
	if invoice.RemainingAmount <= 0 {
		invoice.PaymentStatus = string(models.PaymentStatusCleared)
		invoice.RemainingAmount = 0
	} else {
		invoice.PaymentStatus = string(models.PaymentStatusPartiallyPaid)
	}
	if body.PaymentMethod > 0 {
		invoice.PaymentMethodID = &body.PaymentMethod
	}

	if err := sm.Db.Save(&invoice).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "settlement_failed"}
	}
	return invoice, nil
}

func (sm *StaysModule) stayGetInvoice(c fuego.ContextNoBody) (models.Invoice, error) {
	var zero models.Invoice
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}
	var invoice models.Invoice
	if err := sm.Db.Where("stay_id = ?", id).Preload("Items.Service").Preload("PaymentMethod").First(&invoice).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}
	return invoice, nil
}

type addInvoiceItemDto struct {
	ItemType        string  `json:"itemType"`
	Description     string  `json:"description"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unitPrice"`
	TotalPrice      float64 `json:"totalPrice"`
}

func (sm *StaysModule) stayAddInvoiceItem(c fuego.ContextWithBody[addInvoiceItemDto]) (models.InvoiceItem, error) {
	var zero models.InvoiceItem
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	var invoice models.Invoice
	if err := sm.Db.Where("stay_id = ?", id).First(&invoice).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}

	item := models.InvoiceItem{
		InvoiceID:       invoice.ID,
		StayID:          id,
		ItemType:        body.ItemType,
		Quantity:        body.Quantity,
		UnitPrice:       body.UnitPrice,
		TotalPrice:      body.TotalPrice,
		Description:     body.Description,
		RemainingAmount: body.TotalPrice,
		PaymentStatus:   string(models.PaymentStatusUnpaid),
	}
	if err := sm.Db.Create(&item).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "create_failed"}
	}

	invoice.TotalAmount += body.TotalPrice
	invoice.RemainingAmount += body.TotalPrice
	sm.Db.Save(&invoice)

	return item, nil
}

type payInvoiceDto struct {
	Amount        float64 `json:"amount"`
	PaymentMethod uint    `json:"paymentMethod"`
}

func (sm *StaysModule) stayPayInvoice(c fuego.ContextWithBody[payInvoiceDto]) (models.Invoice, error) {
	var zero models.Invoice
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	var invoice models.Invoice
	if err := sm.Db.Where("stay_id = ?", id).First(&invoice).Error; err != nil {
		return zero, fuego.NotFoundError{}
	}

	if body.Amount <= 0 {
		return zero, fuego.BadRequestError{Title: "invalid_amount"}
	}

	invoice.PaidAmount += body.Amount
	invoice.RemainingAmount = invoice.TotalAmount - invoice.PaidAmount
	if invoice.RemainingAmount <= 0 {
		invoice.PaymentStatus = string(models.PaymentStatusCleared)
		invoice.RemainingAmount = 0
	} else {
		invoice.PaymentStatus = string(models.PaymentStatusPartiallyPaid)
	}
	if body.PaymentMethod > 0 {
		invoice.PaymentMethodID = &body.PaymentMethod
	}

	if err := sm.Db.Save(&invoice).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "payment_failed"}
	}
	return invoice, nil
}

func (sm *StaysModule) generateInvoiceForStay(stay *models.Stay) error {
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
	if err := sm.Db.Create(&invoice).Error; err != nil {
		return err
	}

	var items []models.InvoiceItem

	// Room charge
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
		breakfastTotal := 10.0 * float64(nights) // default rate
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
		if err := sm.Db.CreateInBatches(items, len(items)).Error; err != nil {
			return err
		}
	}

	invoice.RemainingAmount = invoice.TotalAmount
	return sm.Db.Save(&invoice).Error
}

func generateAcceptanceID() string {
	return fmt.Sprintf("STY-%d", time.Now().Unix())
}
