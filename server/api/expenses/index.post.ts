import {
  expenseCategories,
  paymentMethods,
  expensesPaymentStatuses,
  tables,
} from "~~/server/db/schema";
import { z } from "zod";
import { db } from "@nuxthub/db";

const createExpenseSchema = z.object({
  expenseDate: z.string().transform((val) => new Date(val)),
  description: z.string().min(1, "Description is required"),
  amount: z.number().positive("Amount must be positive"),
  category: z.enum(expenseCategories),
  vendor: z.string().optional(),
  reference: z.string().optional(),
  paymentMethod: z.enum(paymentMethods).optional(),
  paymentStatus: z.enum(expensesPaymentStatuses).default("pending"),
  accountId: z.number().int().positive().optional(),
  receiptNumber: z.string().optional(),
  notes: z.string().optional(),
});

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  const validation = createExpenseSchema.safeParse(body);
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
    const [expense] = await db
      .insert(tables.expenses)
      .values({
        expenseDate: data.expenseDate,
        description: data.description,
        amount: data.amount as any, //FIXME: it should be inferred properly but isn't
        category: data.category,
        vendor: data.vendor,
        reference: data.reference,
        paymentMethod: data.paymentMethod,
        paymentStatus: data.paymentStatus,
        accountId: data.accountId,
        receiptNumber: data.receiptNumber,
        notes: data.notes,
        createdBy,
      })
      .returning();

    return {
      data: expense,
      message: "Expense created successfully",
    };
  } catch (error: unknown) {
    console.error("Failed to create expense:", error);
    throw createError({
      statusCode: 500,
      message: "Failed to create expense",
    });
  }
});
