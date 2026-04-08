import { tables } from "~~/server/db/schema";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  const [vehicle] = await db
    .insert(tables.vehicles)
    .values({
      guestId: body.guestId,
      licensePlate: body.licensePlate,
      vehicleType: body.vehicleType || "car",
      make: body.make,
      model: body.model,
      color: body.color,
      isRegistered: body.isRegistered ? 1 : 0,
      notes: body.notes,
    })
    .returning();

  return vehicle;
});
