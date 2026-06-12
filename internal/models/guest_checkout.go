package models

import (
	"time"
)

// GuestCheckout represents a unified checkout record for a guest
// aggregating stay invoice, parking, and restaurant charges.
type GuestCheckout struct {
	Base
	GuestID         uint       `gorm:"not null;index" json:"guestId"`
	Guest           Guest      `gorm:"foreignKey:GuestID" json:"guest,omitempty"`
	StayID          uint       `gorm:"not null;index" json:"stayId"`
	Stay            Stay       `gorm:"foreignKey:StayID" json:"stay,omitempty"`
	TotalRoom       float64    `gorm:"not null;default:0" json:"totalRoom"`
	TotalParking    float64    `gorm:"not null;default:0" json:"totalParking"`
	TotalRestaurant float64    `gorm:"not null;default:0" json:"totalRestaurant"`
	TotalPaid       float64    `gorm:"not null;default:0" json:"totalPaid"`
	PaymentMethodID *uint      `json:"paymentMethodId"`
	PaymentMethod   *PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"paymentMethod,omitempty" translate:"true"`
	ProcessedByID   *uint      `json:"processedById"`
	CheckedOutAt    time.Time  `json:"checkedOutAt"`
}
