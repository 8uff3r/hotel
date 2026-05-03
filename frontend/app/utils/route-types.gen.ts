/* Do not change, this code is generated from Golang structs */

export interface PermissionCategory {
  id?: number;
  slug: string;
  label: string;
}
export interface Permission {
  id?: number;
  resource: string;
  action: string;
  categoryId: number;
  category: PermissionCategory;
}
export interface UserPermission {
  id?: number;
  userId: number;
  hotelId?: string;
  permissionId: number;
  permission: Permission;
  granted: boolean;
}
export interface Hotel {
  id: string;
  name: string;
  address: string;
  phone: string;
  email: string;
}
export interface UserHotel {
  id?: number;
  userId: number;
  hotelId: string;
  hotel: Hotel;
}
export interface User {
  id?: number;
  email: string;
  firstName: string;
  lastName: string;
  userHotels: UserHotel[];
  isActive: boolean;
  permissions?: UserPermission[];
}
export interface Time {}
export interface Session {
  id: string;
  userId: number;
  expiresAt: Time;
  createdAt: Time;
}

export interface RoomStatus {
  id?: number;
  slug: string;
  label: string;
  colorHex?: string;
}
export interface RoomType {
  id?: number;
  slug: string;
  label: string;
  colorHex?: string;
}
export interface Amenity {
  id?: number;
  slug: string;
  label: string;
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
}
export interface GuestCompanion {
  id?: number;
  guestId: number;
  firstName: string;
  lastName: string;
  nationalId: string;
  idNumber: string;
  relation: string;
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
  hotelId?: string;
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
  reservations?: Reservation[];
  companions?: GuestCompanion[];
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
export interface PaymentStatus {
  id?: number;
  slug: string;
  label: string;
  colorHex?: string;
}
export interface ExpenseCategory {
  id?: number;
  slug: string;
  label: string;
}
export interface Expense {
  id?: number;
  expenseDate: Time;
  description: string;
  amount: number;
  vendor: string;
  reference: string;
  notes: string;
  createdBy?: number;
  hotelId?: number;
  hotel: Hotel;
  accountId?: number;
  account: Account;
  categoryId: number;
  category: ExpenseCategory;
  paymentStatusId: number;
  paymentStatus: PaymentStatus;
  paymentMethodId: number;
  paymentMethod: string;
}
export interface IncomeCategory {
  id?: number;
  slug: string;
  label: string;
}
export interface Income {
  id?: number;
  incomeDate: Time;
  description: string;
  amount: number;
  source: string;
  reference: string;
  notes: string;
  createdBy?: number;
  hotelId?: number;
  hotel: Hotel;
  accountId?: number;
  account: Account;
  reservationId?: number;
  reservation: Reservation;
  paymentStatusId: number;
  paymentStatus: PaymentStatus;
  paymentMethodId: number;
  paymentMethod: string;
  categoryId: number;
  category: IncomeCategory;
}
export interface ParkingLotStatus {
  id?: number;
  slug: string;
  label: string;
}
export interface ParkingLot {
  id?: number;
  name: string;
  location: string;
  totalSpots: number;
  hourlyRate: number;
  dailyRate: number;
  description: string;
  statusId: number;
  status: ParkingLotStatus;
  hotelId?: number;
  hotel: Hotel;
}
export interface ParkingSpotStatus {
  id?: number;
  slug: string;
  label: string;
  colorHex?: string;
}
export interface ParkingSpotType {
  id?: number;
  slug: string;
  label: string;
}
export interface ParkingSpot {
  id?: number;
  spotNumber: string;
  floor: string;
  isCovered: boolean;
  description: string;
  lotId?: number;
  lot: ParkingLot;
  spotTypeId: number;
  spotType: ParkingSpotType;
  statusId: number;
  status: ParkingSpotStatus;
}
export interface VehicleType {
  id?: number;
  slug: string;
  label: string;
}
export interface Vehicle {
  id?: number;
  licensePlate: string;
  make: string;
  model: string;
  color: string;
  isRegistered: boolean;
  notes: string;
  vehicleType: number;
  vehicle: VehicleType;
  guestId: number;
  guest: Guest;
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

export interface ParkingStats {
  lots: number;
  spots: number;
  availableSpots: number;
}
export interface UserHotelInfo {
  hotelId: string;
  hotel: Hotel;
}
export interface SanitizedUser {
  id?: number;
  email: string;
  firstName: string;
  lastName: string;
  userHotels: UserHotelInfo[];
}

export interface PermissionTemplate {
  id?: number;
  slug: string;
  label: string;
  description: string;
  permissions: Permission[];
}

export interface UserTemplate {
  id?: number;
  userId: number;
  templateId: number;
  template: PermissionTemplate;
}
export interface TravelReason {
  id?: number;
  slug: string;
  label: string;
  sanaId: string;
  sanaName: string;
}
export interface FamilyRelationship {
  id?: number;
  slug: string;
  label: string;
  sanaId: string;
  sanaName: string;
}
export interface Nationality {
  id?: number;
  slug: string;
  label: string;
  sanaId: string;
  sanaName: string;
}
export interface Country {
  id?: number;
  slug: string;
  label: string;
  sanaId: string;
  sanaName: string;
  isIran: boolean;
}
export interface Occupation {
  id?: number;
  slug: string;
  label: string;
  sanaId: string;
  sanaName: string;
}
export interface InventoryItemStatus {
  id?: number;
  slug: string;
  label: string;
  colorHex?: string;
}
export interface InventoryItemCategory {
  id?: number;
  slug: string;
  label: string;
}
export interface InventoryItemUnit {
  id?: number;
  slug: string;
  label: string;
}
export interface InventoryItem {
  id?: number;
  name: string;
  quantity: number;
  unitCost: number;
  reorderLevel: number;
  isActive: boolean;
  description: string;
  unitId: number;
  unit: InventoryItemUnit;
  categoryId: number;
  category: InventoryItemCategory;
  statusId: number;
  status: InventoryItemStatus;
  hotelId?: number;
}
export interface RestaurantBillStatus {
  id?: number;
  slug: string;
  label: string;
  colorHex?: string;
}
export interface RestaurantBill {
  id?: number;
  billDate: Time;
  subtotal: number;
  taxAmount: number;
  discountAmount: number;
  totalAmount: number;
  isExternal: boolean;
  externalRestaurant: string;
  notes: string;
  settled: boolean;
  settledAt?: Time;
  settledBy?: number;
  hotelId?: number;
  reservationId?: number;
  reservation: Reservation;
  guestId?: number;
  guest: Guest;
  roomId?: number;
  room: Room;
  statusId: number;
  status: RestaurantBillStatus;
}
export interface MealTransaction {
  id?: number;
  itemName: string;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
  isExternal: boolean;
  notes: string;
  inventoryItemId?: number;
  inventoryItem: InventoryItem;
  billId: number;
  hotelId?: number;
  hotel: Hotel;
}
export interface RestaurantStats {
  totalBills: number;
  totalRevenue: number;
  internalRevenue: number;
  externalRevenue: number;
  totalMeals: number;
}
