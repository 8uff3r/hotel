import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = Number(getRouterParam(event, "id"));
  const body = await readBody(event);

  const [parkingLot] = await db
    .update(tables.parkingLots)
    .set({
      name: body.name,
      location: body.location,
      totalSpots: body.totalSpots,
      hourlyRate: body.hourlyRate,
      dailyRate: body.dailyRate,
      status: body.status,
      description: body.description,
      updatedAt: new Date(),
    })
    .where(eq(tables.parkingLots.id, id))
    .returning();

  if (!parkingLot) {
    throw createError({
      statusCode: 404,
      message: "Parking lot not found",
    });
  }

  return parkingLot;
});
