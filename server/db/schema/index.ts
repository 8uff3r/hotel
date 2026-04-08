import { vehicles, vehicleTypes } from "./vehicles";
import { users } from "./users";
import { userRolesTable, userRoles } from "./userRoles";
import { roomTypes } from "./roomTypes";
import { rooms, roomStatuses, roomsRoomTypes as roomTypeEnum } from "./rooms";
import { reservations, reservationStatuses, reservationPaymentStatuses } from "./reservations";
import {
  parkingTransactions,
  parkingTransactionStatuses,
  paymentStatuses,
} from "./parkingTransactions";
import { parkingSpots, parkingSpotStatuses, spotTypes } from "./parkingSpots";
import { parkingLots, parkingLotStatuses } from "./parkingLots";
import { journalLines } from "./journalLines";
import { journalEntries } from "./journalEntries";
import { income } from "./income";
import { hotels } from "./hotels";
import { guests } from "./guests";
import {
  expenses,
  expenseCategories,
  paymentMethods,
  expensesPaymentStatuses as expensePaymentStatuses,
} from "./expenses";
import { accounts } from "./accounts";

// Export all tables
export const tables = {
  users,
  hotels,
  rooms,
  guests,
  reservations,
  roomTypes,
  accounts,
  journalEntries,
  journalLines,
  expenseCategories,
  paymentMethods,
  expensePaymentStatuses,
  income,
  expenses,
  parkingLots,
  parkingSpots,
  vehicles,
  parkingTransactions,
  userRoles: userRolesTable,
};

// Export types
export * from "./users";
export * from "./hotels";
export * from "./rooms";
export * from "./guests";
export * from "./reservations";
export * from "./roomTypes";
export * from "./accounts";
export * from "./journalEntries";
export * from "./journalLines";
export * from "./expenses";
export * from "./income";
export * from "./parkingLots";
export * from "./parkingSpots";
export * from "./vehicles";
export * from "./parkingTransactions";
export * from "./userRoles";
