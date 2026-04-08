import { reservationStatuses, tables, reservationPaymentStatuses } from "~~/server/db/schema";
import { eq, gte, lte, sql, and } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  // Pagination
  const page = parseInt(query.page as string) || 1;
  const limit = Math.min(parseInt(query.limit as string) || 50, 100);
  const offset = (page - 1) * limit;

  // Build conditions
  const conditions = [];

  if (query.status) {
    conditions.push(
      eq(tables.reservations.status, query.status as (typeof reservationStatuses)[number])
    );
  }

  if (query.paymentStatus) {
    conditions.push(
      eq(
        tables.reservations.paymentStatus,
        query.paymentStatus as (typeof reservationPaymentStatuses)[number]
      )
    );
  }

  if (query.checkInFrom) {
    conditions.push(gte(tables.reservations.checkInDate, new Date(query.checkInFrom as string)));
  }

  if (query.checkInTo) {
    conditions.push(lte(tables.reservations.checkInDate, new Date(query.checkInTo as string)));
  }

  const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

  // Get reservations with room info
  const reservations = await db
    .select({
      id: tables.reservations.id,
      hotelId: tables.reservations.hotelId,
      guestId: tables.reservations.guestId,
      roomId: tables.reservations.roomId,
      checkInDate: tables.reservations.checkInDate,
      checkOutDate: tables.reservations.checkOutDate,
      actualCheckIn: tables.reservations.actualCheckIn,
      actualCheckOut: tables.reservations.actualCheckOut,
      status: tables.reservations.status,
      totalAmount: tables.reservations.totalAmount,
      paidAmount: tables.reservations.paidAmount,
      paymentStatus: tables.reservations.paymentStatus,
      specialRequests: tables.reservations.specialRequests,
      numberOfGuests: tables.reservations.numberOfGuests,
      createdAt: tables.reservations.createdAt,
      roomNumber: tables.rooms.roomNumber,
      roomType: tables.rooms.roomType,
      guest: {
        id: tables.guests.id,
        firstName: tables.guests.firstName,
        lastName: tables.guests.lastName,
        email: tables.guests.email,
        phone: tables.guests.phone,
      },
    })
    .from(tables.reservations)
    .leftJoin(tables.rooms, eq(tables.reservations.roomId, tables.rooms.id))
    .leftJoin(tables.guests, eq(tables.reservations.guestId, tables.guests.id))
    .where(whereClause)
    .limit(limit)
    .offset(offset);

  // Get count
  const countResult = await db
    .select({ count: sql<number>`count(*)` })
    .from(tables.reservations)
    .leftJoin(tables.rooms, eq(tables.reservations.roomId, tables.rooms.id))
    .leftJoin(tables.guests, eq(tables.reservations.guestId, tables.guests.id))
    .where(whereClause);

  const total = Number(countResult[0]?.count || 0);

  return {
    data: reservations,
    pagination: {
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    },
  };
});
