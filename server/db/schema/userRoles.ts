import { pgTable, serial, text, timestamp, integer, index, uniqueIndex } from "drizzle-orm/pg-core";
import { users } from "./users";

export const userRoles = ["admin", "manager", "receptionist", "staff"] as const;

export const userRolesTable = pgTable(
  "user_roles",
  {
    id: serial("id").primaryKey(),
    userId: integer("user_id").references(() => users.id, { onDelete: "cascade" }),
    role: text("role", { enum: userRoles }).notNull(),
    createdAt: timestamp("created_at").notNull().defaultNow(),
  },
  (table) => [
    uniqueIndex("user_roles_user_role_idx").on(table.userId, table.role),
    index("user_roles_user_id_idx").on(table.userId),
  ]
);

export type UserRole = typeof userRolesTable.$inferSelect;
export type NewUserRole = typeof userRolesTable.$inferInsert;
