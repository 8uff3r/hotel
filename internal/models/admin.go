package models

import "gorm.io/gorm"

// Admin represents a system administrator who can manage multiple hotels or all hotels
type Admin struct {
	Base
	FirstName     string    `gorm:"not null" json:"firstName"`
	LastName      string    `gorm:"not null" json:"lastName"`
	ContactNumber string    `json:"contactNumber"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	Username      string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash  string    `gorm:"not null" json:"-"`
	Role          string    `json:"role"`
	IsActive      bool      `gorm:"not null;default:true" json:"isActive"`
	IsSuperAdmin  bool      `gorm:"not null;default:false" json:"isSuperAdmin"`
	DeletedAt     gorm.DeletedAt

	AdminHotels []AdminHotel `gorm:"foreignKey:AdminID" json:"adminHotels,omitempty"`
}

type AdminHotel struct {
	Base
	AdminID uint  `gorm:"not null;index" json:"adminId"`
	Admin   Admin `gorm:"foreignKey:AdminID" json:"-"`
	HotelID string `gorm:"not null;index" json:"hotelId"`
	Hotel   Hotel `gorm:"foreignKey:HotelID" json:"hotel"`
}

type SanitizedAdmin struct {
	ID            uint             `json:"id"`
	FirstName     string           `json:"firstName"`
	LastName      string           `json:"lastName"`
	ContactNumber string           `json:"contactNumber"`
	Email         string           `json:"email"`
	Username      string           `json:"username"`
	Role          string           `json:"role"`
	IsActive      bool             `json:"isActive"`
	IsSuperAdmin  bool             `json:"isSuperAdmin"`
	AdminHotels   []AdminHotelInfo `json:"adminHotels"`
}

type AdminHotelInfo struct {
	HotelID string `json:"hotelId"`
	Hotel   Hotel  `json:"hotel"`
}

func SanitizeAdmin(a *Admin) SanitizedAdmin {
	var hotels []AdminHotelInfo
	for _, ah := range a.AdminHotels {
		hotels = append(hotels, AdminHotelInfo{
			HotelID: ah.HotelID,
			Hotel:   ah.Hotel,
		})
	}
	return SanitizedAdmin{
		ID:            a.ID,
		FirstName:     a.FirstName,
		LastName:      a.LastName,
		ContactNumber: a.ContactNumber,
		Email:         a.Email,
		Username:      a.Username,
		Role:          a.Role,
		IsActive:      a.IsActive,
		IsSuperAdmin:  a.IsSuperAdmin,
		AdminHotels:   hotels,
	}
}
