import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = Number(getRouterParam(event, "id"));

  const [parkingSpot] = await db
    .delete(tables.parkingSpots)
    .where(eq(tables.parkingSpots.id, id))
    .returning();

  if (!parkingSpot) {
    throw createError({
      statusCode: 404,
      message: "Parking spot not found",
    });
  }

  return { success: true };
});
