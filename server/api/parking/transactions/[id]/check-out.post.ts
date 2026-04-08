import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = Number(getRouterParam(event, "id"));
  const body = await readBody(event);

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

  if (transaction.status !== "active") {
    throw createError({
      statusCode: 400,
      message: "Transaction is not active",
    });
  }

  const exitTime = new Date();
  const entryTime = new Date(transaction.entryTime);
  const hoursParked = (exitTime.getTime() - entryTime.getTime()) / (1000 * 60 * 60);

  const [lot] = await db
    .select()
    .from(tables.parkingLots)
    .where(eq(tables.parkingLots.id, transaction.lotId!));

  let amountDue = 0;
  if (lot) {
    const rate =
      body.rateType === "daily"
        ? parseFloat(lot.dailyRate)
        : parseFloat(lot.hourlyRate) * hoursParked;
    amountDue = Math.ceil(rate * 100) / 100;
  }

  const [updatedTransaction] = await db
    .update(tables.parkingTransactions)
    .set({
      exitTime,
      hoursParked,
      amountDue,
      amountPaid: body.amountPaid || amountDue,
      paymentStatus: body.amountPaid && body.amountPaid >= amountDue ? "paid" : "pending",
      paymentMethod: body.paymentMethod,
      status: "completed",
      notes: body.notes,
      updatedAt: new Date(),
    })
    .where(eq(tables.parkingTransactions.id, id))
    .returning();

  if (transaction.spotId) {
    await db
      .update(tables.parkingSpots)
      .set({ status: "available", updatedAt: new Date() })
      .where(eq(tables.parkingSpots.id, transaction.spotId));
  }

  return updatedTransaction;
});
