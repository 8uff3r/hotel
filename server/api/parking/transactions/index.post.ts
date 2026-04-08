import { tables } from "~~/server/db/schema";
import { eq, and } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  if (body.spotId) {
    const [spot] = await db
      .select()
      .from(tables.parkingSpots)
      .where(eq(tables.parkingSpots.id, body.spotId));

    if (!spot) {
      throw createError({
        statusCode: 404,
        message: "Parking spot not found",
      });
    }

    if (spot.status === "occupied") {
      throw createError({
        statusCode: 400,
        message: "Parking spot is already occupied",
      });
    }

    await db
      .update(tables.parkingSpots)
      .set({ status: "occupied", updatedAt: new Date() })
      .where(eq(tables.parkingSpots.id, body.spotId));
  }

  const [transaction] = await db
    .insert(tables.parkingTransactions)
    .values({
      lotId: body.lotId,
      spotId: body.spotId,
      guestId: body.guestId,
      reservationId: body.reservationId,
      licensePlate: body.licensePlate,
      entryTime: new Date(),
      rateApplied: body.rateApplied || "hourly",
      amountDue: 0,
      amountPaid: 0,
      status: "active",
      paymentStatus: "pending",
    })
    .returning();

  return transaction;
});
