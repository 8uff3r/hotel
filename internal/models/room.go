package models

import "gorm.io/gorm"

type Room struct {
	Base
	HotelID     *string   `json:"hotelId,omitempty"`
	Name        string    `json:"name"`
	RoomNumber  string    `gorm:"not null;uniqueIndex:idx_room_number_hotel" json:"roomNumber"`
	FloorID     *uint     `json:"floorId"`
	Floor       *Floor    `gorm:"foreignKey:FloorID" json:"floor,omitzero"`
	Capacity    int       `gorm:"not null;default:2" json:"capacity"`
	BasePrice   float64   `gorm:"not null;default:0" json:"basePrice"`
	Amenities   []Amenity `gorm:"many2many:room_amenities;" json:"amenities" translate:"true"`
	Facilities  string    `json:"facilities"` // comma-separated or JSON list of facility names
	Description string    `json:"description"`

	TypeID uint     `gorm:"not null" json:"roomTypeId"`
	Type   RoomType `gorm:"foreignKey:TypeID" json:"roomType,omitzero" translate:"true"`

	StatusID uint       `gorm:"not null" json:"statusId"`
	Status   RoomStatus `gorm:"foreignKey:StatusID" json:"status,omitzero" translate:"true"`

	Pictures []RoomPicture `gorm:"foreignKey:RoomID" json:"pictures,omitempty"`
}

type RoomType struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type RoomStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type Amenity struct {
	Base
	TranslateBase
}

type Floor struct {
	Base
	HotelID     string `gorm:"not null;index" json:"hotelId"`
	Number      int    `gorm:"not null" json:"number"`
	Description string `json:"description"`
}

type RoomStatusSlug string

const (
	RoomStatusAvailable   RoomStatusSlug = "available"
	RoomStatusOccupied    RoomStatusSlug = "occupied"
	RoomStatusReserved    RoomStatusSlug = "reserved"
	RoomStatusUnderRepair RoomStatusSlug = "under_repair"
	RoomStatusCleaning    RoomStatusSlug = "cleaning"
)

type RoomPicture struct {
	Base
	RoomID      uint   `gorm:"not null;index" json:"roomId"`
	URL         string `gorm:"not null" json:"url"`
	Description string `json:"description"`
}

func (r *Room) BeforeSave(tx *gorm.DB) (err error) {
	if r.Status.Slug != "" && r.StatusID == 0 {
		var status RoomStatus
		if err := tx.Where("slug = ?", r.Status.Slug).First(&status).Error; err != nil {
			return err
		}
		r.StatusID = status.ID
		r.Status = status
	}
	return nil
}
