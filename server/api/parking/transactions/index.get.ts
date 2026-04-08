import { tables, parkingTransactionStatuses, paymentStatuses } from "~~/server/db/schema";
import { eq, like, and, sql, gte, lte } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const page = parseInt(query.page as string) || 1;
  const limit = Math.min(parseInt(query.limit as string) || 50, 100);
  const offset = (page - 1) * limit;

  const conditions = [];

  if (query.lotId) {
    conditions.push(eq(tables.parkingTransactions.lotId, parseInt(query.lotId as string)));
  }

  if (query.guestId) {
    conditions.push(eq(tables.parkingTransactions.guestId, parseInt(query.guestId as string)));
  }

  if (query.status) {
    conditions.push(
      eq(
        tables.parkingTransactions.status,
        query.status as (typeof parkingTransactionStatuses)[number]
      )
    );
  }

  if (query.paymentStatus) {
    conditions.push(
      eq(
        tables.parkingTransactions.paymentStatus,
        query.paymentStatus as (typeof paymentStatuses)[number]
      )
    );
  }

  if (query.search) {
    conditions.push(like(tables.parkingTransactions.licensePlate, `%${query.search}%`));
  }

  if (query.startDate) {
    conditions.push(gte(tables.parkingTransactions.entryTime, new Date(query.startDate as string)));
  }

  if (query.endDate) {
    conditions.push(lte(tables.parkingTransactions.entryTime, new Date(query.endDate as string)));
  }

  const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

  const transactions = await db
    .select()
    .from(tables.parkingTransactions)
    .where(whereClause)
    .limit(limit)
    .offset(offset);

  const countResult = await db
    .select({ count: sql<number>`count(*)` })
    .from(tables.parkingTransactions)
    .where(whereClause);

  const total = Number(countResult[0]?.count || 0);

  return {
    data: transactions,
    pagination: {
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    },
  };
});
