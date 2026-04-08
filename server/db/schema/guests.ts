import { pgTable, serial, text, timestamp, index } from "drizzle-orm/pg-core";

export const guests = pgTable(
  "guests",
  {
    id: serial("id").primaryKey(),
    firstName: text("first_name").notNull(),
    lastName: text("last_name").notNull(),
    email: text("email"),
    phone: text("phone"),
    idType: text("id_type", { enum: ["passport", "national_id", "driver_license", "other"] }),
    idNumber: text("id_number"),
    address: text("address"),
    city: text("city"),
    country: text("country"),
    notes: text("notes"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [
    index("guests_email_idx").on(table.email),
    index("guests_phone_idx").on(table.phone),
    index("guests_name_idx").on(table.lastName, table.firstName),
  ]
);

export type Guest = typeof guests.$inferSelect;
export type NewGuest = typeof guests.$inferInsert;
