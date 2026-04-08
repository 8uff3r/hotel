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

  const [room] = await db.select().from(tables.rooms).where(eq(tables.rooms.id, roomId)).limit(1);

  if (!room) {
    throw createError({
      statusCode: 404,
      message: "Room not found",
    });
  }

  // Parse JSON fields
  return {
    data: {
      ...room,
      amenities: room.amenities ? JSON.parse(room.amenities) : [],
      images: room.images ? JSON.parse(room.images) : [],
    },
  };
});
