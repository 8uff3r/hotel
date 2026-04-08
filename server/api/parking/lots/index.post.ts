import { tables } from "~~/server/db/schema";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  const [parkingLot] = await db
    .insert(tables.parkingLots)
    .values({
      hotelId: body.hotelId,
      name: body.name,
      location: body.location,
      totalSpots: body.totalSpots || 0,
      hourlyRate: body.hourlyRate || "0",
      dailyRate: body.dailyRate || "0",
      status: body.status || "active",
      description: body.description,
    })
    .returning();

  return parkingLot;
});
