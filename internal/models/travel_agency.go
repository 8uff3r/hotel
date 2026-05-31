package models

type TravelAgency struct {
	Base
	Name   string `gorm:"not null" json:"name"`
	CEOFirstName string `json:"ceoFirstName"`
	CEOLastName  string `json:"ceoLastName"`
	Province     string `json:"province"`
	City         string `json:"city"`
	Status       string `gorm:"not null;default:enabled" json:"status"`
}
