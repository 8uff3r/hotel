package models

import (
	"time"

	"gorm.io/gorm"
)

type GuestStatus string

const (
	GuestStatusWaiting    GuestStatus = "waiting"
	GuestStatusResident   GuestStatus = "resident"
	GuestStatusCheckedOut GuestStatus = "checked_out"
	GuestStatusCancelled  GuestStatus = "cancelled"
	GuestStatusAbsence    GuestStatus = "absence"
)

type Guest struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FirstName    string    `gorm:"not null" json:"firstName" validate:"required,min=2,max=50"`
	LastName     string    `gorm:"not null" json:"lastName" validate:"required,min=2,max=50"`
	FatherName   string    `gorm:"not null" json:"fatherName"`
	NationalID   string    `gorm:"index" json:"nationalId"`
	IDNumber     string    `json:"idNumber"`
	Passport     string    `json:"passport"`
	Gender       string    `json:"gender"`
	DateOfBirth  time.Time `json:"dateOfBirth"`
	PlaceOfBirth string    `json:"placeOfBirth"`
	Phone        string    `json:"phone"`
	Telephone    string    `json:"telephone"`
	Address      string    `json:"address"`
	PostalCode   string    `json:"postalCode"`
	Occupation   string    `json:"occupation"`
	Job          string    `json:"job"`
	Email        string    `json:"email"`
	Landline     string    `json:"landline"`
	Status       string    `gorm:"not null;default:'waiting'" json:"status"`

	HotelID string `json:"hotelId"`
	Hotel   Hotel  `gorm:"foreignKey:HotelID"`

	NationalityID uint    `json:"nationalityID" validate:"required"`
	Nationality   Country `gorm:"foreignKey:NationalityID" json:"nationality,omitzero"`

	Companions []GuestCompanion `gorm:"foreignKey:GuestID" json:"companions,omitempty"`

	DeletedAt gorm.DeletedAt
}

type GuestCompanion struct {
	Base
	GuestID     uint               `gorm:"not null;index" json:"guestId"`
	FirstName   string             `gorm:"not null" json:"firstName"`
	LastName    string             `gorm:"not null" json:"lastName"`
	FatherName  string             `json:"fatherName"`
	NationalID  string             `json:"nationalId"`
	IDNumber    string             `json:"idNumber"`
	Gender      string             `json:"gender"`
	RelationID  uint               `gorm:"not null" json:"relationId"`
	Relation    FamilyRelationship `gorm:"foreignKey:RelationID" json:"relation,omitzero"`
	DateOfBirth time.Time          `json:"dateOfBirth"`
	Phone       string             `json:"phone"`
	Mobile      string             `json:"mobile"`

	NationalityID uint    `json:"nationalityID"`
	Nationality   Country `gorm:"foreignKey:NationalityID" json:"nationality,omitzero"`
}
