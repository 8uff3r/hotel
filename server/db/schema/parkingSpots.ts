import { pgTable, serial, text, timestamp, integer, boolean, index } from "drizzle-orm/pg-core";
import { parkingLots } from "./parkingLots";

export const parkingSpotStatuses = ["available", "occupied", "reserved", "maintenance"] as const;
export const spotTypes = ["standard", "handicap", "electric", "compact", "large"] as const;

export const parkingSpots = pgTable(
  "parking_spots",
  {
    id: serial("id").primaryKey(),
    lotId: integer("lot_id").references(() => parkingLots.id),
    spotNumber: text("spot_number").notNull(),
    floor: text("floor"),
    spotType: text("spot_type", { enum: spotTypes }).notNull().default("standard"),
    status: text("status", { enum: parkingSpotStatuses }).notNull().default("available"),
    isCovered: boolean("is_covered").notNull().default(false),
    description: text("description"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [
    index("parking_spots_lot_id_idx").on(table.lotId),
    index("parking_spots_status_idx").on(table.status),
    index("parking_spots_spot_number_idx").on(table.spotNumber),
  ]
);

export type ParkingSpot = typeof parkingSpots.$inferSelect;
export type NewParkingSpot = typeof parkingSpots.$inferInsert;
