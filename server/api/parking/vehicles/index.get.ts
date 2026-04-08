import { tables, vehicleTypes } from "~~/server/db/schema";
import { eq, like, and, sql } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const page = parseInt(query.page as string) || 1;
  const limit = Math.min(parseInt(query.limit as string) || 50, 100);
  const offset = (page - 1) * limit;

  const conditions = [];

  if (query.guestId) {
    conditions.push(eq(tables.vehicles.guestId, parseInt(query.guestId as string)));
  }

  if (query.vehicleType) {
    conditions.push(
      eq(tables.vehicles.vehicleType, query.vehicleType as (typeof vehicleTypes)[number])
    );
  }

  if (query.search) {
    conditions.push(like(tables.vehicles.licensePlate, `%${query.search}%`));
  }

  const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

  const vehicles = await db
    .select()
    .from(tables.vehicles)
    .where(whereClause)
    .limit(limit)
    .offset(offset);

  const countResult = await db
    .select({ count: sql<number>`count(*)` })
    .from(tables.vehicles)
    .where(whereClause);

  const total = Number(countResult[0]?.count || 0);

  return {
    data: vehicles,
    pagination: {
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    },
  };
});
