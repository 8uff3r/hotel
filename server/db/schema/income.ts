import {
  pgTable,
  serial,
  text,
  timestamp,
  integer,
  decimal,
  index,
  foreignKey,
} from "drizzle-orm/pg-core";
import { users } from "./users";
import { hotels } from "./hotels";
import { accounts } from "./accounts";

export const incomeCategories = [
  "room_revenue",
  "food_beverage",
  "laundry",
  "spa",
  "meeting_rooms",
  "other",
] as const;

export const incomePaymentMethods = [
  "cash",
  "credit_card",
  "debit_card",
  "bank_transfer",
  "other",
] as const;

export const incomePaymentStatuses = ["pending", "received", "refunded"] as const;

export const income = pgTable(
  "income",
  {
    id: serial("id").primaryKey(),
    hotelId: integer("hotel_id").references(() => hotels.id),
    incomeDate: timestamp("income_date").notNull(),
    description: text("description").notNull(),
    amount: decimal("amount", { precision: 15, scale: 2 }).notNull(),
    category: text("category", {
      enum: incomeCategories,
    }).notNull(),
    source: text("source"),
    reference: text("reference"),
    paymentMethod: text("payment_method", {
      enum: incomePaymentMethods,
    }),
    paymentStatus: text("payment_status", { enum: incomePaymentStatuses })
      .notNull()
      .default("pending"),
    accountId: integer("account_id").references(() => accounts.id),
    reservationId: integer("reservation_id"),
    notes: text("notes"),
    createdBy: integer("created_by").references(() => users.id),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [
    index("income_hotel_id_idx").on(table.hotelId),
    index("income_income_date_idx").on(table.incomeDate),
    index("income_category_idx").on(table.category),
    index("income_payment_status_idx").on(table.paymentStatus),
    index("income_reservation_id_idx").on(table.reservationId),
  ]
);

export type Income = typeof income.$inferSelect;
export type NewIncome = typeof income.$inferInsert;
