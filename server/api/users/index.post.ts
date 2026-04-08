import { tables } from "~~/server/db/schema";
import { z } from "zod";
import { eq } from "drizzle-orm";
import bcrypt from "bcrypt";
import { db } from "@nuxthub/db";

const createUserSchema = z.object({
  email: z.email("Invalid email address"),
  password: z.string().min(6, "Password must be at least 6 characters"),
  firstName: z.string().min(1, "First name is required"),
  lastName: z.string().min(1, "Last name is required"),
  roles: z
    .array(z.enum(["admin", "manager", "receptionist", "staff"]))
    .min(1, "At least one role is required"),
  isActive: z.boolean().default(true),
});

export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  // Validate input
  const validation = createUserSchema.safeParse(body);
  if (!validation.success) {
    const firstIssue = validation.error.issues[0];
    throw createError({
      statusCode: 400,
      message: firstIssue?.message || "Invalid input",
    });
  }

  const { email, password, firstName, lastName, roles, isActive } = validation.data;

  // Check if email already exists
  const [existingUser] = await db
    .select()
    .from(tables.users)
    .where(eq(tables.users.email, email))
    .limit(1);

  if (existingUser) {
    throw createError({
      statusCode: 400,
      message: "Email already exists",
    });
  }

  // Hash password
  const passwordHash = await bcrypt.hash(password, 10);

  // Create user
  const [newUser] = await db
    .insert(tables.users)
    .values({
      email,
      passwordHash,
      firstName,
      lastName,
      isActive,
    })
    .returning();

  if (!newUser) {
    throw createError({
      statusCode: 500,
      message: "Failed to create user",
    });
  }

  // Assign roles
  for (const role of roles) {
    await db.insert(tables.userRoles).values({
      userId: newUser.id,
      role,
    });
  }

  // Return user without password hash
  const { passwordHash: _, ...userWithoutPassword } = newUser;

  return {
    user: { ...userWithoutPassword, roles } as any,
    message: "User created successfully",
  };
});
