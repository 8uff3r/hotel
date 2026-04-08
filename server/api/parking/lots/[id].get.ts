import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = Number(getRouterParam(event, "id"));

  const [parkingLot] = await db
    .select()
    .from(tables.parkingLots)
    .where(eq(tables.parkingLots.id, id));

  if (!parkingLot) {
    throw createError({
      statusCode: 404,
      message: "Parking lot not found",
    });
  }

  return parkingLot;
});
