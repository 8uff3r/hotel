/* Do not change, this code is generated from Golang structs */

export interface User {
  id?: number;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
  isActive: boolean;
}
export interface Time {}
export interface Session {
  id: string;
  userId: number;
  expiresAt: Time;
  createdAt: Time;
}
export interface Hotel {
  id?: number;
  name: string;
  address: string;
  phone: string;
  email: string;
}
export interface Amenity {
  id?: number;
  name: string;
}
export interface Room {
  id?: number;
  hotelId?: number;
  roomNumber: string;
  roomType: string;
  floor?: number;
  capacity: number;
  basePrice: number;
  status: string;
  amenities: Amenity[];
  description: string;
}
export interface Guest {
  id?: number;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  idType: string;
  idNumber: string;
  address: string;
  city: string;
  country: string;
  notes: string;
}
export interface Reservation {
  id?: number;
  hotelId?: number;
  guestId: number;
  roomId: number;
  checkInDate: Time;
  checkOutDate: Time;
  actualCheckIn?: Time;
  actualCheckOut?: Time;
  status: string;
  totalAmount: number;
  paidAmount: number;
  paymentStatus: string;
  specialRequests: string;
  numberOfGuests: number;
  createdBy?: number;
}
export interface Account {
  id?: number;
  hotelId?: number;
  accountCode: string;
  accountName: string;
  accountType: string;
  accountSubType: string;
  parentId?: number;
  isActive: boolean;
  isSystem: boolean;
  description: string;
  normalBalance: string;
}
export interface Expense {
  id?: number;
  hotelId?: number;
  expenseDate: Time;
  description: string;
  amount: number;
  category: string;
  vendor: string;
  reference: string;
  paymentMethod: string;
  paymentStatus: string;
  accountId?: number;
  notes: string;
  createdBy?: number;
}
export interface Income {
  id?: number;
  hotelId?: number;
  incomeDate: Time;
  description: string;
  amount: number;
  category: string;
  source: string;
  reference: string;
  paymentMethod: string;
  paymentStatus: string;
  accountId?: number;
  reservationId?: number;
  notes: string;
  createdBy?: number;
}
export interface ParkingLot {
  id?: number;
  hotelId?: number;
  name: string;
  location: string;
  totalSpots: number;
  hourlyRate: number;
  dailyRate: number;
  status: string;
  description: string;
}
export interface ParkingSpot {
  id?: number;
  lotId?: number;
  spotNumber: string;
  floor: string;
  spotType: string;
  status: string;
  isCovered: boolean;
  description: string;
}
export interface Vehicle {
  id?: number;
  guestId?: number;
  licensePlate: string;
  vehicleType: string;
  make: string;
  model: string;
  color: string;
  isRegistered: boolean;
  notes: string;
}
export interface ParkingTransaction {
  id?: number;
  lotId?: number;
  spotId?: number;
  guestId?: number;
  reservationId?: number;
  licensePlate: string;
  entryTime: Time;
  exitTime?: Time;
  hoursParked?: number;
  rateApplied?: number;
  amountDue: number;
  amountPaid: number;
  status: string;
  paymentStatus: string;
  paymentMethod: string;
  notes: string;
}
