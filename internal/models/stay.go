package models

import (
	"time"

	"gorm.io/gorm"
)

// StayStatus represents the status of a stay/reception record
type StayStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type StayStatusSlug string

const (
	StayStatusWaiting    StayStatusSlug = "waiting"
	StayStatusResident   StayStatusSlug = "resident"
	StayStatusCheckedOut StayStatusSlug = "checked_out"
	StayStatusCancelled  StayStatusSlug = "cancelled"
	StayStatusNoShow     StayStatusSlug = "no_show"
)

type StayType string

const (
	StayTypeNormal        StayType = "normal"
	StayTypePreviousNight StayType = "previous_night"
	StayTypeEarlyCheckIn  StayType = "early_checkin"
	StayTypeHalfDay       StayType = "half_day"
)

// Stay represents an actual hotel stay/reception (check-in)
type Stay struct {
	Base
	HotelID       string       `gorm:"not null;index" json:"hotelId"`
	GuestID       uint         `gorm:"not null;index" json:"guestId"`
	Guest         Guest        `gorm:"foreignKey:GuestID" json:"guest"`
	RoomID        uint         `gorm:"not null;index" json:"roomId"`
	Room          Room         `gorm:"foreignKey:RoomID" json:"room"`
	ReservationID *uint        `json:"reservationId"` // links to original reservation if applicable
	Reservation   *Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`

	AcceptanceID           string     `gorm:"uniqueIndex;not null" json:"acceptanceId"`
	StayType               string     `gorm:"not null" json:"stayType"`
	EntryDate              time.Time  `gorm:"not null" json:"entryDate"`
	DepartureDate          time.Time  `json:"departureDate"`
	ScheduledEntryDate     time.Time  `json:"scheduledEntryDate"`
	ScheduledDepartureDate time.Time  `json:"scheduledDepartureDate"`
	ActualCheckIn          *time.Time `json:"actualCheckIn"`
	ActualCheckOut         *time.Time `json:"actualCheckOut"`
	DurationOfStay         int        `json:"durationOfStay"`
	NumberOfPeople         int        `gorm:"not null;default:1" json:"numberOfPeople"`

	Origin          string  `json:"origin"`
	Destination     string  `json:"destination"`
	PurposeOfTravel string  `json:"purposeOfTravel"`
	RoomPrice       float64 `gorm:"not null;default:0" json:"roomPrice"`
	Breakfast       bool    `gorm:"not null;default:false" json:"breakfast"`
	HalfBoard       bool    `gorm:"not null;default:false" json:"halfBoard"`
	FullBoard       bool    `gorm:"not null;default:false" json:"fullBoard"`
	Parking         bool    `gorm:"not null;default:false" json:"parking"`
	Notes           string  `json:"notes"`

	StayCode string `gorm:"index" json:"stayCode"`

	TravelAgencyID *uint         `json:"travelAgencyId"`
	TravelAgency   *TravelAgency `gorm:"foreignKey:TravelAgencyID" json:"travelAgency,omitempty"`

	EarlyCheckInFee float64 `gorm:"not null;default:0" json:"earlyCheckInFee"`
	HalfDayFee      float64 `gorm:"not null;default:0" json:"halfDayFee"`

	StatusID uint       `gorm:"not null" json:"statusId"`
	Status   StayStatus `gorm:"foreignKey:StatusID" json:"status,omitzero" translate:"true"`

	Invoice *Invoice `gorm:"foreignKey:StayID" json:"invoice,omitempty"`
}

func (s *Stay) AfterSave(tx *gorm.DB) error {
	return UpdateGuestStatus(tx, s.GuestID)
}
