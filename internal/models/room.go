package models

type Room struct {
	Base
	HotelID     *uint     `json:"hotelId,omitempty"`
	Name        string    `json:"name"`
	RoomNumber  string    `gorm:"not null" json:"roomNumber" validate:"required,min=1"`
	Floor       int       `json:"floor"`
	Capacity    int       `gorm:"not null;default:2" json:"capacity"`
	BasePrice   float64   `gorm:"not null;default:0" json:"basePrice"`
	Amenities   []Amenity `gorm:"many2many:room_amenities;" json:"amenities" translate:"true"`
	Description string    `json:"description"`

	TypeID uint     `gorm:"not null" json:"roomTypeId"`
	Type   RoomType `gorm:"foreignKey:TypeID" json:"roomType,omitzero" translate:"true"`

	StatusID uint       `gorm:"not null" json:"statusId"`
	Status   RoomStatus `gorm:"foreignKey:StatusID" json:"status,omitzero" translate:"true"`
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
