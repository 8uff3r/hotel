package models

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
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	FirstName    string `gorm:"not null" json:"firstName"`
	LastName     string `gorm:"not null" json:"lastName"`
	Role         string `gorm:"not null;default:staff" json:"role"`
	IsActive     bool   `gorm:"not null;default:true" json:"isActive"`
}

type Session struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	User      User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type Hotel struct {
	Base
	Name    string `gorm:"not null" json:"name"`
	Address string `gorm:"not null" json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type Room struct {
	Base
	HotelID     *uint     `json:"hotelId"`
	RoomNumber  string    `gorm:"not null" json:"roomNumber"`
	RoomType    string    `gorm:"not null;default:single" json:"roomType"`
	Floor       *int      `json:"floor"`
	Capacity    int       `gorm:"not null;default:2" json:"capacity"`
	BasePrice   float64   `gorm:"not null;default:0" json:"basePrice"`
	Status      string    `gorm:"not null;default:available" json:"status"`
	Amenities   []Amenity `gorm:"many2many:room_amenities;" json:"amenities"`
	Description string    `json:"description"`

	GuestID uint `gorm:"not null;index" json:"guestId"`
}

type Translation map[string]string

type Amenity struct {
	Base
	Name        string      `json:"name"`
	Translation Translation `gorm:"type:jsonb" json:"translation"`
}

type ParkingSpotType struct {
	Base
	Name        string      `json:"name"`
	Translation Translation `gorm:"type:jsonb" json:"translation"`
}

type ParkingSpotStatus struct {
	Base
	Name        string      `json:"name"`
	Translation Translation `gorm:"type:jsonb" json:"translation"`
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

	Reservations []Reservation `gorm:"foreignKey:GuestID" json:"reservation,omitempty"`
}

type Reservation struct {
	Base
	ID      uint   `gorm:"primaryKey" json:"id"`
	GuestID uint   `gorm:"not null;index" json:"guestId"`
	Rooms   []Room `gorm:"many2many:reservation_rooms;"`

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

func AllPtr() []any {
	return []any{
		&User{}, &Session{}, &Hotel{}, &Room{}, &Guest{}, &Reservation{}, &Account{},
		&Expense{}, &Income{}, &ParkingLot{}, &ParkingSpot{}, &Vehicle{}, &ParkingTransaction{}, &Amenity{},
		&ParkingSpotType{}, &ParkingSpotStatus{},
	}
}

func All() []any {
	return []any{
		User{}, Session{}, Hotel{}, Room{}, Guest{}, Reservation{}, Account{},
		Expense{}, Income{}, ParkingLot{}, ParkingSpot{}, Vehicle{}, ParkingTransaction{}, Amenity{}, ParkingSpotType{}, ParkingSpotStatus{}, ParkingStats{},
	}
}
