import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = getRouterParam(event, "id");

  if (!id || isNaN(Number(id))) {
    throw createError({
      statusCode: 400,
      message: "Invalid reservation ID",
    });
  }

  try {
    // Check if reservation exists
    const [reservation] = await db
      .select()
      .from(tables.reservations)
      .where(eq(tables.reservations.id, Number(id)))
      .limit(1);

    if (!reservation) {
      throw createError({
        statusCode: 404,
        message: "Reservation not found",
      });
    }

    // Check if reservation can be checked out
    if (reservation.status !== "checked_in") {
      throw createError({
        statusCode: 400,
        message: `Cannot check out reservation with status: ${reservation.status}`,
      });
    }

    // Update reservation status
    const [updatedReservation] = await db
      .update(tables.reservations)
      .set({
        status: "checked_out",
        actualCheckOut: new Date(),
        updatedAt: new Date(),
      })
      .where(eq(tables.reservations.id, Number(id)))
      .returning();

    // Make room available again
    await db
      .update(tables.rooms)
      .set({ status: "available", updatedAt: new Date() })
      .where(eq(tables.rooms.id, reservation.roomId));

    return {
      data: updatedReservation,
      message: "Guest checked out successfully",
    };
  } catch (error: unknown) {
    if (error && typeof error === "object" && "statusCode" in error) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "Failed to check out reservation",
    });
  }
});
