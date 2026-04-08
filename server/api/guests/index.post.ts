import { tables } from "~~/server/db/schema";
import { z } from "zod";
import { db } from "@nuxthub/db";

const createGuestSchema = z.object({
  firstName: z.string().min(1, "First name is required"),
  lastName: z.string().min(1, "Last name is required"),
  email: z.string().email("Invalid email address").optional().or(z.literal("")),
  phone: z.string().optional(),
  idType: z.enum(["passport", "national_id", "driver_license", "other"]).optional(),
  idNumber: z.string().optional(),
  address: z.string().optional(),
  city: z.string().optional(),
  country: z.string().optional(),
  notes: z.string().optional(),
});

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  const validation = createGuestSchema.safeParse(body);
  if (!validation.success) {
    const firstIssue = validation.error.issues[0];
    throw createError({
      statusCode: 400,
      message: firstIssue?.message || "Invalid input",
    });
  }

  const data = validation.data;

  try {
    const [guest] = await db.insert(tables.guests).values(data).returning();

    return {
      data: guest,
      message: "Guest created successfully",
    };
  } catch (error: unknown) {
    if (error && typeof error === "object" && "statusCode" in error) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "Failed to create guest",
    });
  }
});
