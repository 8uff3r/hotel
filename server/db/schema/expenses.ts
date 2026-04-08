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

export const expenseCategories = [
  "food_beverage",
  "housekeeping",
  "maintenance",
  "utilities",
  "salaries",
  "marketing",
  "supplies",
  "insurance",
  "taxes",
  "rent",
  "other",
] as const;

export const paymentMethods = [
  "cash",
  "credit_card",
  "debit_card",
  "bank_transfer",
  "check",
  "other",
] as const;

export const expensesPaymentStatuses = ["pending", "paid", "cancelled"] as const;

export const expenses = pgTable(
  "expenses",
  {
    id: serial("id").primaryKey(),
    hotelId: integer("hotel_id").references(() => hotels.id),
    expenseDate: timestamp("expense_date").notNull(),
    description: text("description").notNull(),
    amount: decimal("amount", { precision: 15, scale: 2 }).notNull(),
    category: text("category", {
      enum: expenseCategories,
    }).notNull(),
    vendor: text("vendor"),
    reference: text("reference"),
    paymentMethod: text("payment_method", {
      enum: paymentMethods,
    }),
    paymentStatus: text("payment_status", { enum: expensesPaymentStatuses })
      .notNull()
      .default("pending"),
    accountId: integer("account_id").references(() => accounts.id),
    receiptNumber: text("receipt_number"),
    notes: text("notes"),
    createdBy: integer("created_by").references(() => users.id),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [
    index("expenses_hotel_id_idx").on(table.hotelId),
    index("expenses_expense_date_idx").on(table.expenseDate),
    index("expenses_category_idx").on(table.category),
    index("expenses_payment_status_idx").on(table.paymentStatus),
  ]
);

export type Expense = typeof expenses.$inferSelect;
export type NewExpense = typeof expenses.$inferInsert;
