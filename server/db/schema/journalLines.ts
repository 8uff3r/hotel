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
import { journalEntries } from "./journalEntries";
import { accounts } from "./accounts";

export const journalLines = pgTable(
  "journal_lines",
  {
    id: serial("id").primaryKey(),
    entryId: integer("entry_id")
      .references(() => journalEntries.id)
      .notNull(),
    accountId: integer("account_id")
      .references(() => accounts.id)
      .notNull(),
    description: text("description"),
    debit: decimal("debit", { precision: 15, scale: 2 }).notNull().default("0"),
    credit: decimal("credit", { precision: 15, scale: 2 }).notNull().default("0"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
  },
  (table) => [
    index("journal_lines_entry_id_idx").on(table.entryId),
    index("journal_lines_account_id_idx").on(table.accountId),
  ]
);

export type JournalLine = typeof journalLines.$inferSelect;
export type NewJournalLine = typeof journalLines.$inferInsert;
