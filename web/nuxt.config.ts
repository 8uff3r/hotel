// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2026-04-12",
  devtools: { enabled: true },

  modules: [
    "@nuxt/image",
    "@nuxt/ui",
    "@nuxt/icon",
    "@nuxt/fonts",
    "@pinia/nuxt",
    "@nuxtjs/i18n",
    "@vueuse/nuxt",
    "pinia-plugin-persistedstate/nuxt",
  ],

  i18n: {
    baseUrl: "",
    defaultLocale: "fa",
    defaultDirection: "rtl",
    detectBrowserLanguage: false,
    strategy: "no_prefix",
    langDir: "locales",
    locales: [
      { code: "fa", dir: "rtl", file: "fa.json", name: "فارسی" },
      { code: "en", dir: "ltr", file: "en.json", name: "English" },
      // { code: "ar", dir: "rtl", file: "ar.json" },
    ],
  },
  css: ["~/assets/css/main.css"],
  ssr: false,

  fonts: {
    provider: "local",
    defaults: {
      weights: [100, 200, 300, 400, 500, 600, 700, 800, 900],
    },
    families: [
      { name: "Vazirmatn FD", provider: "local" },
      { name: "Vazirmatn", provider: "local" },
    ],
  },
  icon: {
    provider: "none",
    clientBundle: {
      scan: {
        globExclude: ["dist", "build", "coverage", "test", "tests", ".*"],
        globInclude: ["app/**/*.{vue,jsx,tsx,md,mdc,mdx}", "node_modules/@nuxt/ui/**/*"],
      },
    },
    collections: ["lucide"],
    fallbackToApi: false,
    size: "16px",
  },

  vite: {
    server: {
      proxy: {
        "/api": {
          target: "http://127.0.0.1:8080",
          changeOrigin: true,
        },
        "/healthz": {
          target: "http://127.0.0.1:8080",
          changeOrigin: true,
        },
        "/readyz": {
          target: "http://127.0.0.1:8080",
          changeOrigin: true,
        },
      },
    },
  },
  nitro: {
    preset: "static",
  },

  // Runtime config for environment variables
  runtimeConfig: {
    // Server-side only
    sessionSecret: process.env.NUXT_SESSION_SECRET || "change-this-secret-in-production",
    adminPassword: process.env.NUXT_ADMIN_PASSWORD,
    adminEmail: process.env.NUXT_ADMIN_EMAIL || "admin@hotel.parsiansh.ir",

    // Public runtime config
    public: {
      hotelName: process.env.NUXT_PUBLIC_HOTEL_NAME || "HMS",
      backendUrl: process.env.BACKEND_URL || "http://127.0.0.1:8080",
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
