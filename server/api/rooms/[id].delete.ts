import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = getRouterParam(event, "id");

  if (!id) {
    throw createError({
      statusCode: 400,
      message: "Room ID is required",
    });
  }

  const roomId = parseInt(id);

  if (isNaN(roomId)) {
    throw createError({
      statusCode: 400,
      message: "Invalid room ID",
    });
  }

  // Check if room exists
  const [existingRoom] = await db
    .select()
    .from(tables.rooms)
    .where(eq(tables.rooms.id, roomId))
    .limit(1);

  if (!existingRoom) {
    throw createError({
      statusCode: 404,
      message: "Room not found",
    });
  }

  // Check if room has active reservations
  const activeReservations = await db
    .select()
    .from(tables.reservations)
    .where(eq(tables.reservations.roomId, roomId))
    .limit(1);

  if (activeReservations.length > 0) {
    throw createError({
      statusCode: 409,
      message: "Cannot delete room with active reservations",
    });
  }

  await db.delete(tables.rooms).where(eq(tables.rooms.id, roomId));

  return {
    message: "Room deleted successfully",
  };
});
