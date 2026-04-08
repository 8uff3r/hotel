import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = Number(getRouterParam(event, "id"));

  const [transaction] = await db
    .select()
    .from(tables.parkingTransactions)
    .where(eq(tables.parkingTransactions.id, id));

  if (!transaction) {
    throw createError({
      statusCode: 404,
      message: "Transaction not found",
    });
  }

  return transaction;
});
