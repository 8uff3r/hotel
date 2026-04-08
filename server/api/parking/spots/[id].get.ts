import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = Number(getRouterParam(event, "id"));

  const [parkingSpot] = await db
    .select()
    .from(tables.parkingSpots)
    .where(eq(tables.parkingSpots.id, id));

  if (!parkingSpot) {
    throw createError({
      statusCode: 404,
      message: "Parking spot not found",
    });
  }

  return parkingSpot;
});
