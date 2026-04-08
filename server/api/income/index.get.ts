import { incomeCategories, incomePaymentStatuses, tables } from "~~/server/db/schema";
import { eq, like, desc, gte, lte, and } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const page = parseInt(query.page as string) || 1;
  const limit = parseInt(query.limit as string) || 20;
  const offset = (page - 1) * limit;
  const search = query.search as string | undefined;
  const category = query.category as (typeof incomeCategories)[number] | "all" | undefined;
  const paymentStatus = query.paymentStatus as
    | (typeof incomePaymentStatuses)[number]
    | "all"
    | undefined;
  const dateFrom = query.dateFrom as string | undefined;
  const dateTo = query.dateTo as string | undefined;

  try {
    const conditions: any[] = [];

    if (search) {
      conditions.push(like(tables.income.description, `%${search}%`));
    }

    if (category && category !== "all") {
      conditions.push(eq(tables.income.category, category));
    }

    if (paymentStatus && paymentStatus !== "all") {
      conditions.push(eq(tables.income.paymentStatus, paymentStatus));
    }

    if (dateFrom) {
      conditions.push(gte(tables.income.incomeDate, new Date(dateFrom)));
    }

    if (dateTo) {
      conditions.push(lte(tables.income.incomeDate, new Date(dateTo)));
    }

    const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

    const [data, countResult] = await Promise.all([
      db
        .select()
        .from(tables.income)
        .where(whereClause)
        .orderBy(desc(tables.income.incomeDate))
        .limit(limit)
        .offset(offset),
      db.select({ count: tables.income.id }).from(tables.income).where(whereClause),
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
    console.error("Failed to fetch income:", error);
    throw createError({
      statusCode: 500,
      message: "Failed to fetch income",
    });
  }
});
