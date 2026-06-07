package models

import (
	"time"
)

type ReservationStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type Reservation struct {
	Base
	HotelID *string `json:"hotelId"`
	GuestID uint    `gorm:"not null;index" json:"guestId"`
	Guest   Guest   `gorm:"foreignKey:GuestID" json:"guest"`
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
	Parking   bool `gorm:"not null;default:false" json:"parking"`
	FullBoard bool `gorm:"not null;default:false" json:"fullBoard"`
	HalfBoard bool `gorm:"not null;default:false" json:"halfBoard"`

	RoomPrice float64 `gorm:"not null;default:0" json:"roomPrice"`

	StatusID *uint             `json:"statusId"`
	Status   ReservationStatus `gorm:"foreignKey:StatusID" json:"status"`

	PaymentDeadline *time.Time `json:"paymentDeadline"` // reservation expires if not paid by this time
	Notes         string     `json:"notes"`

	Payment Payment `gorm:"foreignKey:ReservationID" json:"payment"`
}

type Payment struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	ReservationID uint    `gorm:"uniqueIndex;not null" json:"reservationId"`
	Referrer      string  `json:"referrer"`
	Amount        float64 `gorm:"not null" json:"amount"`

	StatusID uint          `json:"statusId"`
	Status   PaymentStatus `gorm:"foreignKey:StatusID" json:"status"`

	MethodID uint          `gorm:"not null" json:"methodId"`
	Method   PaymentMethod `gorm:"foreignKey:MethodID" json:"method,omitzero" translate:"true"`
}
