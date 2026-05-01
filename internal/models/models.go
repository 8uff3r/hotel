package models

//go:generate go run ./generate.go && go fmt ./generate.go

import (
	"time"
)

type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBanned   UserStatus = "banned"
)

type Base struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type User struct {
	Base
	Email        string      `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string      `gorm:"not null" json:"-"`
	FirstName    string      `gorm:"not null" json:"firstName"`
	LastName     string      `gorm:"not null" json:"lastName"`
	UserHotels   []UserHotel `gorm:"foreignKey:UserID" json:"userHotels"`
	IsActive     bool        `gorm:"not null;default:true" json:"isActive"`

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
	ExpiresAt time.Time `gorm:"not null;index" json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type Hotel struct {
	ID      string `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Address string `gorm:"not null" json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type Room struct {
	Base
	HotelID     *uint     `json:"hotelId,omitempty"`
	Name        string    `json:"name"`
	RoomNumber  string    `gorm:"not null" json:"roomNumber" validate:"required,min=1"`
	Floor       int       `json:"floor"`
	Capacity    int       `gorm:"not null;default:2" json:"capacity"`
	BasePrice   float64   `gorm:"not null;default:0" json:"basePrice"`
	Amenities   []Amenity `gorm:"many2many:room_amenities;" json:"amenities" translate:"true"`
	Description string    `json:"description"`

	TypeID uint     `gorm:"not null" json:"roomTypeId"`
	Type   RoomType `gorm:"foreignKey:TypeID" json:"roomType,omitzero" translate:"true"`

	StatusID uint       `gorm:"not null" json:"statusId"`
	Status   RoomStatus `gorm:"foreignKey:StatusID" json:"status,omitzero" translate:"true"`
}

type RoomType struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type RoomStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}

type Amenity struct {
	Base
	TranslateBase
}

type ParkingSpotType struct {
	Base
	TranslateBase
}

type ParkingSpotStatus struct {
	Base
	TranslateBase
}

type Guest struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	HotelID      *uint     `json:"hotelId"`
	FirstName    string    `gorm:"not null" json:"firstName" validate:"required,min=2,max=50"`
	LastName     string    `gorm:"not null" json:"lastName"`
	FatherName   string    `json:"fatherName"`
	NationalID   string    `gorm:"index" json:"nationalId"`
	IDNumber     string    `json:"idNumber"`
	Nationality  string    `json:"nationality"`
	Gender       string    `json:"gender"`
	DateOfBirth  time.Time `json:"dateOfBirth"`
	PlaceOfBirth string    `json:"placeOfBirth"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	PostalCode   string    `json:"postalCode"`
	Occupation   string    `json:"occupation"`

	Reservations []Reservation    `gorm:"foreignKey:GuestID" json:"reservations,omitempty"`
	Companions   []GuestCompanion `gorm:"foreignKey:GuestID" json:"companions,omitempty"`
}

type GuestCompanion struct {
	Base
	GuestID    uint   `gorm:"not null;index" json:"guestId"`
	FirstName  string `gorm:"not null" json:"firstName"`
	LastName   string `gorm:"not null" json:"lastName"`
	NationalID string `json:"nationalId"`
	IDNumber   string `json:"idNumber"`
	Relation   string `json:"relation"`
}

type Reservation struct {
	Base
	HotelID *string `json:"hotelId"`
	GuestID uint    `gorm:"not null;index" json:"guestId"`
	Rooms   []Room  `gorm:"many2many:reservation_rooms;" json:"rooms"`

	ReservationCode string `gorm:"index" json:"reservationCode"`

	EntryDate      time.Time `gorm:"not null" json:"entryDate"`
	DepartureDate  time.Time `json:"departureDate"`
	DurationOfStay int       `json:"durationOfStay"`
	NumberOfPeople int       `gorm:"not null;default:1" json:"numberOfPeople"`

	Origin          string `json:"origin"`
	Destination     string `json:"destination"`
	PurposeOfTravel string `json:"purposeOfTravel"`

	Breakfast bool `gorm:"not null;default:false" json:"breakfast"`
	Guide     bool `gorm:"not null;default:false" json:"guide"`

	RoomPrice float64 `gorm:"not null;default:0" json:"roomPrice"`

	UserCheckIn  string `json:"userCheckIn"`
	UserCheckOut string `json:"userCheckOut"`

	Notes string `json:"notes"`

	Payment Payment `gorm:"foreignKey:ReservationID" json:"payment"`
}

type Payment struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	ReservationID uint   `gorm:"uniqueIndex;not null" json:"reservationId"`
	IsCash        bool   `gorm:"not null;default:false" json:"isCash"`
	Agency        bool   `gorm:"not null;default:false" json:"agency"`
	Referrer      string `json:"referrer"`
	ContractType  string `json:"contractType"`
}

type Account struct {
	Base
	HotelID        *uint  `json:"hotelId"`
	AccountCode    string `gorm:"not null" json:"accountCode"`
	AccountName    string `gorm:"not null" json:"accountName"`
	AccountType    string `gorm:"not null" json:"accountType"`
	AccountSubType string `json:"accountSubType"`
	ParentID       *uint  `json:"parentId"`
	IsActive       bool   `gorm:"not null;default:true" json:"isActive"`
	IsSystem       bool   `gorm:"not null;default:false" json:"isSystem"`
	Description    string `json:"description"`
	NormalBalance  string `gorm:"not null;default:debit" json:"normalBalance"`
}

type Expense struct {
	Base
	HotelID       *uint     `json:"hotelId"`
	ExpenseDate   time.Time `gorm:"not null" json:"expenseDate"`
	Description   string    `gorm:"not null" json:"description"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Category      string    `gorm:"not null" json:"category"`
	Vendor        string    `json:"vendor"`
	Reference     string    `json:"reference"`
	PaymentMethod string    `json:"paymentMethod"`
	PaymentStatus string    `gorm:"not null;default:pending" json:"paymentStatus"`
	AccountID     *uint     `json:"accountId"`
	Notes         string    `json:"notes"`
	CreatedBy     *uint     `json:"createdBy"`
}

type Income struct {
	Base
	HotelID       *uint     `json:"hotelId"`
	IncomeDate    time.Time `gorm:"not null" json:"incomeDate"`
	Description   string    `gorm:"not null" json:"description"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Category      string    `gorm:"not null" json:"category"`
	Source        string    `json:"source"`
	Reference     string    `json:"reference"`
	PaymentMethod string    `json:"paymentMethod"`
	PaymentStatus string    `gorm:"not null;default:pending" json:"paymentStatus"`
	AccountID     *uint     `json:"accountId"`
	ReservationID *uint     `json:"reservationId"`
	Notes         string    `json:"notes"`
	CreatedBy     *uint     `json:"createdBy"`
}

type ParkingLot struct {
	Base
	HotelID     *uint   `json:"hotelId"`
	Name        string  `gorm:"not null" json:"name"`
	Location    string  `json:"location"`
	TotalSpots  int     `gorm:"not null;default:0" json:"totalSpots"`
	HourlyRate  float64 `gorm:"not null;default:0" json:"hourlyRate"`
	DailyRate   float64 `gorm:"not null;default:0" json:"dailyRate"`
	Status      string  `gorm:"not null;default:active" json:"status"`
	Description string  `json:"description"`
}

type ParkingSpot struct {
	Base
	LotID       *uint  `json:"lotId"`
	SpotNumber  string `gorm:"not null" json:"spotNumber"`
	Floor       string `json:"floor"`
	SpotType    string `gorm:"not null;default:standard" json:"spotType"`
	Status      string `gorm:"not null;default:available" json:"status"`
	IsCovered   bool   `gorm:"not null;default:false" json:"isCovered"`
	Description string `json:"description"`
}

type Vehicle struct {
	Base
	GuestID      *uint  `json:"guestId"`
	LicensePlate string `gorm:"not null" json:"licensePlate"`
	VehicleType  string `gorm:"not null;default:car" json:"vehicleType"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	IsRegistered bool   `gorm:"not null;default:true" json:"isRegistered"`
	Notes        string `json:"notes"`
}

type ParkingTransaction struct {
	Base
	LotID         *uint      `json:"lotId"`
	SpotID        *uint      `json:"spotId"`
	GuestID       *uint      `json:"guestId"`
	ReservationID *uint      `json:"reservationId"`
	LicensePlate  string     `gorm:"not null" json:"licensePlate"`
	EntryTime     time.Time  `gorm:"not null" json:"entryTime"`
	ExitTime      *time.Time `json:"exitTime"`
	HoursParked   *float64   `json:"hoursParked"`
	RateApplied   *float64   `json:"rateApplied"`
	AmountDue     float64    `gorm:"not null;default:0" json:"amountDue"`
	AmountPaid    float64    `gorm:"not null;default:0" json:"amountPaid"`
	Status        string     `gorm:"not null;default:active" json:"status"`
	PaymentStatus string     `gorm:"not null;default:pending" json:"paymentStatus"`
	PaymentMethod string     `json:"paymentMethod"`
	Notes         string     `json:"notes"`
}

type ParkingStats struct {
	Lots           int64 `json:"lots"`
	Spots          int64 `json:"spots"`
	AvailableSpots int64 `json:"availableSpots"`
}

func AllForDb() []any {
	return []any{
		&User{}, &Session{}, &Hotel{}, &Room{}, &Guest{}, &Reservation{}, &Payment{}, &Account{},
		&Expense{}, &Income{}, &ParkingLot{}, &ParkingSpot{}, &Vehicle{}, &ParkingTransaction{}, &Amenity{},
		&ParkingSpotType{}, &ParkingSpotStatus{}, &UserHotel{}, &Permission{}, &PermissionTemplate{},
		&UserPermission{}, &UserTemplate{}, &SanaGuest{}, &SanaRoomRack{},
	}
}

type SanitizedUser struct {
	ID         uint            `json:"id"`
	Email      string          `json:"email"`
	FirstName  string          `json:"firstName"`
	LastName   string          `json:"lastName"`
	UserHotels []UserHotelInfo `json:"userHotels"`
}

type UserHotelInfo struct {
	HotelID string `json:"hotelId"`
	Hotel   Hotel  `json:"hotel"`
}

func AllForTypeGen() []any {
	return []any{
		User{}, Session{}, Hotel{}, Room{}, Guest{}, GuestCompanion{}, Reservation{}, Account{},
		Expense{}, Income{}, ParkingLot{}, ParkingSpot{}, Vehicle{}, ParkingTransaction{}, Amenity{}, ParkingSpotType{}, ParkingSpotStatus{}, ParkingStats{}, SanitizedUser{}, UserHotelInfo{},
		Permission{}, PermissionTemplate{}, UserPermission{}, UserTemplate{},
	}
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
