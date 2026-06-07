package models

import "gorm.io/gorm"

type Hotel struct {
	ID               string  `gorm:"primaryKey" json:"id"`
	Code             string  `gorm:"uniqueIndex;not null" json:"code"`
	Name             string  `gorm:"not null" json:"name"`
	Address          string  `gorm:"not null" json:"address"`
	Phone            string  `json:"phone"`
	Email            string  `json:"email"`
	TotalCapacity    int     `json:"totalCapacity"`
	NumberOfFloors   int     `json:"numberOfFloors"`
	CEOName          string  `json:"ceoName"`
	NearbyFacilities string  `gorm:"type:text" json:"nearbyFacilities"` // JSON array stored as text
	DeletedAt        gorm.DeletedAt

	Pictures []HotelPicture `gorm:"foreignKey:HotelID" json:"pictures,omitempty"`
	Setting  *HotelSetting  `gorm:"foreignKey:HotelID" json:"setting,omitempty"`
}

type HotelPicture struct {
	Base
	HotelID     string `gorm:"not null;index" json:"hotelId"`
	URL         string `gorm:"not null" json:"url"`
	Description string `json:"description"`
}

type HotelSetting struct {
	Base
	HotelID               string `gorm:"not null;uniqueIndex" json:"hotelId"`
	StandardCheckInTime  string `gorm:"not null;default:'14:00'" json:"standardCheckInTime"`
	StandardCheckOutTime string `gorm:"not null;default:'12:00'" json:"standardCheckOutTime"`
	NightAuditHour       string `gorm:"not null;default:'03:00'" json:"nightAuditHour"`
}
