import bcrypt from "bcrypt";
import { db, schema } from "@nuxthub/db";

export default defineTask({
  meta: {
    name: "db:seed",
    description: "Seed database with initial data",
  },
  async run() {
    console.log("Seeding database...");

    const config = useRuntimeConfig();
    if (!config.adminPassword) {
      throw new Error("ADMIN_PASSWORD environment variable is not set");
    }

    const users = [
      {
        email: config.adminEmail || "admin@hotel.parsiansh.ir",
        passwordHash: await bcrypt.hash(config.adminPassword, 10),
        firstName: "admin",
        lastName: "",
        isActive: true,
      },
    ];

    await db.insert(schema.users).values(users);

    return { result: "Database seeded successfully" };
  },
});
