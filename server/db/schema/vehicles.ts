import { pgTable, serial, text, timestamp, integer, index } from "drizzle-orm/pg-core";
import { guests } from "./guests";

export const vehicleTypes = ["car", "motorcycle", "truck", "van", "other"] as const;

export const vehicles = pgTable(
  "vehicles",
  {
    id: serial("id").primaryKey(),
    guestId: integer("guest_id").references(() => guests.id),
    licensePlate: text("license_plate").notNull(),
    vehicleType: text("vehicle_type", { enum: vehicleTypes }).notNull().default("car"),
    make: text("make"),
    model: text("model"),
    color: text("color"),
    isRegistered: integer("is_registered").notNull().default(1),
    notes: text("notes"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [
    index("vehicles_guest_id_idx").on(table.guestId),
    index("vehicles_license_plate_idx").on(table.licensePlate),
  ]
);

export type Vehicle = typeof vehicles.$inferSelect;
export type NewVehicle = typeof vehicles.$inferInsert;
