import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = getRouterParam(event, "id");

  if (!id || isNaN(Number(id))) {
    throw createError({
      statusCode: 400,
      message: "Invalid guest ID",
    });
  }

  try {
    const [guest] = await db
      .select()
      .from(tables.guests)
      .where(eq(tables.guests.id, Number(id)))
      .limit(1);

    if (!guest) {
      throw createError({
        statusCode: 404,
        message: "Guest not found",
      });
    }

    return {
      data: guest,
    };
  } catch (error: unknown) {
    if (error && typeof error === "object" && "statusCode" in error) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "Failed to fetch guest",
    });
  }
});
