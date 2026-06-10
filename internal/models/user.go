package models

import (
	"time"

	"gorm.io/gorm"
)

type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBanned   UserStatus = "banned"
)

type User struct {
	Base
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	FirstName    string         `gorm:"not null" json:"firstName"`
	LastName     string         `gorm:"not null" json:"lastName"`
	Username     string         `gorm:"uniqueIndex" json:"username"`
	ContactNumber string        `json:"contactNumber"`
	Role         string         `json:"role"`
	Status       string         `gorm:"not null;default:'active'" json:"status"`
	HotelID      string         `gorm:"not null;index" json:"hotelId"`
	Hotel        Hotel          `gorm:"foreignKey:HotelID" json:"hotel,omitempty"`
	UserHotels   []UserHotel    `gorm:"foreignKey:UserID" json:"userHotels"`
	IsActive     bool           `gorm:"not null;default:true" json:"isActive"`
	Roles        []UserTemplate `gorm:"foreignKey:UserID" json:"role,omitempty"`
	DeletedAt    gorm.DeletedAt

	Permissions []UserPermission `gorm:"foreignKey:UserID" json:"permissions,omitempty"`
}

type UserHotel struct {
	Base
	UserID  uint   `gorm:"not null;index" json:"userId"`
	User    User   `gorm:"foreignKey:UserID" json:"-"`
	HotelID string `gorm:"not null;index" json:"hotelId"`
	Hotel   Hotel  `gorm:"foreignKey:HotelID" json:"hotel"`
}

type Session struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	User      User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	AdminID   *uint     `gorm:"index" json:"adminId"`
	Admin     *Admin    `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	IsAdmin   bool      `gorm:"not null;default:false" json:"isAdmin"`
}

type SanitizedUser struct {
	ID            uint                 `json:"id"`
	Email         string               `json:"email"`
	FirstName     string               `json:"firstName"`
	LastName      string               `json:"lastName"`
	Username      string               `json:"username"`
	ContactNumber string               `json:"contactNumber"`
	Role          string               `json:"role"`
	Status        string               `json:"status"`
	HotelID       string               `json:"hotelId"`
	UserHotels    []UserHotelInfo      `json:"userHotels"`
	Roles         []PermissionTemplate `json:"roles,omitempty"`
	IsAdmin       bool                 `json:"isAdmin"`
	AdminHotels   []AdminHotelInfo     `json:"adminHotels,omitempty"`
}

type UserHotelInfo struct {
	HotelID string `json:"hotelId"`
	Hotel   Hotel  `json:"hotel"`
}

type SanitizedUserWithPermissions struct {
	ID          uint                 `json:"id"`
	Email       string               `json:"email"`
	FirstName   string               `json:"firstName"`
	LastName    string               `json:"lastName"`
	UserHotels  []UserHotelInfo      `json:"userHotels"`
	Permissions []UserPermissionInfo `json:"permissions"`
}

type UserPermissionInfo struct {
	Label        string             `json:"label"`
	PermissionID uint               `json:"permissionId"`
	Page         string             `json:"page"`
	Action       PermissionAction   `json:"action"`
	Category     PermissionCategory `json:"category"`
	Granted      bool               `json:"granted"`
}

type UserTemplateInfo struct {
	TranslateBase
	TemplateID  uint   `json:"templateId"`
	Description string `json:"description"`
}
