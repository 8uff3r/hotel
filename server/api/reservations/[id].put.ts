import { tables } from "~~/server/db/schema";
import { z } from "zod";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

const updateReservationSchema = z.object({
  guestId: z.number().int().positive().optional(),
  roomId: z.number().int().positive().optional(),
  checkInDate: z
    .string()
    .transform((val) => new Date(val))
    .optional(),
  checkOutDate: z
    .string()
    .transform((val) => new Date(val))
    .optional(),
  numberOfGuests: z.number().int().positive().optional(),
  specialRequests: z.string().optional(),
  totalAmount: z.number().min(0).optional(),
  paidAmount: z.number().min(0).optional(),
  paymentStatus: z.enum(["pending", "partial", "paid", "refunded"]).optional(),
  status: z
    .enum(["pending", "confirmed", "checked_in", "checked_out", "cancelled", "no_show"])
    .optional(),
});

export default defineEventHandler(async (event) => {
  const id = getRouterParam(event, "id");

  if (!id || isNaN(Number(id))) {
    throw createError({
      statusCode: 400,
      message: "Invalid reservation ID",
    });
  }

  const body = await readBody(event);

  // Validate input
  const validation = updateReservationSchema.safeParse(body);
  if (!validation.success) {
    const firstIssue = validation.error.issues[0];
    throw createError({
      statusCode: 400,
      message: firstIssue?.message || "Invalid input",
    });
  }

  const data = validation.data;

  try {
    // Check if reservation exists
    const [existingReservation] = await db
      .select()
      .from(tables.reservations)
      .where(eq(tables.reservations.id, Number(id)))
      .limit(1);

    if (!existingReservation) {
      throw createError({
        statusCode: 404,
        message: "Reservation not found",
      });
    }

    // If changing room, check if new room exists and is available
    if (data.roomId && data.roomId !== existingReservation.roomId) {
      const [newRoom] = await db
        .select()
        .from(tables.rooms)
        .where(eq(tables.rooms.id, data.roomId))
        .limit(1);

      if (!newRoom) {
        throw createError({
          statusCode: 404,
          message: "Room not found",
        });
      }

      if (newRoom.status !== "available" && newRoom.status !== "occupied") {
        throw createError({
          statusCode: 409,
          message: `Room is not available (current status: ${newRoom.status})`,
        });
      }
    }

    // If changing dates, validate them
    if (data.checkInDate && data.checkOutDate) {
      if (data.checkInDate >= data.checkOutDate) {
        throw createError({
          statusCode: 400,
          message: "Check-out date must be after check-in date",
        });
      }
    }

    // Build update object
    const updateData: Record<string, unknown> = {
      updatedAt: new Date(),
    };

    if (data.guestId !== undefined) updateData.guestId = data.guestId;
    if (data.roomId !== undefined) updateData.roomId = data.roomId;
    if (data.checkInDate !== undefined) updateData.checkInDate = data.checkInDate;
    if (data.checkOutDate !== undefined) updateData.checkOutDate = data.checkOutDate;
    if (data.numberOfGuests !== undefined) updateData.numberOfGuests = data.numberOfGuests;
    if (data.specialRequests !== undefined) updateData.specialRequests = data.specialRequests;
    if (data.totalAmount !== undefined) updateData.totalAmount = data.totalAmount;
    if (data.paidAmount !== undefined) updateData.paidAmount = data.paidAmount;
    if (data.paymentStatus !== undefined) updateData.paymentStatus = data.paymentStatus;
    if (data.status !== undefined) updateData.status = data.status;

    const [updatedReservation] = await db
      .update(tables.reservations)
      .set(updateData)
      .where(eq(tables.reservations.id, Number(id)))
      .returning();

    // If status changed to cancelled, make room available again
    if (data.status === "cancelled" && existingReservation.status !== "cancelled") {
      await db
        .update(tables.rooms)
        .set({ status: "available", updatedAt: new Date() })
        .where(eq(tables.rooms.id, existingReservation.roomId));
    }

    // If status changed to checked_out, make room available
    if (data.status === "checked_out" && existingReservation.status !== "checked_out") {
      await db
        .update(tables.rooms)
        .set({ status: "available", updatedAt: new Date() })
        .where(eq(tables.rooms.id, existingReservation.roomId));
    }

    return {
      data: updatedReservation,
      message: "Reservation updated successfully",
    };
  } catch (error: unknown) {
    if (error && typeof error === "object" && "statusCode" in error) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "Failed to update reservation",
    });
  }
});
