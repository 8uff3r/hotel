package sana

import (
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

	h.Post(s, "/guests/:id/sync", m.syncGuest)
	h.Post(s, "/rooms/:id/sync", m.syncRoom)

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

func (m SanaModule) getSanaGuests(c h.ContextNoBody) ([]SanaGuestResponse, error) {
	var guests []models.SanaGuest
	err := m.db.Preload("Guest").Find(&guests).Error
	if err != nil {
		return nil, err
	}

	var response []SanaGuestResponse
	for _, g := range guests {
		gi := &GuestInfo{}
		if g.GuestID > 0 {
			gi.FirstName = g.Guest.FirstName
			gi.LastName = g.Guest.LastName
			gi.NationalID = g.Guest.NationalID
		}
		syncTime := ""
		if !g.SyncTime.IsZero() {
			syncTime = g.SyncTime.Format("2006-01-02T15:04:05Z")
		}
		response = append(response, SanaGuestResponse{
			ID:              g.ID,
			GuestID:         g.GuestID,
			Guest:           gi,
			RecordMosafer:   g.RecordMosafer,
			ShomarePaziresh: g.ShomarePaziresh,
			ShomareOtagh:    g.ShomareOtagh,
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
	var roomRacks []models.SanaRoomRack
	err := m.db.Find(&roomRacks).Error
	if err != nil {
		return nil, err
	}

	var response []SanaRoomResponse
	for _, r := range roomRacks {
		ri := &RoomInfo{}
		lastSyncTime := ""
		if !r.LastSyncTime.IsZero() {
			lastSyncTime = r.LastSyncTime.Format("2006-01-02T15:04:05Z")
		}
		response = append(response, SanaRoomResponse{
			ID:           r.ID,
			HotelID:      r.HotelID,
			Room:         ri,
			Rac:          r.Rac,
			LastSyncTime: lastSyncTime,
			IsSynced:     r.IsSynced,
			LastError:    r.LastError,
		})
	}
	return response, nil
}

func (m SanaModule) syncGuest(c fuego.ContextNoBody) (SanaGuestResponse, error) {
	id := c.PathParam("id")
	var guest models.SanaGuest
	if err := m.db.First(&guest, id).Error; err != nil {
		return SanaGuestResponse{}, err
	}

	client := sana.NewClient(sana.Config{
		KelidVahed:     m.sanaCfg.KelidVahed,
		KelidPeimankar: m.sanaCfg.KelidPeimankar,
		CodeVahed:      m.sanaCfg.CodeVahed,
	})

	var g models.Guest
	if err := m.db.First(&g, guest.GuestID).Error; err != nil {
		return SanaGuestResponse{}, err
	}

	gender := 1
	if g.Gender == "female" {
		gender = 2
	}

	input := sana.SabtMosaferinInput{
		NameMosafer:      g.FirstName,
		FamilMosafer:     g.LastName,
		NamePedar:        g.FatherName,
		ShomareShenasaee: g.NationalID,
		TarikhTavalod:    g.DateOfBirth.Format("2006/01/02"),
		ID_Jensiat:       gender,
		TedadHamrah:      0,
		MosafereKhareji:  false,
		ShomareOtagh:     guest.ShomareOtagh,
		NameKarbareSabt:  "system",
		ID_NoeDadeh:      sana.DataTypeRoomReception,
	}

	result, err := client.Sabt_Mosaferin(input)
	if err != nil {
		return SanaGuestResponse{}, err
	}

	guest.RecordMosafer = result.RecordMosafer
	guest.ShomarePaziresh = result.ShomarePaziresh
	guest.ShomareOtagh = result.ShomareOtagh

	m.db.Save(&guest)

	syncTime := ""
	if !guest.SyncTime.IsZero() {
		syncTime = guest.SyncTime.Format("2006-01-02T15:04:05Z")
	}

	return SanaGuestResponse{
		ID:              guest.ID,
		GuestID:         guest.GuestID,
		RecordMosafer:   guest.RecordMosafer,
		ShomarePaziresh: guest.ShomarePaziresh,
		ShomareOtagh:    guest.ShomareOtagh,
		SyncTime:        syncTime,
	}, nil
}

func (m SanaModule) syncRoom(c fuego.ContextNoBody) (SanaRoomResponse, error) {
	id := c.PathParam("id")
	var roomRack models.SanaRoomRack
	if err := m.db.First(&roomRack, id).Error; err != nil {
		return SanaRoomResponse{}, err
	}

	client := sana.NewClient(sana.Config{
		KelidVahed:     m.sanaCfg.KelidVahed,
		KelidPeimankar: m.sanaCfg.KelidPeimankar,
		CodeVahed:      m.sanaCfg.CodeVahed,
	})

	input := sana.SabtChidemanInput{
		Sakhteman:  1,
		TedadOtagh: 0,
		Floors:     []sana.Floor{},
	}

	result, err := client.SabtChidemanVahed(input)
	if err != nil {
		roomRack.LastError = err.Error()
		m.db.Save(&roomRack)
		return SanaRoomResponse{}, err
	}

	roomRack.IsSynced = result.IsOK
	roomRack.LastError = ""
	m.db.Save(&roomRack)

	lastSyncTime := ""
	if !roomRack.LastSyncTime.IsZero() {
		lastSyncTime = roomRack.LastSyncTime.Format("2006-01-02T15:04:05Z")
	}

	return SanaRoomResponse{
		ID:           roomRack.ID,
		HotelID:      roomRack.HotelID,
		Rac:          roomRack.Rac,
		LastSyncTime: lastSyncTime,
		IsSynced:     roomRack.IsSynced,
		LastError:    roomRack.LastError,
	}, nil
}

func (m SanaModule) syncAll(c h.ContextNoBody) (map[string]string, error) {
	result := map[string]string{
		"status": "completed",
	}
	return result, nil
}
