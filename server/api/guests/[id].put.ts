import { tables } from "~~/server/db/schema";
import { z } from "zod";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

const updateGuestSchema = z.object({
  firstName: z.string().min(1).optional(),
  lastName: z.string().min(1).optional(),
  email: z.string().email().optional().or(z.literal("")),
  phone: z.string().optional().or(z.literal("")),
  idType: z.enum(["passport", "national_id", "driver_license", "other"]).optional(),
  idNumber: z.string().optional().or(z.literal("")),
  address: z.string().optional().or(z.literal("")),
  city: z.string().optional().or(z.literal("")),
  country: z.string().optional().or(z.literal("")),
  notes: z.string().optional().or(z.literal("")),
});

export default defineEventHandler(async (event) => {
  const id = getRouterParam(event, "id");

  if (!id || isNaN(Number(id))) {
    throw createError({
      statusCode: 400,
      message: "Invalid guest ID",
    });
  }

  const body = await readBody(event);

  const validation = updateGuestSchema.safeParse(body);
  if (!validation.success) {
    const firstIssue = validation.error.issues[0];
    throw createError({
      statusCode: 400,
      message: firstIssue?.message || "Invalid input",
    });
  }

  const data = validation.data;

  try {
    const [existingGuest] = await db
      .select()
      .from(tables.guests)
      .where(eq(tables.guests.id, Number(id)))
      .limit(1);

    if (!existingGuest) {
      throw createError({
        statusCode: 404,
        message: "Guest not found",
      });
    }

    const updateData: Record<string, unknown> = {
      updatedAt: new Date(),
    };

    if (data.firstName !== undefined) updateData.firstName = data.firstName;
    if (data.lastName !== undefined) updateData.lastName = data.lastName;
    if (data.email !== undefined) updateData.email = data.email || null;
    if (data.phone !== undefined) updateData.phone = data.phone || null;
    if (data.idType !== undefined) updateData.idType = data.idType;
    if (data.idNumber !== undefined) updateData.idNumber = data.idNumber || null;
    if (data.address !== undefined) updateData.address = data.address || null;
    if (data.city !== undefined) updateData.city = data.city || null;
    if (data.country !== undefined) updateData.country = data.country || null;
    if (data.notes !== undefined) updateData.notes = data.notes || null;

    const [updatedGuest] = await db
      .update(tables.guests)
      .set(updateData)
      .where(eq(tables.guests.id, Number(id)))
      .returning();

    return {
      data: updatedGuest,
      message: "Guest updated successfully",
    };
  } catch (error: unknown) {
    if (error && typeof error === "object" && "statusCode" in error) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "Failed to update guest",
    });
  }
});
