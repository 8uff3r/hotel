import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";

export default defineEventHandler(async (event) => {
  const id = Number(getRouterParam(event, "id"));

  const [vehicle] = await db.select().from(tables.vehicles).where(eq(tables.vehicles.id, id));

  if (!vehicle) {
    throw createError({
      statusCode: 404,
      message: "Vehicle not found",
    });
  }

  return vehicle;
});
