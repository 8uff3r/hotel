import { pgTable, serial, text, timestamp, integer, boolean, index } from "drizzle-orm/pg-core";
import { hotels } from "./hotels";

export const accountsAccountTypes = ["asset", "liability", "equity", "revenue", "expense"] as const;

export const accounts = pgTable(
  "accounts",
  {
    id: serial("id").primaryKey(),
    hotelId: integer("hotel_id").references(() => hotels.id),
    accountCode: text("account_code").notNull(),
    accountName: text("account_name").notNull(),
    accountType: text("account_type", {
      enum: accountsAccountTypes,
    }).notNull(),
    accountSubType: text("account_sub_type"),
    parentId: integer("parent_id"),
    isActive: boolean("is_active").notNull().default(true),
    isSystem: boolean("is_system").notNull().default(false),
    description: text("description"),
    normalBalance: text("normal_balance", { enum: ["debit", "credit"] })
      .notNull()
      .default("debit"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [
    index("accounts_hotel_id_idx").on(table.hotelId),
    index("accounts_account_code_idx").on(table.accountCode),
    index("accounts_account_type_idx").on(table.accountType),
    index("accounts_parent_id_idx").on(table.parentId),
  ]
);

export type Account = typeof accounts.$inferSelect;
export type NewAccount = typeof accounts.$inferInsert;
