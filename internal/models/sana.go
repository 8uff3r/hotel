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
	Rac          string    `gorm:"type:text" json:"rac"`
	LastSyncTime time.Time `json:"lastSyncTime"`
	IsSynced     bool      `gorm:"default:false" json:"isSynced"`
	LastError    string    `json:"lastError,omitempty"`
}
