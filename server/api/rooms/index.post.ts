import { tables } from "~~/server/db/schema";
import { z } from "zod";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

const createRoomSchema = z.object({
  roomNumber: z.string().min(1, "Room number is required"),
  roomType: z.enum(["single", "double", "suite", "deluxe"]).default("single"),
  floor: z.number().int().positive().optional(),
  capacity: z.number().int().positive().default(2),
  basePrice: z.number().positive().default(0),
  status: z.enum(["available", "occupied", "maintenance", "out_of_order"]).default("available"),
  amenities: z.array(z.string()).optional(),
  description: z.string().optional(),
  images: z.array(z.string()).optional(),
});

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  // Validate input
  const validation = createRoomSchema.safeParse(body);
  if (!validation.success) {
    const firstIssue = validation.error.issues[0];
    throw createError({
      statusCode: 400,
      message: firstIssue?.message || "Invalid input",
    });
  }

  const data = validation.data;

  try {
    // Check if room number already exists
    const existingRoom = await db
      .select()
      .from(tables.rooms)
      .where(eq(tables.rooms.roomNumber, data.roomNumber))
      .limit(1);

    if (existingRoom.length > 0) {
      throw createError({
        statusCode: 409,
        message: "Room number already exists",
      });
    }

    // Convert amenities and images to JSON strings
    const insertData = {
      ...data,
      amenities: data.amenities ? JSON.stringify(data.amenities) : null,
      images: data.images ? JSON.stringify(data.images) : null,
    };

    const [room] = await db.insert(tables.rooms).values(insertData).returning();

    return {
      data: room,
      message: "Room created successfully",
    };
  } catch (error: unknown) {
    if (error && typeof error === "object" && "statusCode" in error) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "Failed to create room",
    });
  }
});
