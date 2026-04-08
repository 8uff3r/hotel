import { tables } from "~~/server/db/schema";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  const [parkingSpot] = await db
    .insert(tables.parkingSpots)
    .values({
      lotId: body.lotId,
      spotNumber: body.spotNumber,
      floor: body.floor,
      spotType: body.spotType || "standard",
      status: body.status || "available",
      isCovered: body.isCovered || false,
      description: body.description,
    })
    .returning();

  return parkingSpot;
});
