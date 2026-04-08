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

    // Check if reservation can be checked in
    if (reservation.status !== "confirmed") {
      throw createError({
        statusCode: 400,
        message: `Cannot check in reservation with status: ${reservation.status}`,
      });
    }

    // Update reservation status
    const [updatedReservation] = await db
      .update(tables.reservations)
      .set({
        status: "checked_in",
        actualCheckIn: new Date(),
        updatedAt: new Date(),
      })
      .where(eq(tables.reservations.id, Number(id)))
      .returning();

    return {
      data: updatedReservation,
      message: "Guest checked in successfully",
    };
  } catch (error: unknown) {
    if (error && typeof error === "object" && "statusCode" in error) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "Failed to check in reservation",
    });
  }
});
