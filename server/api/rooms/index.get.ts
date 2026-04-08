import { roomStatuses, roomsRoomTypes, tables } from "~~/server/db/schema";
import { eq, like, and, sql } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  // Pagination
  const page = parseInt(query.page as string) || 1;
  const limit = Math.min(parseInt(query.limit as string) || 50, 100);
  const offset = (page - 1) * limit;

  // Build conditions
  const conditions = [];

  if (query.status) {
    conditions.push(eq(tables.rooms.status, query.status as (typeof roomStatuses)[number]));
  }

  if (query.roomType) {
    conditions.push(eq(tables.rooms.roomType, query.roomType as (typeof roomsRoomTypes)[number]));
  }

  if (query.floor) {
    conditions.push(eq(tables.rooms.floor, parseInt(query.floor as string)));
  }

  if (query.search) {
    conditions.push(like(tables.rooms.roomNumber, `%${query.search}%`));
  }

  // Execute query with filters
  const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

  const rooms = await db.select().from(tables.rooms).where(whereClause).limit(limit).offset(offset);

  // Get count
  const countResult = await db
    .select({ count: sql<number>`count(*)` })
    .from(tables.rooms)
    .where(whereClause);

  const total = Number(countResult[0]?.count || 0);

  return {
    data: rooms,
    pagination: {
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    },
  };
});
