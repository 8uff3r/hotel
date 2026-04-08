import { pgTable, serial, text, timestamp, real, integer, index } from "drizzle-orm/pg-core";
import { hotels } from "./hotels";

export const roomStatuses = ["available", "occupied", "maintenance", "out_of_order"] as const;
export const roomsRoomTypes = ["single", "double", "suite", "deluxe"] as const;

export const rooms = pgTable(
  "rooms",
  {
    id: serial("id").primaryKey(),
    hotelId: integer("hotel_id").references(() => hotels.id),
    roomNumber: text("room_number").notNull(),
    roomType: text("room_type", { enum: roomsRoomTypes }).notNull().default("single"),
    floor: integer("floor"),
    capacity: integer("capacity").notNull().default(2),
    basePrice: real("base_price").notNull().default(0),
    status: text("status", { enum: roomStatuses }).notNull().default("available"),
    amenities: text("amenities"),
    description: text("description"),
    images: text("images"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [
    index("rooms_hotel_id_idx").on(table.hotelId),
    index("rooms_room_number_idx").on(table.roomNumber),
    index("rooms_status_idx").on(table.status),
    index("rooms_room_type_idx").on(table.roomType),
  ]
);

export type Room = typeof rooms.$inferSelect;
export type NewRoom = typeof rooms.$inferInsert;
