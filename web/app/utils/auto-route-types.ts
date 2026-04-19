/* Do not change, this code is generated from Golang structs */

export interface Role {
  id?: number;
  name: string;
  translation?: { [key: string]: string };
}
export interface User {
  id?: number;
  email: string;
  firstName: string;
  lastName: string;
  roles: Role[];
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
export interface RoomStatus {
  id?: number;
  name: string;
  translation?: { [key: string]: string };
  colorHex?: string;
}
export interface RoomType {
  id?: number;
  name: string;
  translation?: { [key: string]: string };
  colorHex?: string;
}
export interface Amenity {
  id?: number;
  name: string;
  translation?: { [key: string]: string };
}
export interface Room {
  id?: number;
  hotelId?: number;
  name: string;
  roomNumber: string;
  floor: number;
  capacity: number;
  basePrice: number;
  amenities: Amenity[];
  description: string;
  roomTypeId: number;
  roomType: RoomType;
  statusId: number;
  status: RoomStatus;
  guestId: number;
}
export interface Payment {
  id?: number;
  reservationId: number;
  isCash: boolean;
  agency: boolean;
  referrer: string;
  contractType: string;
}
export interface Reservation {
  id?: number;
  guestId: number;
  rooms: Room[];
  reservationCode: string;
  entryDate: Time;
  departureDate: Time;
  durationOfStay: number;
  numberOfPeople: number;
  origin: string;
  destination: string;
  purposeOfTravel: string;
  breakfast: boolean;
  guide: boolean;
  roomPrice: number;
  userCheckIn: string;
  userCheckOut: string;
  notes: string;
  payment: Payment;
}
export interface Guest {
  id?: number;
  hotelId?: number;
  firstName: string;
  lastName: string;
  fatherName: string;
  nationalId: string;
  idNumber: string;
  nationality: string;
  gender: string;
  dateOfBirth: Time;
  placeOfBirth: string;
  phone: string;
  address: string;
  postalCode: string;
  occupation: string;
  reservation?: Reservation[];
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

export interface ParkingSpotType {
  id?: number;
  name: string;
  translation?: { [key: string]: string };
}
export interface ParkingSpotStatus {
  id?: number;
  name: string;
  translation?: { [key: string]: string };
}
export interface ParkingStats {
  lots: number;
  spots: number;
  availableSpots: number;
}
export interface SanitizedUser {
  id?: number;
  email: string;
  firstName: string;
  lastName: string;
  roles: Role[];
}
