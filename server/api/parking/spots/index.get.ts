import { tables, parkingSpotStatuses, spotTypes } from "~~/server/db/schema";
import { eq, like, and, sql } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const page = parseInt(query.page as string) || 1;
  const limit = Math.min(parseInt(query.limit as string) || 50, 100);
  const offset = (page - 1) * limit;

  const conditions = [];

  if (query.lotId) {
    conditions.push(eq(tables.parkingSpots.lotId, parseInt(query.lotId as string)));
  }

  if (query.status) {
    conditions.push(
      eq(tables.parkingSpots.status, query.status as (typeof parkingSpotStatuses)[number])
    );
  }

  if (query.spotType) {
    conditions.push(eq(tables.parkingSpots.spotType, query.spotType as (typeof spotTypes)[number]));
  }

  if (query.search) {
    conditions.push(like(tables.parkingSpots.spotNumber, `%${query.search}%`));
  }

  const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

  const parkingSpots = await db
    .select()
    .from(tables.parkingSpots)
    .where(whereClause)
    .limit(limit)
    .offset(offset);

  const countResult = await db
    .select({ count: sql<number>`count(*)` })
    .from(tables.parkingSpots)
    .where(whereClause);

  const total = Number(countResult[0]?.count || 0);

  return {
    data: parkingSpots,
    pagination: {
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    },
  };
});
