package models

type InventoryItemUnit struct {
	Base
	TranslateBase
}

type InventoryItemStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type InventoryItemCategory struct {
	Base
	TranslateBase
}

type InventoryItem struct {
	Base
	Name         string  `gorm:"not null" json:"name"`
	Quantity     float64 `gorm:"not null;default:0" json:"quantity"`
	UnitCost     float64 `gorm:"not null;default:0" json:"unitCost"`
	ReorderLevel float64 `gorm:"not null;default:0" json:"reorderLevel"`
	IsActive     bool    `gorm:"not null;default:true" json:"isActive"`
	Description  string  `json:"description"`

	UnitID uint              `gorm:"not null" json:"unitId"`
	Unit   InventoryItemUnit `gorm:"foreignKey:UnitID" json:"unit,omitzero"`

	CategoryID uint                  `gorm:"not null" json:"categoryId"`
	Category   InventoryItemCategory `gorm:"foreignKey:CategoryID" json:"category,omitzero"`

	StatusID uint                `gorm:"not null" json:"statusId"`
	Status   InventoryItemStatus `gorm:"foreignKey:StatusID" json:"status,omitzero"`

	HotelID *uint `gorm:"not null" json:"hotelId"`
	Hotel   Hotel `gorm:"foreignKey:HotelID" json:",omitzero"`
}
