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

export const journalEntries = pgTable(
  "journal_entries",
  {
    id: serial("id").primaryKey(),
    hotelId: integer("hotel_id").references(() => hotels.id),
    entryNumber: text("entry_number").notNull(),
    entryDate: timestamp("entry_date").notNull(),
    description: text("description").notNull(),
    reference: text("reference"),
    status: text("status", { enum: ["draft", "posted", "voided"] })
      .notNull()
      .default("draft"),
    totalDebit: decimal("total_debit", { precision: 15, scale: 2 }).notNull().default("0"),
    totalCredit: decimal("total_credit", { precision: 15, scale: 2 }).notNull().default("0"),
    createdBy: integer("created_by").references(() => users.id),
    postedAt: timestamp("posted_at"),
    voidedAt: timestamp("voided_at"),
    voidReason: text("void_reason"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [
    index("journal_entries_hotel_id_idx").on(table.hotelId),
    index("journal_entries_entry_number_idx").on(table.entryNumber),
    index("journal_entries_entry_date_idx").on(table.entryDate),
    index("journal_entries_status_idx").on(table.status),
  ]
);

export type JournalEntry = typeof journalEntries.$inferSelect;
export type NewJournalEntry = typeof journalEntries.$inferInsert;
