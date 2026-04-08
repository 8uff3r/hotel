import { pgTable, serial, text, timestamp, index } from "drizzle-orm/pg-core";

export const hotels = pgTable(
  "hotels",
  {
    id: serial("id").primaryKey(),
    name: text("name").notNull(),
    address: text("address").notNull(),
    phone: text("phone"),
    email: text("email"),
    logoUrl: text("logo_url"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [index("hotels_name_idx").on(table.name)]
);

export type Hotel = typeof hotels.$inferSelect;
export type NewHotel = typeof hotels.$inferInsert;
