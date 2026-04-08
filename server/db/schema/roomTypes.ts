import { pgTable, serial, text, real, timestamp, index } from "drizzle-orm/pg-core";

export const roomTypes = pgTable(
  "room_types",
  {
    id: serial("id").primaryKey(),
    name: text("name").notNull().unique(),
    description: text("description"),
    basePriceMultiplier: real("base_price_multiplier").notNull().default(1.0),
    amenities: text("amenities"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [index("room_types_name_idx").on(table.name)]
);

export type RoomType = typeof roomTypes.$inferSelect;
export type NewRoomType = typeof roomTypes.$inferInsert;
