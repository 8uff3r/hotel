import { tables } from "~~/server/db/schema";
import { eq } from "drizzle-orm";
import { db } from "@nuxthub/db";
import { verifyJWT } from "../../utils/jwt";

export default defineEventHandler(async (event) => {
  try {
    // Get JWT from cookie
    const token = getCookie(event, "auth_token");

    if (!token) {
      throw createError({
        statusCode: 401,
        message: "Unauthorized",
      });
    }

    // Verify JWT
    const payload = verifyJWT(token);

    if (!payload) {
      throw createError({
        statusCode: 401,
        message: "Invalid token",
      });
    }

    // Fetch user from database
    const [user] = await db
      .select()
      .from(tables.users)
      .where(eq(tables.users.id, payload.userId))
      .limit(1);

    if (!user) {
      throw createError({
        statusCode: 404,
        message: "User not found",
      });
    }

    // Check if user is active
    if (!user.isActive) {
      throw createError({
        statusCode: 403,
        message: "Account is deactivated",
      });
    }

    // Return user without password hash
    const { passwordHash: _, ...userWithoutPassword } = user;

    // Get roles from user_roles table
    const userRoles = await db
      .select()
      .from(tables.userRoles)
      .where(eq(tables.userRoles.userId, user.id));

    const roles = userRoles.length > 0 ? userRoles.map((r) => r.role) : [user.role || "staff"];

    return {
      user: { ...userWithoutPassword, roles },
    };
  } catch (error: unknown) {
    if (error && typeof error === "object" && "statusCode" in error) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "An error occurred while fetching user",
    });
  }
});
