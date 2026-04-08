import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = Number(getRouterParam(event, "id"));
  const body = await readBody(event);

  const [vehicle] = await db
    .update(tables.vehicles)
    .set({
      guestId: body.guestId,
      licensePlate: body.licensePlate,
      vehicleType: body.vehicleType,
      make: body.make,
      model: body.model,
      color: body.color,
      isRegistered: body.isRegistered ? 1 : 0,
      notes: body.notes,
      updatedAt: new Date(),
    })
    .where(eq(tables.vehicles.id, id))
    .returning();

  if (!vehicle) {
    throw createError({
      statusCode: 404,
      message: "Vehicle not found",
    });
  }

  return vehicle;
});
