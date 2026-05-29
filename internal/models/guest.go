package models

import (
	"time"
)

type Guest struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	HotelID      *string   `json:"hotelId"`
	FirstName    string    `gorm:"not null" json:"firstName" validate:"required,min=2,max=50"`
	LastName     string    `gorm:"not null" json:"lastName"`
	FatherName   string    `json:"fatherName"`
	NationalID   string    `gorm:"index" json:"nationalId"`
	IDNumber     string    `json:"idNumber"`
	Nationality  string    `json:"nationality"`
	Gender       string    `json:"gender"`
	DateOfBirth  time.Time `json:"dateOfBirth"`
	PlaceOfBirth string    `json:"placeOfBirth"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	PostalCode   string    `json:"postalCode"`
	Occupation   string    `json:"occupation"`
	Email        string    `json:"email"`
	Landline     string    `json:"landline"`

	Reservations []Reservation    `gorm:"foreignKey:GuestID" json:"reservations,omitempty"`
	Companions   []GuestCompanion `gorm:"foreignKey:GuestID" json:"companions,omitempty"`
}

type GuestCompanion struct {
	Base
	GuestID    uint   `gorm:"not null;index" json:"guestId"`
	FirstName  string `gorm:"not null" json:"firstName"`
	LastName   string `gorm:"not null" json:"lastName"`
	NationalID string `json:"nationalId"`
	IDNumber   string `json:"idNumber"`
	Relation   string `json:"relation"`
}
