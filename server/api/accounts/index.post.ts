import { tables } from "~~/server/db/schema";
import { z } from "zod";
import { db } from "@nuxthub/db";

const createAccountSchema = z.object({
  accountCode: z.string().min(1, "Account code is required"),
  accountName: z.string().min(1, "Account name is required"),
  accountType: z.enum(["asset", "liability", "equity", "revenue", "expense"]),
  accountSubType: z.string().optional(),
  parentId: z.number().int().positive().optional(),
  isActive: z.boolean().default(true),
  description: z.string().optional(),
  normalBalance: z.enum(["debit", "credit"]).default("debit"),
});

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  const validation = createAccountSchema.safeParse(body);
  if (!validation.success) {
    const firstIssue = validation.error.issues[0];
    throw createError({
      statusCode: 400,
      message: firstIssue?.message || "Invalid input",
    });
  }

  const data = validation.data;

  try {
    const [account] = await db
      .insert(tables.accounts)
      .values({
        accountCode: data.accountCode,
        accountName: data.accountName,
        accountType: data.accountType,
        accountSubType: data.accountSubType,
        parentId: data.parentId,
        isActive: data.isActive,
        description: data.description,
        normalBalance: data.normalBalance,
      })
      .returning();

    return {
      data: account,
      message: "Account created successfully",
    };
  } catch (error: unknown) {
    console.error("Failed to create account:", error);
    throw createError({
      statusCode: 500,
      message: "Failed to create account",
    });
  }
});
