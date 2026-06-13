package models

import (
	"time"

	"gorm.io/gorm"
)

type GuestStatusSlug string

const (
	GuestStatusWaiting    GuestStatusSlug = "waiting"
	GuestStatusResident   GuestStatusSlug = "resident"
	GuestStatusCheckedOut GuestStatusSlug = "checked_out"
	GuestStatusCancelled  GuestStatusSlug = "cancelled"
)

type GuestStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type Guest struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FirstName    string    `gorm:"not null" json:"firstName" validate:"required,min=2,max=50"`
	LastName     string    `gorm:"not null" json:"lastName" validate:"required,min=2,max=50"`
	FatherName   string    `gorm:"not null" json:"fatherName"`
	IDNumber     string    `json:"idNumber"`
	Gender       string    `json:"gender"`
	DateOfBirth  time.Time `json:"dateOfBirth" validate:"required"`
	PlaceOfBirth string    `json:"placeOfBirth"`
	Phone        string    `json:"phone"`
	Telephone    string    `json:"telephone"`
	Address      string    `json:"address"`
	PostalCode   string    `json:"postalCode"`
	Occupation   string    `json:"occupation"`
	Job          string    `json:"job"`
	Email        string    `json:"email"`
	Landline     string    `json:"landline"`

	StatusID *uint       `json:"statusId"`
	Status   GuestStatus `gorm:"foreignKey:StatusID" json:"status" translate:"true"`

	HotelID string `json:"hotelId"`
	Hotel   Hotel  `gorm:"foreignKey:HotelID"`

	NationalityID uint    `json:"nationalityID" validate:"required"`
	Nationality   Country `gorm:"foreignKey:NationalityID" json:"nationality,omitzero"`

	Companions []GuestCompanion `gorm:"foreignKey:GuestID" json:"companions,omitempty"`

	DeletedAt gorm.DeletedAt
}

func (g *Guest) BeforeSave(tx *gorm.DB) error {
	if g.Status.Slug != "" {
		var status GuestStatus
		if err := tx.Where("slug = ?", g.Status.Slug).First(&status).Error; err != nil {
			return err
		}
		g.StatusID = &status.ID
	}
	return nil
}

// UpdateGuestStatus updates a guest's status based on their current stays and reservations.
// Priority: resident (active stay) > waiting (active reservation) > checked_out (completed stay) > cancelled
func UpdateGuestStatus(db *gorm.DB, guestID uint) error {
	var guest Guest
	if err := db.First(&guest, guestID).Error; err != nil {
		return err
	}

	// 1. Check for active (resident) stays
	var activeStayCount int64
	db.Model(&Stay{}).Joins("JOIN stay_statuses ON stays.status_id = stay_statuses.id").
		Where("stays.guest_id = ? AND stay_statuses.slug = ?", guestID, string(StayStatusResident)).
		Count(&activeStayCount)
	if activeStayCount > 0 {
		return setGuestStatus(db, &guest, string(GuestStatusResident))
	}

	// 2. Check for active reservations (waiting state)
	var activeReservationCount int64
	db.Model(&Reservation{}).Joins("JOIN reservation_statuses ON reservations.status_id = reservation_statuses.id").
		Where("reservations.guest_id = ? AND reservation_statuses.slug IN ?", guestID, []string{"accepted", "verified", "awaiting_payment"}).
		Count(&activeReservationCount)
	if activeReservationCount > 0 {
		return setGuestStatus(db, &guest, string(GuestStatusWaiting))
	}

	// 3. Check for checked_out stays
	var checkedOutStayCount int64
	db.Model(&Stay{}).Joins("JOIN stay_statuses ON stays.status_id = stay_statuses.id").
		Where("stays.guest_id = ? AND stay_statuses.slug = ?", guestID, "checked_out").
		Count(&checkedOutStayCount)
	if checkedOutStayCount > 0 {
		return setGuestStatus(db, &guest, string(GuestStatusCheckedOut))
	}

	// 4. Check for cancelled states (stays or reservations)
	var cancelledStayCount int64
	db.Model(&Stay{}).Joins("JOIN stay_statuses ON stays.status_id = stay_statuses.id").
		Where("stays.guest_id = ? AND stay_statuses.slug IN ?", guestID, []string{"cancelled", "no_show"}).
		Count(&cancelledStayCount)

	var cancelledReservationCount int64
	db.Model(&Reservation{}).Joins("JOIN reservation_statuses ON reservations.status_id = reservation_statuses.id").
		Where("reservations.guest_id = ? AND reservation_statuses.slug IN ?", guestID, []string{"cancelled", "expired"}).
		Count(&cancelledReservationCount)

	if cancelledStayCount > 0 || cancelledReservationCount > 0 {
		return setGuestStatus(db, &guest, string(GuestStatusCancelled))
	}

	// 5. Default to waiting
	return setGuestStatus(db, &guest, string(GuestStatusWaiting))
}

func setGuestStatus(db *gorm.DB, guest *Guest, slug string) error {
	var status GuestStatus
	if err := db.Where("slug = ?", slug).First(&status).Error; err != nil {
		return err
	}
	return db.Model(guest).Update("status_id", &status.ID).Error
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
	DateOfBirth time.Time          `json:"dateOfBirth" validate:"required"`
	Phone       string             `json:"phone"`
	Mobile      string             `json:"mobile"`

	NationalityID uint    `json:"nationalityID"`
	Nationality   Country `gorm:"foreignKey:NationalityID" json:"nationality,omitzero"`
}
