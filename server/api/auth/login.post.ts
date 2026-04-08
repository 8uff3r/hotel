import { tables } from "~~/server/db/schema";
import { z } from "zod";
import { eq } from "drizzle-orm";
import bcrypt from "bcrypt";
import { db } from "@nuxthub/db";
import { createJWT } from "../../utils/jwt";

const loginSchema = z.object({
  email: z.email("Invalid email address"),
  password: z.string().min(6, "Password must be at least 6 characters"),
});

export default defineEventHandler(async (event) => {
  const body = await readBody(event);
  console.log(process.env.DATABASE_URL);

  // Validate input
  const validation = loginSchema.safeParse(body);
  if (!validation.success) {
    const firstIssue = validation.error.issues[0];
    throw createError({
      statusCode: 400,
      message: firstIssue?.message || "Invalid input",
    });
  }

  const { email, password } = validation.data;

  try {
    // Find user by email
    const [user] = await db
      .select()
      .from(tables.users)
      .where(eq(tables.users.email, email))
      .limit(1);

    if (!user) {
      throw createError({
        statusCode: 401,
        message: "Invalid email or password",
      });
    }

    // Check if user is active
    if (!user.isActive) {
      throw createError({
        statusCode: 403,
        message: "Account is deactivated. Please contact administrator.",
      });
    }

    // Verify password
    const isValidPassword = await bcrypt.compare(password, user.passwordHash);

    if (!isValidPassword) {
      throw createError({
        statusCode: 401,
        message: "Invalid email or password",
      });
    }

    // Get roles from user_roles table
    const userRoles = await db
      .select()
      .from(tables.userRoles)
      .where(eq(tables.userRoles.userId, user.id));

    const roles = userRoles.length > 0 ? userRoles.map((r) => r.role) : [user.role || "staff"];

    // Create JWT token
    const token = createJWT({
      userId: user.id,
      email: user.email,
      roles: roles,
    });

    // Set JWT as httpOnly cookie
    setCookie(event, "auth_token", token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      maxAge: 60 * 60 * 24 * 7, // 7 days
      path: "/",
    });

    // Return user without password hash
    const { passwordHash: _, ...userWithoutPassword } = user;

    return {
      user: { ...userWithoutPassword, roles },
    };
  } catch (error: any) {
    console.log(error);
    if (error.statusCode) {
      throw error;
    }

    throw createError({
      statusCode: 500,
      message: "An error occurred during login",
    });
  }
});
