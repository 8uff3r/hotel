package sana

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"hotel/internal/config"
	"hotel/internal/httpapi"
	"hotel/internal/models"
	"hotel/internal/sana"

	"github.com/go-fuego/fuego"
	h "github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type SanaModule struct {
	db      *gorm.DB
	sanaCfg config.SanaConfig
}

func New(db *gorm.DB, sanaCfg config.SanaConfig) SanaModule {
	return SanaModule{db: db, sanaCfg: sanaCfg}
}

func (m SanaModule) RegisterRoutes(api *httpapi.API, s *fuego.Server) {
	h.Get(s, "/travel-reasons", m.getTravelReasons)
	h.Get(s, "/family-relationships", m.getFamilyRelationships)
	h.Get(s, "/nationalities", m.getNationalities)
	h.Get(s, "/countries", m.getCountries)
	h.Get(s, "/cities", m.getSanaCities)
	h.Get(s, "/occupations", m.getOccupations)

	h.Get(s, "/guests", m.getSanaGuests)
	h.Get(s, "/rooms", m.getSanaRooms)

	h.Post(s, "/guests/{id}/sync", m.syncGuest)
	h.Post(s, "/rooms/{id}/sync", m.syncRoom)

	h.Post(s, "/sync-all", m.syncAll)
}

func (m SanaModule) getTravelReasons(c h.ContextNoBody) ([]models.TravelReason, error) {
	var reasons []models.TravelReason
	err := m.db.Find(&reasons).Error
	return reasons, err
}

func (m SanaModule) getFamilyRelationships(c h.ContextNoBody) ([]models.FamilyRelationship, error) {
	var relationships []models.FamilyRelationship
	err := m.db.Find(&relationships).Error
	return relationships, err
}

func (m SanaModule) getNationalities(c h.ContextNoBody) ([]models.Nationality, error) {
	var nationalities []models.Nationality
	err := m.db.Find(&nationalities).Error
	return nationalities, err
}

func (m SanaModule) getCountries(c h.ContextNoBody) ([]models.Country, error) {
	var countries []models.Country
	err := m.db.Find(&countries).Error
	if err != nil {
		return nil, err
	}
	lang := c.Header("Accept-Language")
	if lang == "" {
		lang = "fa"
	}
	models.ApplyTranslations(&countries, lang)
	return countries, nil
}

func (m SanaModule) getSanaCities(c h.ContextNoBody) ([]models.SanaCity, error) {
	var cities []models.SanaCity
	err := m.db.Find(&cities).Error
	if err != nil {
		return nil, err
	}
	lang := c.Header("Accept-Language")
	if lang == "" {
		lang = "fa"
	}
	models.ApplyTranslations(&cities, lang)
	return cities, nil
}

func (m SanaModule) getOccupations(c h.ContextNoBody) ([]models.Occupation, error) {
	var occupations []models.Occupation
	err := m.db.Find(&occupations).Error
	return occupations, err
}

type SanaGuestResponse struct {
	ID              uint       `json:"id"`
	GuestID         uint       `json:"guestId"`
	Guest           *GuestInfo `json:"guest,omitempty"`
	RecordMosafer   int        `json:"recordMosafer"`
	ShomarePaziresh string     `json:"shomarePaziresh"`
	ShomareOtagh    string     `json:"shomareOtagh"`
	SyncTime        string     `json:"syncTime,omitempty"`
}

type GuestInfo struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	NationalID string `json:"nationalId"`
}

func getActiveStayRoomNumber(db *gorm.DB, guestID uint) string {
	var stay models.Stay
	if err := db.Joins("JOIN stay_statuses ON stays.status_id = stay_statuses.id").
		Where("stays.guest_id = ? AND stay_statuses.slug = ?", guestID, string(models.StayStatusResident)).
		First(&stay).Error; err != nil {
		return ""
	}
	var room models.Room
	if err := db.First(&room, stay.RoomID).Error; err != nil {
		return ""
	}
	return room.RoomNumber
}

func (m SanaModule) getSanaGuests(c h.ContextNoBody) ([]SanaGuestResponse, error) {
	// Find guests with active (resident) stays
	var guests []models.Guest
	err := m.db.
		Joins("JOIN stays ON stays.guest_id = guests.id").
		Joins("JOIN stay_statuses ON stays.status_id = stay_statuses.id").
		Where("stay_statuses.slug = ?", string(models.StayStatusResident)).
		Group("guests.id").
		Find(&guests).Error
	if err != nil {
		return nil, err
	}

	var response []SanaGuestResponse
	for _, g := range guests {
		gi := &GuestInfo{
			FirstName:  g.FirstName,
			LastName:   g.LastName,
			NationalID: g.NationalID,
		}

		// Find or create SanaGuest record
		var sanaGuest models.SanaGuest
		if err := m.db.Where("guest_id = ?", g.ID).First(&sanaGuest).Error; err != nil {
			sanaGuest.GuestID = g.ID
			sanaGuest.ShomareOtagh = getActiveStayRoomNumber(m.db, g.ID)
			m.db.Create(&sanaGuest)
		}

		// Refresh room number if not set
		if sanaGuest.ShomareOtagh == "" {
			roomNumber := getActiveStayRoomNumber(m.db, g.ID)
			if roomNumber != "" {
				sanaGuest.ShomareOtagh = roomNumber
				m.db.Save(&sanaGuest)
			}
		}

		syncTime := ""
		if !sanaGuest.SyncTime.IsZero() {
			syncTime = sanaGuest.SyncTime.Format("2006-01-02T15:04:05Z")
		}
		response = append(response, SanaGuestResponse{
			ID:              sanaGuest.ID,
			GuestID:         g.ID,
			Guest:           gi,
			RecordMosafer:   sanaGuest.RecordMosafer,
			ShomarePaziresh: sanaGuest.ShomarePaziresh,
			ShomareOtagh:    sanaGuest.ShomareOtagh,
			SyncTime:        syncTime,
		})
	}
	return response, nil
}

type SanaRoomResponse struct {
	ID           uint      `json:"id"`
	HotelID      string    `json:"hotelId"`
	Room         *RoomInfo `json:"room,omitempty"`
	Rac          string    `json:"rac"`
	LastSyncTime string    `json:"lastSyncTime,omitempty"`
	IsSynced     bool      `json:"isSynced"`
	LastError    string    `json:"lastError,omitempty"`
}

type RoomInfo struct {
	ID         uint   `json:"id"`
	RoomNumber string `json:"roomNumber"`
}

func (m SanaModule) getSanaRooms(c h.ContextNoBody) ([]SanaRoomResponse, error) {
	var rooms []models.Room
	err := m.db.Preload("Floor").Find(&rooms).Error
	if err != nil {
		return nil, err
	}

	var response []SanaRoomResponse
	for _, r := range rooms {
		ri := &RoomInfo{
			ID:         r.ID,
			RoomNumber: r.RoomNumber,
		}

		// Find or create SanaRoomRack record
		var roomRack models.SanaRoomRack
		if err := m.db.Where("room_id = ?", r.ID).First(&roomRack).Error; err != nil {
			roomRack.RoomID = &r.ID
			roomRack.HotelID = "default"
			roomRack.Rac = buildRoomRac(r)
			m.db.Create(&roomRack)
		}

		lastSyncTime := ""
		if !roomRack.LastSyncTime.IsZero() {
			lastSyncTime = roomRack.LastSyncTime.Format("2006-01-02T15:04:05Z")
		}
		response = append(response, SanaRoomResponse{
			ID:           roomRack.ID,
			HotelID:      roomRack.HotelID,
			Room:         ri,
			Rac:          roomRack.Rac,
			LastSyncTime: lastSyncTime,
			IsSynced:     roomRack.IsSynced,
			LastError:    roomRack.LastError,
		})
	}
	return response, nil
}

func buildRoomRac(room models.Room) string {
	floorNum := 0
	if room.Floor != nil {
		floorNum = room.Floor.Number
	}
	return fmt.Sprintf("%d:%s", floorNum, room.RoomNumber)
}

func (m SanaModule) syncGuest(c fuego.ContextNoBody) (SanaGuestResponse, error) {
	idStr := c.PathParam("id")
	sanaGuestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return SanaGuestResponse{}, fuego.BadRequestError{Title: "invalid_id"}
	}

	var sanaGuest models.SanaGuest
	if err := m.db.First(&sanaGuest, sanaGuestID).Error; err != nil {
		return SanaGuestResponse{}, err
	}

	var guest models.Guest
	if err := m.db.First(&guest, sanaGuest.GuestID).Error; err != nil {
		return SanaGuestResponse{}, err
	}

	client := sana.NewClient(sana.Config{
		KelidVahed:     m.sanaCfg.KelidVahed,
		KelidPeimankar: m.sanaCfg.KelidPeimankar,
		CodeVahed:      m.sanaCfg.CodeVahed,
	})

	gender := 1
	if guest.Gender == "female" {
		gender = 2
	}

	// Get active stay and room number
	var activeStay models.Stay
	var shomareOtagh string
	if err := m.db.Joins("JOIN stay_statuses ON stays.status_id = stay_statuses.id").
		Where("stays.guest_id = ? AND stay_statuses.slug = ?", guest.ID, string(models.StayStatusResident)).
		First(&activeStay).Error; err == nil {
		var room models.Room
		if err := m.db.First(&room, activeStay.RoomID).Error; err == nil {
			shomareOtagh = room.RoomNumber
		}
	}
	if shomareOtagh == "" {
		shomareOtagh = sanaGuest.ShomareOtagh
	}

	input := sana.SabtMosaferinInput{
		NameMosafer:      guest.FirstName,
		FamilMosafer:     guest.LastName,
		NamePedar:        guest.FatherName,
		ShomareShenasaee: guest.NationalID,
		TarikhTavalod:    guest.DateOfBirth.Format("2006/01/02"),
		ID_Jensiat:       gender,
		TedadHamrah:      0,
		MosafereKhareji:  false,
		ShomareOtagh:     shomareOtagh,
		NameKarbareSabt:  "system",
		ID_NoeDadeh:      sana.DataTypeRoomReception,
		Tozihat:          "",
		ShomarePaziresh:  activeStay.StayCode,
	}

	result, err := client.SabtMosaferin(input)
	if err != nil {
		return SanaGuestResponse{}, fuego.BadRequestError{Title: "sana_api_error", Detail: err.Error()}
	}

	if !result.IsOK {
		detail := result.RawMessage
		if result.CodeError > 0 {
			detail = fmt.Sprintf("SANA CodeError=%d: %s", result.CodeError, result.RawMessage)
		}
		return SanaGuestResponse{}, fuego.BadRequestError{Title: "sana_api_error", Detail: detail}
	}

	sanaGuest.RecordMosafer = result.RecordMosafer
	sanaGuest.ShomarePaziresh = result.ShomarePaziresh
	sanaGuest.ShomareOtagh = result.ShomareOtagh
	sanaGuest.SyncTime = time.Now()
	m.db.Save(&sanaGuest)

	syncTime := ""
	if !sanaGuest.SyncTime.IsZero() {
		syncTime = sanaGuest.SyncTime.Format("2006-01-02T15:04:05Z")
	}

	return SanaGuestResponse{
		ID:              sanaGuest.ID,
		GuestID:         guest.ID,
		RecordMosafer:   sanaGuest.RecordMosafer,
		ShomarePaziresh: sanaGuest.ShomarePaziresh,
		ShomareOtagh:    sanaGuest.ShomareOtagh,
		SyncTime:        syncTime,
	}, nil
}

func (m SanaModule) syncRoom(c fuego.ContextNoBody) (SanaRoomResponse, error) {
	idStr := c.PathParam("id")
	roomRackID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return SanaRoomResponse{}, fuego.BadRequestError{Title: "invalid_id"}
	}

	var roomRack models.SanaRoomRack
	if err := m.db.First(&roomRack, roomRackID).Error; err != nil {
		return SanaRoomResponse{}, err
	}

	var room models.Room
	if roomRack.RoomID != nil {
		if err := m.db.First(&room, *roomRack.RoomID).Error; err != nil {
			return SanaRoomResponse{}, err
		}
	} else {
		return SanaRoomResponse{}, fuego.BadRequestError{Title: "room_not_linked"}
	}

	client := sana.NewClient(sana.Config{
		KelidVahed:     m.sanaCfg.KelidVahed,
		KelidPeimankar: m.sanaCfg.KelidPeimankar,
		CodeVahed:      m.sanaCfg.CodeVahed,
	})

	// Build full hotel RAC from all rooms
	var allRooms []models.Room
	if err := m.db.Preload("Floor").Find(&allRooms).Error; err != nil {
		return SanaRoomResponse{}, err
	}

	// Group rooms by floor
	floorMap := make(map[int][]string)
	for _, r := range allRooms {
		floorNum := 0
		if r.Floor != nil {
			floorNum = r.Floor.Number
		}
		floorMap[floorNum] = append(floorMap[floorNum], r.RoomNumber)
	}

	// Sort floors
	var floors []int
	for f := range floorMap {
		floors = append(floors, f)
	}
	sort.Ints(floors)

	var sanaFloors []sana.Floor
	for _, f := range floors {
		// Sort rooms within floor
		sort.Strings(floorMap[f])
		sanaFloors = append(sanaFloors, sana.Floor{
			Number: f,
			Rooms:  floorMap[f],
		})
	}

	input := sana.SabtChidemanInput{
		Sakhteman:  1,
		TedadOtagh: len(allRooms),
		Floors:     sanaFloors,
	}

	result, err := client.SabtChidemanVahed(input)
	if err != nil {
		roomRack.LastError = err.Error()
		m.db.Save(&roomRack)
		return SanaRoomResponse{}, fuego.BadRequestError{Title: "sana_api_error", Detail: err.Error()}
	}

	if !result.IsOK {
		roomRack.LastError = result.Message
		m.db.Save(&roomRack)
		return SanaRoomResponse{}, fuego.BadRequestError{Title: "sana_api_error", Detail: result.Message}
	}

	roomRack.IsSynced = true
	roomRack.LastError = ""
	roomRack.LastSyncTime = time.Now()
	m.db.Save(&roomRack)

	// Also update all other room racks for the same hotel to share sync status
	// since SabtChidemanVahed is for the whole hotel
	m.db.Model(&models.SanaRoomRack{}).
		Where("hotel_id = ?", roomRack.HotelID).
		Updates(map[string]interface{}{
			"is_synced":      true,
			"last_sync_time": time.Now(),
			"last_error":     "",
		})

	lastSyncTime := ""
	if !roomRack.LastSyncTime.IsZero() {
		lastSyncTime = roomRack.LastSyncTime.Format("2006-01-02T15:04:05Z")
	}

	return SanaRoomResponse{
		ID:           roomRack.ID,
		HotelID:      roomRack.HotelID,
		Room:         &RoomInfo{ID: room.ID, RoomNumber: room.RoomNumber},
		Rac:          roomRack.Rac,
		LastSyncTime: lastSyncTime,
		IsSynced:     roomRack.IsSynced,
		LastError:    roomRack.LastError,
	}, nil
}

type SanaSyncAllResult struct {
	Status       string   `json:"status"`
	GuestsSynced int      `json:"guestsSynced"`
	GuestsFailed int      `json:"guestsFailed"`
	RoomsSynced  int      `json:"roomsSynced"`
	RoomsFailed  int      `json:"roomsFailed"`
	Errors       []string `json:"errors,omitempty"`
}

func (m SanaModule) syncAll(c h.ContextNoBody) (SanaSyncAllResult, error) {
	client := sana.NewClient(sana.Config{
		KelidVahed:     m.sanaCfg.KelidVahed,
		KelidPeimankar: m.sanaCfg.KelidPeimankar,
		CodeVahed:      m.sanaCfg.CodeVahed,
	})

	var result SanaSyncAllResult

	// Sync all unsynced guests
	var unsyncedGuests []models.SanaGuest
	m.db.Where("sync_time IS NULL OR sync_time = ?", time.Time{}).Find(&unsyncedGuests)

	for _, sg := range unsyncedGuests {
		var guest models.Guest
		if err := m.db.First(&guest, sg.GuestID).Error; err != nil {
			result.GuestsFailed++
			continue
		}

		gender := 1
		if guest.Gender == "female" {
			gender = 2
		}

		var shomareOtagh string
		var activeStay models.Stay
		if err := m.db.Joins("JOIN stay_statuses ON stays.status_id = stay_statuses.id").
			Where("stays.guest_id = ? AND stay_statuses.slug = ?", guest.ID, string(models.StayStatusResident)).
			First(&activeStay).Error; err == nil {
			var room models.Room
			if err := m.db.First(&room, activeStay.RoomID).Error; err == nil {
				shomareOtagh = room.RoomNumber
			}
		}
		if shomareOtagh == "" {
			shomareOtagh = sg.ShomareOtagh
		}

		input := sana.SabtMosaferinInput{
			NameMosafer:      guest.FirstName,
			FamilMosafer:     guest.LastName,
			NamePedar:        guest.FatherName,
			ShomareShenasaee: guest.NationalID,
			TarikhTavalod:    guest.DateOfBirth.Format("2006/01/02"),
			ID_Jensiat:       gender,
			TedadHamrah:      0,
			MosafereKhareji:  false,
			ShomareOtagh:     shomareOtagh,
			NameKarbareSabt:  "system",
			ID_NoeDadeh:      sana.DataTypeRoomReception,
		}

		apiResult, err := client.SabtMosaferin(input)
		if err != nil {
			result.GuestsFailed++
			result.Errors = append(result.Errors, fmt.Sprintf("Guest %s %s: %s", guest.FirstName, guest.LastName, err.Error()))
			continue
		}
		if !apiResult.IsOK {
			result.GuestsFailed++
			detail := apiResult.RawMessage
			if apiResult.CodeError > 0 {
				detail = fmt.Sprintf("SANA CodeError=%d: %s", apiResult.CodeError, apiResult.RawMessage)
			}
			result.Errors = append(result.Errors, fmt.Sprintf("Guest %s %s: %s", guest.FirstName, guest.LastName, detail))
			continue
		}

		sg.RecordMosafer = apiResult.RecordMosafer
		sg.ShomarePaziresh = apiResult.ShomarePaziresh
		sg.ShomareOtagh = apiResult.ShomareOtagh
		sg.SyncTime = time.Now()
		m.db.Save(&sg)
		result.GuestsSynced++
	}

	// Sync room rack for the default hotel
	var allRooms []models.Room
	if err := m.db.Preload("Floor").Find(&allRooms).Error; err == nil {
		floorMap := make(map[int][]string)
		for _, r := range allRooms {
			floorNum := 0
			if r.Floor != nil {
				floorNum = r.Floor.Number
			}
			floorMap[floorNum] = append(floorMap[floorNum], r.RoomNumber)
		}
		var floors []int
		for f := range floorMap {
			floors = append(floors, f)
		}
		sort.Ints(floors)
		var sanaFloors []sana.Floor
		for _, f := range floors {
			sort.Strings(floorMap[f])
			sanaFloors = append(sanaFloors, sana.Floor{
				Number: f,
				Rooms:  floorMap[f],
			})
		}
		input := sana.SabtChidemanInput{
			Sakhteman:  1,
			TedadOtagh: len(allRooms),
			Floors:     sanaFloors,
		}
		apiResult, err := client.SabtChidemanVahed(input)
		if err != nil {
			result.RoomsFailed = 1
			result.Errors = append(result.Errors, fmt.Sprintf("Room rack: %s", err.Error()))
		} else if !apiResult.IsOK {
			result.RoomsFailed = 1
			result.Errors = append(result.Errors, fmt.Sprintf("Room rack: %s", apiResult.Message))
			m.db.Model(&models.SanaRoomRack{}).
				Where("hotel_id = ?", "default").
				Update("last_error", apiResult.Message)
		} else {
			result.RoomsSynced = 1
			m.db.Model(&models.SanaRoomRack{}).
				Where("hotel_id = ?", "default").
				Updates(map[string]interface{}{
					"is_synced":      true,
					"last_sync_time": time.Now(),
					"last_error":     "",
				})
		}
	}

	if len(result.Errors) > 0 {
		result.Status = "partial_failure"
	} else {
		result.Status = "completed"
	}
	return result, nil
}
