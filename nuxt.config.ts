// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },

  modules: [
    "@nuxt/image",
    "@nuxt/ui",
    "@nuxt/icon",
    "@nuxt/fonts",
    "@pinia/nuxt",
    "@nuxthub/core",
    "pinia-plugin-persistedstate/nuxt",
  ],

  css: ["~/assets/css/main.css"],

  fonts: {
    provider: "local",
  },
  icon: {
    provider: "server",
    collections: ["lucide"],
    fallbackToApi: false,
    size: "16px",
  },

  hub: {
    // Database configuration - PostgreSQL with PGlite for local development
    db: {
      applyMigrationsDuringBuild: false,
      dialect: "postgresql",
      casing: "snake_case",
    },
    // Enable blob storage for file uploads
    blob: true,
    // Enable KV storage for caching
    kv: true,
    // Enable cache
    cache: true,
  },
  nitro: {
    preset: "bun",
    experimental: {
      tasks: true,
    },
  },

  // Runtime config for environment variables
  runtimeConfig: {
    // Server-side only
    sessionSecret: process.env.NUXT_SESSION_SECRET || "change-this-secret-in-production",
    adminPassword: process.env.NUXT_ADMIN_PASSWORD,
    adminEmail: process.env.NUXT_ADMIN_EMAIL || "admin@hotel.parsiansh.ir",

    // Public runtime config
    public: {
      hotelName: process.env.NUXT_PUBLIC_HOTEL_NAME || "Hotel Management System",
    },
  },

  // Pinia configuration
  pinia: {
    storesDirs: ["./app/stores/**"],
  },

  // TypeScript configuration
  typescript: {
    strict: true,
    typeCheck: false,
  },
});
