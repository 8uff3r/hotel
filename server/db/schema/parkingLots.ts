import { pgTable, serial, text, timestamp, integer, index } from "drizzle-orm/pg-core";
import { hotels } from "./hotels";

export const parkingLotStatuses = ["active", "full", "closed"] as const;

export const parkingLots = pgTable(
  "parking_lots",
  {
    id: serial("id").primaryKey(),
    hotelId: integer("hotel_id").references(() => hotels.id),
    name: text("name").notNull(),
    location: text("location"),
    totalSpots: integer("total_spots").notNull().default(0),
    hourlyRate: text("hourly_rate").notNull().default("0"),
    dailyRate: text("daily_rate").notNull().default("0"),
    status: text("status", { enum: parkingLotStatuses }).notNull().default("active"),
    description: text("description"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [
    index("parking_lots_hotel_id_idx").on(table.hotelId),
    index("parking_lots_status_idx").on(table.status),
  ]
);

export type ParkingLot = typeof parkingLots.$inferSelect;
export type NewParkingLot = typeof parkingLots.$inferInsert;
