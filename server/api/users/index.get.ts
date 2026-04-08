import { tables } from "~~/server/db/schema";
import { eq, sql } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const page = parseInt(query.page as string) || 1;
  const limit = Math.min(parseInt(query.limit as string) || 50, 100);
  const offset = (page - 1) * limit;

  const conditions = [];

  if (query.search) {
    conditions.push(eq(tables.users.email, `%${query.search}%` as any));
  }

  const whereClause = conditions.length > 0 ? undefined : undefined;

  const users = await db.select().from(tables.users).limit(limit).offset(offset);

  // Get roles for each user
  const usersWithRoles = await Promise.all(
    users.map(async (user) => {
      const userRoles = await db
        .select()
        .from(tables.userRoles)
        .where(eq(tables.userRoles.userId, user.id));

      const roles = userRoles.length > 0 ? userRoles.map((r) => r.role) : [user.role || "staff"];

      const { passwordHash, ...userWithoutPassword } = user;
      return { ...userWithoutPassword, roles };
    })
  );

  const countResult = await db.select({ count: sql<number>`count(*)` }).from(tables.users);

  const total = Number(countResult[0]?.count || 0);

  return {
    data: usersWithRoles,
    pagination: {
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    },
  };
});
