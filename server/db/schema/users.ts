import { pgTable, serial, text, timestamp, boolean, index } from "drizzle-orm/pg-core";

export const users = pgTable(
  "users",
  {
    id: serial("id").primaryKey(),
    email: text("email").notNull().unique(),
    passwordHash: text("password_hash").notNull(),
    firstName: text("first_name").notNull(),
    lastName: text("last_name").notNull(),
    role: text("role", { enum: ["admin", "manager", "receptionist", "staff"] })
      .notNull()
      .default("staff"),
    avatar: text("avatar"),
    isActive: boolean("is_active").notNull().default(true),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => [index("users_email_idx").on(table.email), index("users_role_idx").on(table.role)]
);

export type User = typeof users.$inferSelect;
export type NewUser = typeof users.$inferInsert;
