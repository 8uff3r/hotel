import { pgTable, serial, text, timestamp, integer, real, index } from "drizzle-orm/pg-core";
import { reservations } from "./reservations";
import { parkingSpots } from "./parkingSpots";
import { parkingLots } from "./parkingLots";
import { guests } from "./guests";

export const parkingTransactionStatuses = ["active", "completed", "cancelled"] as const;
export const paymentStatuses = ["pending", "paid", "waived"] as const;

export const parkingTransactions = pgTable(
  "parking_transactions",
  {
    id: serial("id").primaryKey(),
    lotId: integer("lot_id").references(() => parkingLots.id),
    spotId: integer("spot_id").references(() => parkingSpots.id),
    guestId: integer("guest_id").references(() => guests.id),
    reservationId: integer("reservation_id").references(() => reservations.id),
    licensePlate: text("license_plate").notNull(),
    entryTime: timestamp("entry_time").notNull(),
    exitTime: timestamp("exit_time"),
    hoursParked: real("hours_parked"),
    rateApplied: text("rate_applied"),
    amountDue: real("amount_due").notNull().default(0),
    amountPaid: real("amount_paid").notNull().default(0),
    status: text("status", { enum: parkingTransactionStatuses }).notNull().default("active"),
    paymentStatus: text("payment_status", { enum: paymentStatuses }).notNull().default("pending"),
    paymentMethod: text("payment_method"),
    notes: text("notes"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [
    index("parking_transactions_lot_id_idx").on(table.lotId),
    index("parking_transactions_guest_id_idx").on(table.guestId),
    index("parking_transactions_reservation_id_idx").on(table.reservationId),
    index("parking_transactions_status_idx").on(table.status),
    index("parking_transactions_license_plate_idx").on(table.licensePlate),
  ]
);

export type ParkingTransaction = typeof parkingTransactions.$inferSelect;
export type NewParkingTransaction = typeof parkingTransactions.$inferInsert;
