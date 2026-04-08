import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = Number(getRouterParam(event, "id"));
  const body = await readBody(event);

  const [parkingSpot] = await db
    .update(tables.parkingSpots)
    .set({
      lotId: body.lotId,
      spotNumber: body.spotNumber,
      floor: body.floor,
      spotType: body.spotType,
      status: body.status,
      isCovered: body.isCovered,
      description: body.description,
      updatedAt: new Date(),
    })
    .where(eq(tables.parkingSpots.id, id))
    .returning();

  if (!parkingSpot) {
    throw createError({
      statusCode: 404,
      message: "Parking spot not found",
    });
  }

  return parkingSpot;
});
