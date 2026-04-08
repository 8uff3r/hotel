import { expenseCategories, expensesPaymentStatuses, tables } from "~~/server/db/schema";
import { eq, like, desc, gte, lte, and, SQL } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const page = parseInt(query.page as string) || 1;
  const limit = parseInt(query.limit as string) || 20;
  const offset = (page - 1) * limit;
  const search = query.search as string | undefined;

  const category = query.category as (typeof expenseCategories)[number] | "all" | undefined;
  const paymentStatus = query.paymentStatus as
    | (typeof expensesPaymentStatuses)[number]
    | "all"
    | undefined;
  const dateFrom = query.dateFrom as string | undefined;
  const dateTo = query.dateTo as string | undefined;

  try {
    const conditions: SQL[] = [];

    if (search) {
      conditions.push(like(tables.expenses.description, `%${search}%`));
    }

    if (category && category !== "all") {
      conditions.push(eq(tables.expenses.category, category));
    }

    if (paymentStatus && paymentStatus !== "all") {
      conditions.push(eq(tables.expenses.paymentStatus, paymentStatus));
    }

    if (dateFrom) {
      conditions.push(gte(tables.expenses.expenseDate, new Date(dateFrom)));
    }

    if (dateTo) {
      conditions.push(lte(tables.expenses.expenseDate, new Date(dateTo)));
    }

    const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

    const [data, countResult] = await Promise.all([
      db
        .select()
        .from(tables.expenses)
        .where(whereClause)
        .orderBy(desc(tables.expenses.expenseDate))
        .limit(limit)
        .offset(offset),
      db.select({ count: tables.expenses.id }).from(tables.expenses).where(whereClause),
    ])

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
    console.error("Failed to fetch expenses:", error);
    throw createError({
      statusCode: 500,
      message: "Failed to fetch expenses",
    });
  }
});
