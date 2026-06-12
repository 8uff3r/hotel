package models

import (
	"time"
)

// InvoiceItemType represents the type of line item on an invoice
type InvoiceItemType string

const (
	InvoiceItemTypeRoomCharge    InvoiceItemType = "room_charge"
	InvoiceItemTypeBreakfast     InvoiceItemType = "breakfast"
	InvoiceItemTypeHalfBoard     InvoiceItemType = "half_board"
	InvoiceItemTypeFullBoard     InvoiceItemType = "full_board"
	InvoiceItemTypeParking       InvoiceItemType = "parking"
	InvoiceItemTypeRoomService   InvoiceItemType = "room_service"
	InvoiceItemTypeOther         InvoiceItemType = "other"
)

// PaymentStatusInvoice represents the payment status of an invoice or invoice item
type PaymentStatusInvoice string

const (
	PaymentStatusUnpaid       PaymentStatusInvoice = "unpaid"
	PaymentStatusPartiallyPaid PaymentStatusInvoice = "partially_paid"
	PaymentStatusCleared      PaymentStatusInvoice = "cleared"
)

// Invoice represents the bill for a stay
type Invoice struct {
	Base
	StayID          uint               `gorm:"not null;uniqueIndex" json:"stayId"`
	Stay            Stay               `gorm:"foreignKey:StayID" json:"-"`
	HotelID         string             `gorm:"not null" json:"hotelId"`
	TotalAmount     float64            `gorm:"not null;default:0" json:"totalAmount"`
	PaidAmount      float64            `gorm:"not null;default:0" json:"paidAmount"`
	RemainingAmount float64            `gorm:"not null;default:0" json:"remainingAmount"`
	PaymentStatus   string             `gorm:"not null;default:'unpaid'" json:"paymentStatus"`
	PaymentMethodID *uint              `json:"paymentMethodId"`
	PaymentMethod   *PaymentMethod     `gorm:"foreignKey:PaymentMethodID" json:"paymentMethod,omitempty" translate:"true"`
	CheckoutID      *uint              `json:"checkoutId"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`

	Items []InvoiceItem `gorm:"foreignKey:InvoiceID" json:"items,omitempty"`
}

// InvoiceItem represents a single line item on an invoice
type InvoiceItem struct {
	Base
	InvoiceID       uint            `gorm:"not null;index" json:"invoiceId"`
	StayID          uint            `gorm:"not null" json:"stayId"`
	ItemType        string          `gorm:"not null" json:"itemType"` // ROOM_CHARGE, BREAKFAST, etc.
	ServiceID       *uint           `json:"serviceId"` // nullable, for ROOM_SERVICE items
	Service         *Service        `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	Quantity        int             `gorm:"not null;default:1" json:"quantity"`
	UnitPrice       float64         `gorm:"not null;default:0" json:"unitPrice"`
	TotalPrice      float64         `gorm:"not null;default:0" json:"totalPrice"`
	Description     string          `json:"description"`
	PaidAmount      float64         `gorm:"not null;default:0" json:"paidAmount"`
	RemainingAmount float64         `gorm:"not null;default:0" json:"remainingAmount"`
	PaymentStatus   string          `gorm:"not null;default:'unpaid'" json:"paymentStatus"`
	PaymentMethodID *uint           `json:"paymentMethodId"`
	PaymentMethod   *PaymentMethod  `gorm:"foreignKey:PaymentMethodID" json:"paymentMethod,omitempty" translate:"true"`
}

// Service represents a billable service catalog item per hotel
type Service struct {
	Base
	HotelID     string  `gorm:"not null;index" json:"hotelId"`
	Name        string  `gorm:"not null" json:"name"`
	Unit        string  `json:"unit"`
	BaseAmount  float64 `gorm:"not null;default:0" json:"baseAmount"`
	Description string  `json:"description"`
}
