import { tables } from "~~/server/db/schema";
import { z } from "zod";
import { db } from "@nuxthub/db";

const createIncomeSchema = z.object({
  incomeDate: z.string().transform((val) => new Date(val)),
  description: z.string().min(1, "Description is required"),
  amount: z.number().positive("Amount must be positive"),
  category: z.enum(["room_revenue", "food_beverage", "laundry", "spa", "meeting_rooms", "other"]),
  source: z.string().optional(),
  reference: z.string().optional(),
  paymentMethod: z.enum(["cash", "credit_card", "debit_card", "bank_transfer", "other"]).optional(),
  paymentStatus: z.enum(["pending", "received", "refunded"]).default("pending"),
  accountId: z.number().int().positive().optional(),
  reservationId: z.number().int().positive().optional(),
  notes: z.string().optional(),
});

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  const validation = createIncomeSchema.safeParse(body);
  if (!validation.success) {
    const firstIssue = validation.error.issues[0];
    throw createError({
      statusCode: 400,
      message: firstIssue?.message || "Invalid input",
    });
  }

  const data = validation.data;

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
    const [income] = await db
      .insert(tables.income)
      .values({
        incomeDate: data.incomeDate,
        description: data.description,
        amount: data.amount as any,
        category: data.category,
        source: data.source,
        reference: data.reference,
        paymentMethod: data.paymentMethod,
        paymentStatus: data.paymentStatus,
        accountId: data.accountId,
        reservationId: data.reservationId,
        notes: data.notes,
        createdBy,
      })
      .returning();

    return {
      data: income,
      message: "Income recorded successfully",
    };
  } catch (error: unknown) {
    console.error("Failed to create income:", error);
    throw createError({
      statusCode: 500,
      message: "Failed to create income",
    });
  }
});
