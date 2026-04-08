package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	FirstName    string    `gorm:"not null" json:"first_name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	Role         string    `gorm:"not null;default:staff" json:"role"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Session struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Hotel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Address   string    `gorm:"not null" json:"address"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Room struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	HotelID     *uint     `json:"hotel_id"`
	RoomNumber  string    `gorm:"not null" json:"room_number"`
	RoomType    string    `gorm:"not null;default:single" json:"room_type"`
	Floor       *int      `json:"floor"`
	Capacity    int       `gorm:"not null;default:2" json:"capacity"`
	BasePrice   float64   `gorm:"not null;default:0" json:"base_price"`
	Status      string    `gorm:"not null;default:available" json:"status"`
	Amenities   string    `json:"amenities"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Guest struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	IDType    string    `json:"id_type"`
	IDNumber  string    `json:"id_number"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Reservation struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	HotelID         *uint      `json:"hotel_id"`
	GuestID         uint       `gorm:"not null" json:"guest_id"`
	RoomID          uint       `gorm:"not null" json:"room_id"`
	CheckInDate     time.Time  `gorm:"not null" json:"check_in_date"`
	CheckOutDate    time.Time  `gorm:"not null" json:"check_out_date"`
	ActualCheckIn   *time.Time `json:"actual_check_in"`
	ActualCheckOut  *time.Time `json:"actual_check_out"`
	Status          string     `gorm:"not null;default:pending" json:"status"`
	TotalAmount     float64    `gorm:"not null;default:0" json:"total_amount"`
	PaidAmount      float64    `gorm:"not null;default:0" json:"paid_amount"`
	PaymentStatus   string     `gorm:"not null;default:pending" json:"payment_status"`
	SpecialRequests string     `json:"special_requests"`
	NumberOfGuests  int        `gorm:"not null;default:1" json:"number_of_guests"`
	CreatedBy       *uint      `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Account struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	HotelID        *uint     `json:"hotel_id"`
	AccountCode    string    `gorm:"not null" json:"account_code"`
	AccountName    string    `gorm:"not null" json:"account_name"`
	AccountType    string    `gorm:"not null" json:"account_type"`
	AccountSubType string    `json:"account_sub_type"`
	ParentID       *uint     `json:"parent_id"`
	IsActive       bool      `gorm:"not null;default:true" json:"is_active"`
	IsSystem       bool      `gorm:"not null;default:false" json:"is_system"`
	Description    string    `json:"description"`
	NormalBalance  string    `gorm:"not null;default:debit" json:"normal_balance"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Expense struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	HotelID       *uint     `json:"hotel_id"`
	ExpenseDate   time.Time `gorm:"not null" json:"expense_date"`
	Description   string    `gorm:"not null" json:"description"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Category      string    `gorm:"not null" json:"category"`
	Vendor        string    `json:"vendor"`
	Reference     string    `json:"reference"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `gorm:"not null;default:pending" json:"payment_status"`
	AccountID     *uint     `json:"account_id"`
	Notes         string    `json:"notes"`
	CreatedBy     *uint     `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Income struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	HotelID       *uint     `json:"hotel_id"`
	IncomeDate    time.Time `gorm:"not null" json:"income_date"`
	Description   string    `gorm:"not null" json:"description"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Category      string    `gorm:"not null" json:"category"`
	Source        string    `json:"source"`
	Reference     string    `json:"reference"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `gorm:"not null;default:pending" json:"payment_status"`
	AccountID     *uint     `json:"account_id"`
	ReservationID *uint     `json:"reservation_id"`
	Notes         string    `json:"notes"`
	CreatedBy     *uint     `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ParkingLot struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	HotelID     *uint     `json:"hotel_id"`
	Name        string    `gorm:"not null" json:"name"`
	Location    string    `json:"location"`
	TotalSpots  int       `gorm:"not null;default:0" json:"total_spots"`
	HourlyRate  float64   `gorm:"not null;default:0" json:"hourly_rate"`
	DailyRate   float64   `gorm:"not null;default:0" json:"daily_rate"`
	Status      string    `gorm:"not null;default:active" json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ParkingSpot struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LotID       *uint     `json:"lot_id"`
	SpotNumber  string    `gorm:"not null" json:"spot_number"`
	Floor       string    `json:"floor"`
	SpotType    string    `gorm:"not null;default:standard" json:"spot_type"`
	Status      string    `gorm:"not null;default:available" json:"status"`
	IsCovered   bool      `gorm:"not null;default:false" json:"is_covered"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Vehicle struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	GuestID      *uint     `json:"guest_id"`
	LicensePlate string    `gorm:"not null" json:"license_plate"`
	VehicleType  string    `gorm:"not null;default:car" json:"vehicle_type"`
	Make         string    `json:"make"`
	Model        string    `json:"model"`
	Color        string    `json:"color"`
	IsRegistered bool      `gorm:"not null;default:true" json:"is_registered"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ParkingTransaction struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	LotID         *uint      `json:"lot_id"`
	SpotID        *uint      `json:"spot_id"`
	GuestID       *uint      `json:"guest_id"`
	ReservationID *uint      `json:"reservation_id"`
	LicensePlate  string     `gorm:"not null" json:"license_plate"`
	EntryTime     time.Time  `gorm:"not null" json:"entry_time"`
	ExitTime      *time.Time `json:"exit_time"`
	HoursParked   *float64   `json:"hours_parked"`
	RateApplied   *float64   `json:"rate_applied"`
	AmountDue     float64    `gorm:"not null;default:0" json:"amount_due"`
	AmountPaid    float64    `gorm:"not null;default:0" json:"amount_paid"`
	Status        string     `gorm:"not null;default:active" json:"status"`
	PaymentStatus string     `gorm:"not null;default:pending" json:"payment_status"`
	PaymentMethod string     `json:"payment_method"`
	Notes         string     `json:"notes"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func All() []any {
	return []any{
		&User{}, &Session{}, &Hotel{}, &Room{}, &Guest{}, &Reservation{}, &Account{},
		&Expense{}, &Income{}, &ParkingLot{}, &ParkingSpot{}, &Vehicle{}, &ParkingTransaction{},
	}
}
