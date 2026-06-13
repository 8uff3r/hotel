package models

import "time"

type SanaGuest struct {
	Base
	GuestID         uint      `gorm:"not null;uniqueIndex" json:"guestId"`
	Guest           Guest     `gorm:"foreignKey:GuestID" json:"-"`
	RecordMosafer   int       `json:"recordMosafer"`
	ShomarePaziresh string    `json:"shomarePaziresh"`
	ShomareOtagh    string    `json:"shomareOtagh"`
	SyncTime        time.Time `json:"syncTime"`
}

type SanaRoomRack struct {
	Base
	HotelID      string    `gorm:"not null;index" json:"hotelId"`
	Hotel        Hotel     `gorm:"foreignKey:HotelID" json:"-"`
	RoomID       *uint     `gorm:"index" json:"roomId"`
	Room         Room      `gorm:"foreignKey:RoomID" json:"-"`
	Rac          string    `gorm:"type:text" json:"rac"`
	LastSyncTime time.Time `json:"lastSyncTime"`
	IsSynced     bool      `gorm:"default:false" json:"isSynced"`
	LastError    string    `json:"lastError,omitempty"`
}

type TravelReason struct {
	Base
	TranslateBase
	Slug     *string `gorm:"uniqueIndex" json:"slug"`
	SanaID   string  `gorm:"uniqueIndex" json:"sanaId"`
	SanaName string  `json:"sanaName"`
}

type FamilyRelationship struct {
	Base
	TranslateBase
	Slug     *string `gorm:"uniqueIndex" json:"slug"`
	SanaID   string  `gorm:"uniqueIndex" json:"sanaId"`
	SanaName string  `json:"sanaName"`
}

type Nationality struct {
	Base
	TranslateBase
	Slug     *string `gorm:"uniqueIndex" json:"slug"`
	SanaID   string  `gorm:"uniqueIndex" json:"sanaId"`
	SanaName string  `json:"sanaName"`
}

type SanaCity struct {
	Base
	TranslateBase
	Slug     *string `gorm:"uniqueIndex" json:"slug"`
	SanaID   string  `gorm:"uniqueIndex" json:"sanaId"`
	SanaName string  `json:"sanaName"`
	IsIran   bool    `gorm:"default:false" json:"isIran"`
}

type Occupation struct {
	Base
	TranslateBase
	Slug     *string `gorm:"uniqueIndex" json:"slug"`
	SanaID   string  `gorm:"uniqueIndex" json:"sanaId"`
	SanaName string  `json:"sanaName"`
}
