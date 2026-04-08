import { tables } from "~~/server/db/schema";
import { like, sql, and, or } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  // Pagination
  const page = parseInt(query.page as string) || 1;
  const limit = Math.min(parseInt(query.limit as string) || 50, 100);
  const offset = (page - 1) * limit;

  // Build conditions
  const conditions = [];

  if (query.search) {
    const searchTerm = `%${query.search}%`;
    conditions.push(
      or(
        like(tables.guests.firstName, searchTerm),
        like(tables.guests.lastName, searchTerm),
        like(tables.guests.email, searchTerm),
        like(tables.guests.phone, searchTerm)
      )
    );
  }

  if (query.country) {
    conditions.push(like(tables.guests.country, `%${query.country}%`));
  }

  const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

  // Execute query
  const guests = await db
    .select()
    .from(tables.guests)
    .where(whereClause)
    .limit(limit)
    .offset(offset);

  // Get count
  const countResult = await db
    .select({ count: sql<number>`count(*)` })
    .from(tables.guests)
    .where(whereClause);

  const total = Number(countResult[0]?.count || 0);

  return {
    data: guests,
    pagination: {
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    },
  };
});
