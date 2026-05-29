package models

type Hotel struct {
	ID      string `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Address string `gorm:"not null" json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}
