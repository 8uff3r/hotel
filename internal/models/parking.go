package models

import (
	"time"
)

type ParkingLotStatus struct {
	Base
	TranslateBase
}

type ParkingLot struct {
	Base
	Name        string  `gorm:"not null" json:"name"`
	Location    string  `json:"location"`
	TotalSpots  int     `gorm:"not null;default:0" json:"totalSpots"`
	HourlyRate  float64 `gorm:"not null;default:0" json:"hourlyRate"`
	DailyRate   float64 `gorm:"not null;default:0" json:"dailyRate"`
	Description string  `json:"description"`

	StatusID uint             `gorm:"not null" json:"statusId"`
	Status   ParkingLotStatus `gorm:"foreignKey:StatusID" json:"status,omitzero"`

	HotelID *string `gorm:"not null" json:"hotelId"`
	Hotel   Hotel   `gorm:"foreignKey:HotelID" json:"hotel"`
}

type ParkingSpotStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type ParkingSpotType struct {
	Base
	TranslateBase
}

type ParkingSpot struct {
	Base
	SpotNumber  string `gorm:"not null" json:"spotNumber"`
	Floor       string `json:"floor"`
	IsCovered   bool   `gorm:"not null;default:false" json:"isCovered"`
	Description string `json:"description"`

	LotID *uint      `gorm:"not null" json:"lotId"`
	Lot   ParkingLot `gorm:"foreignKey:LotID" json:"lot,omitzero"`

	SpotTypeID uint            `gorm:"not null" json:"spotTypeId"`
	SpotType   ParkingSpotType `gorm:"foreignKey:SpotTypeID" json:"spotType,omitzero"`

	StatusID uint              `gorm:"not null" json:"statusId"`
	Status   ParkingSpotStatus `gorm:"foreignKey:StatusID" json:"status,omitzero"`
}

type VehicleType struct {
	Base
	TranslateBase
}

type Vehicle struct {
	Base
	LicensePlate string `gorm:"not null" json:"licensePlate"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	IsRegistered bool   `gorm:"not null;default:true" json:"isRegistered"`
	Notes        string `json:"notes"`

	VehicleTypeID uint        `gorm:"not null" json:"vehicleType"`
	VehicleType   VehicleType `gorm:"not null;foreignKey:VehicleTypeID" json:"vehicle"`

	GuestID uint  `gorm:"not null" json:"guestId"`
	Guest   Guest `gorm:"foreignKey:GuestID" json:"guest,omitzero"`
}

type ParkingTransaction struct {
	Base
	LotID           *uint         `json:"lotId"`
	Lot             ParkingLot    `gorm:"foreignKey:LotID" json:"lot,omitzero"`
	SpotID          *uint         `json:"spotId"`
	Spot            ParkingSpot   `gorm:"foreignKey:SpotID" json:"spot,omitzero"`
	GuestID         *uint         `json:"guestId"`
	Guest           Guest         `gorm:"foreignKey:GuestID" json:"guest,omitzero"`
	ReservationID   *uint         `json:"reservationId"`
	Reservation     Reservation   `gorm:"foreignKey:ReservationID" json:"reservation,omitzero"`
	LicensePlate    string        `gorm:"not null" json:"licensePlate"`
	EntryTime       time.Time     `gorm:"not null" json:"entryTime"`
	ExitTime        *time.Time    `json:"exitTime"`
	HoursParked     *float64      `json:"hoursParked"`
	RateApplied     *float64      `json:"rateApplied"`
	AmountDue       float64       `gorm:"not null;default:0" json:"amountDue"`
	AmountPaid      float64       `gorm:"not null;default:0" json:"amountPaid"`
	Status          string        `gorm:"not null;default:active" json:"status"`
	PaymentStatus   string        `gorm:"not null;default:pending" json:"paymentStatus"`
	PaymentMethodID *uint         `json:"paymentMethodId"`
	PaymentMethod   PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"paymentMethod,omitzero" translate:"true"`
	CheckoutID      *uint          `json:"checkoutId"`
	Notes           string        `json:"notes"`
}

type ParkingStats struct {
	Lots           int64 `json:"lots"`
	Spots          int64 `json:"spots"`
	AvailableSpots int64 `json:"availableSpots"`
}
