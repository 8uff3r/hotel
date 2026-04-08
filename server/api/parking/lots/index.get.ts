import { tables, parkingLotStatuses } from "~~/server/db/schema";
import { eq, like, and, sql } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const page = parseInt(query.page as string) || 1;
  const limit = Math.min(parseInt(query.limit as string) || 50, 100);
  const offset = (page - 1) * limit;

  const conditions = [];

  if (query.status) {
    conditions.push(
      eq(tables.parkingLots.status, query.status as (typeof parkingLotStatuses)[number])
    );
  }

  if (query.search) {
    conditions.push(like(tables.parkingLots.name, `%${query.search}%`));
  }

  const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

  const parkingLots = await db
    .select()
    .from(tables.parkingLots)
    .where(whereClause)
    .limit(limit)
    .offset(offset);

  const countResult = await db
    .select({ count: sql<number>`count(*)` })
    .from(tables.parkingLots)
    .where(whereClause);

  const total = Number(countResult[0]?.count || 0);

  return {
    data: parkingLots,
    pagination: {
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    },
  };
});
