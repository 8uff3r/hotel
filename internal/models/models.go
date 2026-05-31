package models

//go:generate go run ./generate.go && go fmt ./generate.go

func AllForDB() []any {
	return []any{
		&User{}, &Session{}, &Hotel{}, &Room{}, &Guest{}, &Reservation{}, &Payment{}, &Account{},
		&Expense{}, &Income{}, &ParkingLot{}, &ParkingSpot{}, &Vehicle{}, &ParkingTransaction{}, &Amenity{},
		&ParkingLotStatus{}, &ParkingSpotType{}, &ParkingSpotStatus{}, &UserHotel{}, &Permission{}, &PermissionTemplate{}, &UserPermission{}, &UserTemplate{}, &SanaGuest{}, &SanaRoomRack{},
		&TravelReason{}, &FamilyRelationship{}, &Nationality{}, &Country{}, &SanaCity{}, &Occupation{},
		&InventoryItem{}, &RestaurantBill{}, &MealTransaction{}, &VehicleType{},
		&IncomeCategory{}, &ExpenseCategory{}, &PaymentStatus{}, &PaymentMethod{},
		&InventoryItemCategory{}, &InventoryItemUnit{}, &InventoryItemStatus{}, &RestaurantBillStatus{},
		&RoomType{}, &RoomStatus{}, &GuestCompanion{}, &PermissionCategory{}, &ReservationStatus{},
		&TravelAgency{},
	}
}

func AllForTypeGen() []any {
	return []any{
		User{},
		Session{},
		Hotel{},
		Room{},
		RoomType{},
		RoomStatus{},
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
		SanaCity{},
		Occupation{},
		InventoryItem{},
		RestaurantBill{},
		MealTransaction{},
		RestaurantStats{},
		TravelAgency{},
	}
}
