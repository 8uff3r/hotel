package models

import (
	"time"
)

type RestaurantBillStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type RestaurantBill struct {
	Base
	BillDate           time.Time  `gorm:"not null" json:"billDate"`
	Subtotal           float64    `gorm:"not null;default:0" json:"subtotal"`
	TaxAmount          float64    `gorm:"not null;default:0" json:"taxAmount"`
	DiscountAmount     float64    `gorm:"not null;default:0" json:"discountAmount"`
	TotalAmount        float64    `gorm:"not null;default:0" json:"totalAmount"`
	IsExternal         bool       `gorm:"not null;default:false" json:"isExternal"`
	ExternalRestaurant string     `json:"externalRestaurant"`
	Notes              string     `json:"notes"`
	Settled            bool       `gorm:"not null;default:false" json:"settled"`
	SettledAt          *time.Time `json:"settledAt"`
	SettledBy          *uint      `json:"settledBy"`

	HotelID *string `gorm:"not null" json:"hotelId"`
	Hotel   Hotel   `gorm:"foreignKey:HotelID" json:",omitzero"`

	ReservationID *uint       `gorm:"index" json:"reservationId"`
	Reservation   Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitzero"`

	GuestID *uint `gorm:"index" json:"guestId"`
	Guest   Guest `gorm:"foreignKey:GuestID" json:"guest,omitzero"`

	RoomID *uint `json:"roomId"`
	Room   Room  `gorm:"foreignKey:RoomID" json:"room,omitzero"`

	StatusID uint                 `gorm:"not null" json:"statusId"`
	Status   RestaurantBillStatus `gorm:"foreignKey:StatusID" json:"status,omitzero"`

	CheckoutID *uint `json:"checkoutId"`
}

type RestaurantStats struct {
	TotalBills      int64   `json:"totalBills"`
	TotalRevenue    float64 `json:"totalRevenue"`
	InternalRevenue float64 `json:"internalRevenue"`
	ExternalRevenue float64 `json:"externalRevenue"`
	TotalMeals      int64   `json:"totalMeals"`
}

type MealTransaction struct {
	Base
	ItemName   string  `gorm:"not null" json:"itemName"`
	Quantity   float64 `gorm:"not null" json:"quantity"`
	UnitPrice  float64 `gorm:"not null" json:"unitPrice"`
	TotalPrice float64 `gorm:"not null" json:"totalPrice"`
	IsExternal bool    `gorm:"not null;default:false" json:"isExternal"`
	Notes      string  `json:"notes"`

	InventoryItemID *uint         `json:"inventoryItemId"`
	InventoryItem   InventoryItem `gorm:"foreignKey:InventoryItemID" json:"inventoryItem"`

	BillID uint           `gorm:"not null;index" json:"billId"`
	Bill   RestaurantBill `gorm:"foreignKey:BillID" json:"-"`

	HotelID *string `gorm:"not null" json:"hotelId"`
	Hotel   Hotel   `gorm:"foreignKey:HotelID" json:"hotel"`
}
