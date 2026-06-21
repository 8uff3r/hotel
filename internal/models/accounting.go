package models

import (
	"time"
)

type Account struct {
	Base
	HotelID        *string `json:"hotelId"`
	AccountCode    string  `gorm:"not null" json:"accountCode"`
	AccountName    string  `gorm:"not null" json:"accountName"`
	AccountType    string  `gorm:"not null" json:"accountType"`
	AccountSubType string  `json:"accountSubType"`
	ParentID       *uint   `json:"parentId"`
	IsActive       bool    `gorm:"not null;default:true" json:"isActive"`
	IsSystem       bool    `gorm:"not null;default:false" json:"isSystem"`
	Description    string  `json:"description"`
	NormalBalance  string  `gorm:"not null;default:debit" json:"normalBalance"`
}

type PaymentMethod struct {
	Base
	TranslateBase
}

type PaymentStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type ExpenseCategory struct {
	Base
	TranslateBase
}

type Expense struct {
	Base
	ExpenseDate time.Time `gorm:"not null" json:"expenseDate"`
	Description string    `gorm:"not null" json:"description"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Vendor      string    `json:"vendor"`
	Reference   string    `json:"reference"`
	Notes       string    `json:"notes"`
	CreatedBy   *uint     `json:"createdBy"`

	HotelID *string `gorm:"not null" json:"hotelId"`
	Hotel   Hotel   `gorm:"foreignKey:HotelID" json:"hotel"`

	AccountID *uint   `gorm:"not null" json:"accountId"`
	Account   Account `gorm:"foreignKey:AccountID" json:"account,omitzero"`

	CategoryID uint            `gorm:"not null" json:"categoryId"`
	Category   ExpenseCategory `gorm:"foreignKey:CategoryID" json:"category,omitzero" translate:"true"`

	PaymentStatusID uint          `gorm:"not null" json:"paymentStatusId"`
	PaymentStatus   PaymentStatus `gorm:"foreignKey:PaymentStatusID" json:"paymentStatus,omitzero"`

	PaymentMethodID uint          `gorm:"not null" json:"paymentMethodId"`
	PaymentMethod   PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"paymentMethod,omitzero" translate:"true"`
}

type IncomeCategory struct {
	Base
	TranslateBase
}

type Income struct {
	Base
	IncomeDate  time.Time `gorm:"not null" json:"incomeDate"`
	Description string    `gorm:"not null" json:"description"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Source      string    `json:"source"`
	Reference   string    `json:"reference"`
	Notes       string    `json:"notes"`
	CreatedBy   *uint     `json:"createdBy"`

	HotelID string `gorm:"not null" json:"hotelId"`
	Hotel   Hotel  `gorm:"foreignKey:HotelID" json:"hotel"`

	AccountID *uint   `gorm:"not null" json:"accountId"`
	Account   Account `gorm:"foreignKey:AccountID" json:"account,omitzero"`

	ReservationID *uint        `json:"reservationId"`
	Reservation   *Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`

	StayID *uint `json:"stayId"`
	Stay   *Stay `gorm:"foreignKey:StayID" json:"stay,omitempty"`

	InvoiceID *uint    `json:"invoiceId"`
	Invoice   *Invoice `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`

	PaymentStatusID uint          `gorm:"not null" json:"paymentStatusId"`
	PaymentStatus   PaymentStatus `gorm:"foreignKey:PaymentStatusID" json:"paymentStatus,omitzero"`

	PaymentMethodID uint          `gorm:"not null" json:"paymentMethodId"`
	PaymentMethod   PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"paymentMethod,omitzero" translate:"true"`

	CategoryID uint           `gorm:"not null" json:"categoryId"`
	Category   IncomeCategory `gorm:"foreignKey:CategoryID" json:"category,omitzero" translate:"true"`
}
