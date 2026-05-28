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

type Guest struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	HotelID      *string   `json:"hotelId"`
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
	Email        string    `json:"email"`
	Landline     string    `json:"landline"`

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

type PaymentMethod struct {
	Base
	TranslateBase
}

type PaymentStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}
type ExpenseCategory struct {
	Base
	TranslateBase
}

type Expense struct {
	Base
	ExpenseDate time.Time `gorm:"not null" json:"expenseDate"`
	Description string    `gorm:"not null" json:"description"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Vendor      string    `json:"vendor"`
	Reference   string    `json:"reference"`
	Notes       string    `json:"notes"`
	CreatedBy   *uint     `json:"createdBy"`

	HotelID *uint `gorm:"not null" json:"hotelId"`
	Hotel   Hotel `gorm:"foreignKey:HotelID" json:"hotel"`

	AccountID *uint   `gorm:"not null" json:"accountId"`
	Account   Account `gorm:"foreignKey:AccountID" json:"account,omitzero"`

	CategoryID uint            `gorm:"not null" json:"categoryId"`
	Category   ExpenseCategory `gorm:"foreignKey:CategoryID" json:"category,omitzero" translate:"true"`

	PaymentStatusID uint          `gorm:"not null" json:"paymentStatusId"`
	PaymentStatus   PaymentStatus `gorm:"foreignKey:PaymentStatusID" json:"paymentStatus,omitzero"`

	PaymentMethodID uint   `gorm:"not null" json:"paymentMethodId"`
	PaymentMethod   string `gorm:"foreignKey:PaymentMethodID" json:"paymentMethod,omitzero" translate:"true"`
}

type IncomeCategory struct {
	Base
	TranslateBase
}

type Income struct {
	Base
	IncomeDate  time.Time `gorm:"not null" json:"incomeDate"`
	Description string    `gorm:"not null" json:"description"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Source      string    `json:"source"`
	Reference   string    `json:"reference"`
	Notes       string    `json:"notes"`
	CreatedBy   *uint     `json:"createdBy"`

	HotelID *uint `gorm:"not null" json:"hotelId"`
	Hotel   Hotel `gorm:"foreignKey:HotelID" json:"hotel"`

	AccountID *uint   `gorm:"not null" json:"accountId"`
	Account   Account `gorm:"foreignKey:AccountID" json:"account,omitzero"`

	ReservationID *uint       `gorm:"not null" json:"reservationId"`
	Reservation   Reservation `gorm:"foreignKey:ReservationID" json:"reservation"`

	PaymentStatusID uint          `gorm:"not null" json:"paymentStatusId"`
	PaymentStatus   PaymentStatus `gorm:"foreignKey:PaymentStatusID" json:"paymentStatus,omitzero"`

	PaymentMethodID uint   `gorm:"not null" json:"paymentMethodId"`
	PaymentMethod   string `gorm:"foreignKey:PaymentMethodID" json:"paymentMethod,omitzero" translate:"true"`

	CategoryID uint           `gorm:"not null" json:"categoryId"`
	Category   IncomeCategory `gorm:"foreignKey:CategoryID" json:"category,omitzero" translate:"true"`
}

type ParkingLotStatus struct {
	Base
	TranslateBase
}

type ParkingLot struct {
	Base
	Name        string  `gorm:"not null" json:"name"`
	Location    string  `json:"location"`
	TotalSpots  int     `gorm:"not null;default:0" json:"totalSpots"`
	HourlyRate  float64 `gorm:"not null;default:0" json:"hourlyRate"`
	DailyRate   float64 `gorm:"not null;default:0" json:"dailyRate"`
	Description string  `json:"description"`

	StatusID uint             `gorm:"not null" json:"statusId"`
	Status   ParkingLotStatus `gorm:"foreignKey:StatusID" json:"status,omitzero"`

	HotelID *uint `gorm:"not null" json:"hotelId"`
	Hotel   Hotel `gorm:"foreignKey:HotelID" json:"hotel"`
}

type ParkingSpotStatus struct {
	Base
	TranslateBase
	ColorHex string `gorm:"type:char(6);default:null" json:"colorHex,omitempty"`
}
type ParkingSpotType struct {
	Base
	TranslateBase
}

type ParkingSpot struct {
	Base
	SpotNumber  string `gorm:"not null" json:"spotNumber"`
	Floor       string `json:"floor"`
	IsCovered   bool   `gorm:"not null;default:false" json:"isCovered"`
	Description string `json:"description"`

	LotID *uint      `gorm:"not null" json:"lotId"`
	Lot   ParkingLot `gorm:"foreignKey:LotID" json:"lot,omitzero"`

	SpotTypeID uint            `gorm:"not null" json:"spotTypeId"`
	SpotType   ParkingSpotType `gorm:"foreignKey:SpotTypeID" json:"spotType,omitzero"`

	StatusID uint              `gorm:"not null" json:"statusId"`
	Status   ParkingSpotStatus `gorm:"foreignKey:StatusID" json:"status,omitzero"`
}

type VehicleType struct {
	Base
	TranslateBase
}

type Vehicle struct {
	Base
	LicensePlate string `gorm:"not null" json:"licensePlate"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	IsRegistered bool   `gorm:"not null;default:true" json:"isRegistered"`
	Notes        string `json:"notes"`

	VehicleTypeID uint        `gorm:"not null" json:"vehicleType"`
	VehicleType   VehicleType `gorm:"not null;foreignKey:VehicleTypeID" json:"vehicle"`

	GuestID uint  `gorm:"not null" json:"guestId"`
	Guest   Guest `gorm:"foreignKey:GuestID" json:"guest,omitzero"`
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
		&ParkingLotStatus{}, &ParkingSpotType{}, &ParkingSpotStatus{}, &UserHotel{}, &Permission{}, &PermissionTemplate{}, &UserPermission{}, &UserTemplate{}, &SanaGuest{}, &SanaRoomRack{},
		&TravelReason{}, &FamilyRelationship{}, &Nationality{}, &Country{}, &Occupation{},
		&InventoryItem{}, &RestaurantBill{}, &MealTransaction{}, &VehicleType{},
		&IncomeCategory{}, &ExpenseCategory{}, &PaymentStatus{}, &PaymentMethod{},
		&InventoryItemCategory{}, &InventoryItemUnit{}, &InventoryItemStatus{}, &RestaurantBillStatus{},
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
		User{},
		Session{},
		Hotel{},
		Room{},
		Guest{},
		GuestCompanion{},
		Reservation{},
		Account{},
		Expense{},
		Income{},
		ParkingLot{},
		ParkingSpot{},
		Vehicle{},
		ParkingTransaction{},
		Amenity{},
		ParkingSpotType{},
		ParkingSpotStatus{},
		ParkingStats{},
		SanitizedUser{},
		UserHotelInfo{},
		Permission{},
		PermissionTemplate{},
		UserPermission{},
		UserTemplate{},
		TravelReason{},
		FamilyRelationship{},
		Nationality{},
		Country{},
		Occupation{},
		InventoryItem{},
		RestaurantBill{},
		MealTransaction{},
		RestaurantStats{},
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

	HotelID *uint `gorm:"not null" json:"hotelId"`
	Hotel   Hotel `gorm:"foreignKey:HotelID" json:",omitzero"`

	ReservationID *uint       `gorm:"index" json:"reservationId"`
	Reservation   Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitzero"`

	GuestID *uint `gorm:"index" json:"guestId"`
	Guest   Guest `gorm:"foreignKey:GuestID" json:"guest,omitzero"`

	RoomID *uint `json:"roomId"`
	Room   Room  `gorm:"foreignKey:RoomID" json:"room,omitzero"`

	StatusID uint                 `gorm:"not null" json:"statusId"`
	Status   RestaurantBillStatus `gorm:"foreignKey:StatusID" json:"status,omitzero"`
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

	HotelID *uint `gorm:"not null" json:"hotelId"`
	Hotel   Hotel `gorm:"foreignKey:HotelID" json:"hotel"`
}
