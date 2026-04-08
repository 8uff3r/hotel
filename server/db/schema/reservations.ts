import { pgTable, serial, text, timestamp, real, integer, index } from "drizzle-orm/pg-core";
import { users } from "./users";
import { rooms } from "./rooms";
import { hotels } from "./hotels";
import { guests } from "./guests";

export const reservationStatuses = [
  "pending",
  "confirmed",
  "checked_in",
  "checked_out",
  "cancelled",
  "no_show",
] as const;

export const reservationPaymentStatuses = ["pending", "partial", "paid", "refunded"] as const;

export const reservations = pgTable(
  "reservations",
  {
    id: serial("id").primaryKey(),
    hotelId: integer("hotel_id").references(() => hotels.id),
    guestId: integer("guest_id")
      .references(() => guests.id)
      .notNull(),
    roomId: integer("room_id")
      .references(() => rooms.id)
      .notNull(),
    checkInDate: timestamp("check_in_date").notNull(),
    checkOutDate: timestamp("check_out_date").notNull(),
    actualCheckIn: timestamp("actual_check_in"),
    actualCheckOut: timestamp("actual_check_out"),
    status: text("status", {
      enum: reservationStatuses,
    })
      .notNull()
      .default("pending"),
    totalAmount: real("total_amount").notNull().default(0),
    paidAmount: real("paid_amount").notNull().default(0),
    paymentStatus: text("payment_status", { enum: reservationPaymentStatuses })
      .notNull()
      .default("pending"),
    specialRequests: text("special_requests"),
    numberOfGuests: integer("number_of_guests").notNull().default(1),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
    createdBy: integer("created_by").references(() => users.id),
  },
  (table) => [
    index("reservations_hotel_id_idx").on(table.hotelId),
    index("reservations_guest_id_idx").on(table.guestId),
    index("reservations_room_id_idx").on(table.roomId),
    index("reservations_status_idx").on(table.status),
    index("reservations_check_in_date_idx").on(table.checkInDate),
    index("reservations_check_out_date_idx").on(table.checkOutDate),
  ]
);

export type Reservation = typeof reservations.$inferSelect;
export type NewReservation = typeof reservations.$inferInsert;
