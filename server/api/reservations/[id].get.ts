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
    const [reservation] = await db
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
        updatedAt: tables.reservations.updatedAt,
        createdBy: tables.reservations.createdBy,
        // Guest details
        guestFirstName: tables.guests.firstName,
        guestLastName: tables.guests.lastName,
        guestEmail: tables.guests.email,
        guestPhone: tables.guests.phone,
        guestIdType: tables.guests.idType,
        guestIdNumber: tables.guests.idNumber,
        guestAddress: tables.guests.address,
        guestCity: tables.guests.city,
        guestCountry: tables.guests.country,
        // Room details
        roomNumber: tables.rooms.roomNumber,
        roomType: tables.rooms.roomType,
        roomFloor: tables.rooms.floor,
        roomCapacity: tables.rooms.capacity,
        roomBasePrice: tables.rooms.basePrice,
        roomStatus: tables.rooms.status,
        roomAmenities: tables.rooms.amenities,
        roomDescription: tables.rooms.description,
      })
      .from(tables.reservations)
      .leftJoin(tables.guests, eq(tables.reservations.guestId, tables.guests.id))
      .leftJoin(tables.rooms, eq(tables.reservations.roomId, tables.rooms.id))
      .where(eq(tables.reservations.id, Number(id)))
      .limit(1);

    if (!reservation) {
      throw createError({
        statusCode: 404,
        message: "Reservation not found",
      });
    }

    // Transform the flat result into nested objects
    const result = {
      id: reservation.id,
      hotelId: reservation.hotelId,
      guestId: reservation.guestId,
      roomId: reservation.roomId,
      checkInDate: reservation.checkInDate,
      checkOutDate: reservation.checkOutDate,
      actualCheckIn: reservation.actualCheckIn,
      actualCheckOut: reservation.actualCheckOut,
      status: reservation.status,
      totalAmount: reservation.totalAmount,
      paidAmount: reservation.paidAmount,
      paymentStatus: reservation.paymentStatus,
      specialRequests: reservation.specialRequests,
      numberOfGuests: reservation.numberOfGuests,
      createdAt: reservation.createdAt,
      updatedAt: reservation.updatedAt,
      createdBy: reservation.createdBy,
      guest: {
        id: reservation.guestId,
        firstName: reservation.guestFirstName,
        lastName: reservation.guestLastName,
        email: reservation.guestEmail,
        phone: reservation.guestPhone,
        idType: reservation.guestIdType,
        idNumber: reservation.guestIdNumber,
        address: reservation.guestAddress,
        city: reservation.guestCity,
        country: reservation.guestCountry,
      },
      room: {
        id: reservation.roomId,
        roomNumber: reservation.roomNumber,
        roomType: reservation.roomType,
        floor: reservation.roomFloor,
        capacity: reservation.roomCapacity,
        basePrice: reservation.roomBasePrice,
        status: reservation.roomStatus,
        amenities: reservation.roomAmenities,
        description: reservation.roomDescription,
      },
    };

    return {
      data: result,
    };
  } catch (error: unknown) {
    if (error && typeof error === "object" && "statusCode" in error) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "Failed to fetch reservation",
    });
  }
});
