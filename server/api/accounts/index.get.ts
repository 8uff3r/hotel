import { accountsAccountTypes, tables } from "~~/server/db/schema";
import { eq, like, desc, and, SQL } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const page = parseInt(query.page as string) || 1;
  const limit = parseInt(query.limit as string) || 20;
  const offset = (page - 1) * limit;
  const search = query.search as string | undefined;
  const accountType = query.accountType as
    | (typeof accountsAccountTypes)[number]
    | "all"
    | undefined;

  try {
    const conditions: SQL[] = [];

    if (search) {
      conditions.push(like(tables.accounts.accountName, `%${search}%`));
    }

    if (accountType && accountType !== "all") {
      conditions.push(eq(tables.accounts.accountType, accountType));
    }

    const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

    const [data, countResult] = await Promise.all([
      db
        .select()
        .from(tables.accounts)
        .where(whereClause)
        .orderBy(desc(tables.accounts.id))
        .limit(limit)
        .offset(offset),
      db.select({ count: tables.accounts.id }).from(tables.accounts).where(whereClause),
    ]);

    const total = countResult.length;

    return {
      data,
      pagination: {
        page,
        limit,
        total,
        totalPages: Math.ceil(total / limit),
      },
    };
  } catch (error) {
    console.error("Failed to fetch accounts:", error);
    throw createError({
      statusCode: 500,
      message: "Failed to fetch accounts",
    });
  }
});
