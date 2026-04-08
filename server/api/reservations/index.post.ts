import { tables } from "~~/server/db/schema";
import { z } from "zod";
import { eq, gte, lte } from "drizzle-orm";
import { db } from "@nuxthub/db";

const createReservationSchema = z.object({
  guestId: z.number().int().positive("Guest is required"),
  roomId: z.number().int().positive("Room is required"),
  checkInDate: z.string().transform((val) => new Date(val)),
  checkOutDate: z.string().transform((val) => new Date(val)),
  numberOfGuests: z.number().int().positive().default(1),
  specialRequests: z.string().optional(),
  totalAmount: z.number().positive().default(0),
  paidAmount: z.number().min(0).default(0),
});

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  // Validate input
  const validation = createReservationSchema.safeParse(body);
  if (!validation.success) {
    const firstIssue = validation.error.issues[0];
    throw createError({
      statusCode: 400,
      message: firstIssue?.message || "Invalid input",
    });
  }

  const data = validation.data;

  // Validate dates
  if (data.checkInDate >= data.checkOutDate) {
    throw createError({
      statusCode: 400,
      message: "Check-out date must be after check-in date",
    });
  }

  // Get session for createdBy
  const sessionCookie = getCookie(event, "auth_session");
  let createdBy: number | undefined;
  if (sessionCookie) {
    try {
      const session = JSON.parse(sessionCookie);
      createdBy = session.userId;
    } catch {
      // Invalid session, continue without createdBy
    }
  }

  try {
    // Check if room exists and is available
    const [room] = await db
      .select()
      .from(tables.rooms)
      .where(eq(tables.rooms.id, data.roomId))
      .limit(1);

    if (!room) {
      throw createError({
        statusCode: 404,
        message: "Room not found",
      });
    }

    if (room.status !== "available") {
      throw createError({
        statusCode: 409,
        message: `Room is not available (current status: ${room.status})`,
      });
    }

    // Check for overlapping reservations
    const overlapping = await db
      .select()
      .from(tables.reservations)
      .where(eq(tables.reservations.roomId, data.roomId))
      .limit(10);

    const hasOverlap = overlapping.some(
      (r) =>
        r.status !== "cancelled" &&
        r.status !== "checked_out" &&
        r.status !== "no_show" &&
        data.checkInDate < new Date(r.checkOutDate) &&
        data.checkOutDate > new Date(r.checkInDate)
    );

    if (hasOverlap) {
      throw createError({
        statusCode: 409,
        message: "Room is already reserved for the selected dates",
      });
    }

    // Calculate total amount if not provided
    let totalAmount = data.totalAmount;
    if (totalAmount === 0) {
      const nights = Math.ceil(
        (data.checkOutDate.getTime() - data.checkInDate.getTime()) / (1000 * 60 * 60 * 24)
      );
      totalAmount = room.basePrice * nights;
    }

    // Determine payment status
    let paymentStatus: "pending" | "partial" | "paid" = "pending";
    if (data.paidAmount >= totalAmount) {
      paymentStatus = "paid";
    } else if (data.paidAmount > 0) {
      paymentStatus = "partial";
    }

    const [reservation] = await db
      .insert(tables.reservations)
      .values({
        guestId: data.guestId,
        roomId: data.roomId,
        checkInDate: data.checkInDate,
        checkOutDate: data.checkOutDate,
        numberOfGuests: data.numberOfGuests,
        specialRequests: data.specialRequests,
        totalAmount,
        paidAmount: data.paidAmount,
        paymentStatus,
        status: "confirmed",
        createdBy,
      })
      .returning();

    // Update room status to occupied
    await db
      .update(tables.rooms)
      .set({ status: "occupied", updatedAt: new Date() })
      .where(eq(tables.rooms.id, data.roomId));

    return {
      data: reservation,
      message: "Reservation created successfully",
    };
  } catch (error: unknown) {
    if (error && typeof error === "object" && "statusCode" in error) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "Failed to create reservation",
    });
  }
});
