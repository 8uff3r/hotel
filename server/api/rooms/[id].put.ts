import { tables } from "~~/server/db/schema";
import { z } from "zod";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

const updateRoomSchema = z.object({
  roomNumber: z.string().min(1).optional(),
  roomType: z.enum(["single", "double", "suite", "deluxe"]).optional(),
  floor: z.number().int().positive().optional(),
  capacity: z.number().int().positive().optional(),
  basePrice: z.number().positive().optional(),
  status: z.enum(["available", "occupied", "maintenance", "out_of_order"]).optional(),
  amenities: z.array(z.string()).optional(),
  description: z.string().optional(),
  images: z.array(z.string()).optional(),
});

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

  const body = await readBody(event);

  // Validate input
  const validation = updateRoomSchema.safeParse(body);
  if (!validation.success) {
    const firstIssue = validation.error.issues[0];
    throw createError({
      statusCode: 400,
      message: firstIssue?.message || "Invalid input",
    });
  }

  const data = validation.data;

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

  // Check if room number is being changed and if it's already taken
  if (data.roomNumber && data.roomNumber !== existingRoom.roomNumber) {
    const [duplicate] = await db
      .select()
      .from(tables.rooms)
      .where(eq(tables.rooms.roomNumber, data.roomNumber))
      .limit(1);

    if (duplicate) {
      throw createError({
        statusCode: 409,
        message: "Room number already exists",
      });
    }
  }

  // Convert amenities and images to JSON strings if provided
  const updateData: Record<string, unknown> = { ...data };

  if (data.amenities !== undefined) {
    updateData.amenities = JSON.stringify(data.amenities);
  }

  if (data.images !== undefined) {
    updateData.images = JSON.stringify(data.images);
  }

  // Add updatedAt timestamp
  updateData.updatedAt = new Date();

  const [room] = await db
    .update(tables.rooms)
    .set(updateData)
    .where(eq(tables.rooms.id, roomId))
    .returning();

  return {
    data: {
      ...room,
      amenities: room?.amenities ? JSON.parse(room.amenities) : [],
      images: room?.images ? JSON.parse(room.images) : [],
    },
    message: "Room updated successfully",
  };
});
