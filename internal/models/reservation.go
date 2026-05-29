package models

import (
	"time"
)

type Reservation struct {
	Base
	HotelID *string `json:"hotelId"`
	GuestID uint    `gorm:"not null;index" json:"guestId"`
	Rooms   []Room  `gorm:"many2many:reservation_rooms;" json:"rooms"`

	ReservationCode string `gorm:"index" json:"reservationCode"`

	EntryDate      time.Time `gorm:"not null" json:"entryDate"`
	DepartureDate  time.Time `json:"departureDate"`
	DurationOfStay int       `json:"durationOfStay"`
	NumberOfPeople int       `gorm:"not null;default:1" json:"numberOfPeople"`

	Origin          string `json:"origin"`
	Destination     string `json:"destination"`
	PurposeOfTravel string `json:"purposeOfTravel"`

	Breakfast bool `gorm:"not null;default:false" json:"breakfast"`
	Guide     bool `gorm:"not null;default:false" json:"guide"`

	RoomPrice float64 `gorm:"not null;default:0" json:"roomPrice"`

	UserCheckIn  string `json:"userCheckIn"`
	UserCheckOut string `json:"userCheckOut"`

	Notes string `json:"notes"`

	Payment Payment `gorm:"foreignKey:ReservationID" json:"payment"`
}

type Payment struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	ReservationID uint   `gorm:"uniqueIndex;not null" json:"reservationId"`
	IsCash        bool   `gorm:"not null;default:false" json:"isCash"`
	Agency        bool   `gorm:"not null;default:false" json:"agency"`
	Referrer      string `json:"referrer"`
	ContractType  string `json:"contractType"`
}
