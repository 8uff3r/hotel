<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("restaurant.title") }}</h1>
    </div>

    <!-- Sub Navigation -->
    <div class="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-3">
      <UCard
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="cursor-pointer transition-all hover:ring-2 hover:ring-primary-500"
      >
        <div class="flex items-center gap-4">
          <div class="flex h-12 w-12 items-center justify-center rounded-lg" :class="item.bgColor">
            <UIcon :name="item.icon" class="h-6 w-6 text-white" />
          </div>
          <div>
            <h3 class="font-semibold text-gray-900 dark:text-white">
              {{ t(item.label) }}
            </h3>
            <p class="text-sm text-gray-500">{{ item.description }}</p>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Stats Overview -->
    <UCard v-if="statsPending" class="mb-6">
      <div class="animate-pulse">
        <div class="skeleton h-20" />
      </div>
    </UCard>
    <UCard v-else-if="stats" class="mb-6">
      <template #header>
        <span class="text-lg font-semibold">{{ t("restaurant.statistics") }}</span>
      </template>
      <div class="grid grid-cols-2 gap-4 sm:grid-cols-5">
        <div class="text-center">
          <div class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ stats.totalBills }}
          </div>
          <div class="text-sm text-gray-500">{{ t("restaurant.totalBills") }}</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ stats.totalMeals }}
          </div>
          <div class="text-sm text-gray-500">{{ t("restaurant.totalMeals") }}</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-success-600 dark:text-success-400">
            ${{ stats.internalRevenue?.toFixed(2) }}
          </div>
          <div class="text-sm text-gray-500">{{ t("restaurant.internalRevenue") }}</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-warning-600 dark:text-warning-400">
            ${{ stats.externalRevenue?.toFixed(2) }}
          </div>
          <div class="text-sm text-gray-500">{{ t("restaurant.externalRevenue") }}</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-primary-600 dark:text-primary-400">
            ${{ stats.totalRevenue?.toFixed(2) }}
          </div>
          <div class="text-sm text-gray-500">{{ t("restaurant.totalRevenue") }}</div>
        </div>
      </div>
    </UCard>

    <slot />
  </div>
</template>

<script setup lang="ts">
import type { RestaurantStats } from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.restaurant.restaurant.read,
});

const { t } = useI18n();

const navItems = computed(() =>
  [
    {
      label: "restaurant.inventory",
      description: t("restaurant.inventory"),
      icon: "i-lucide-package",
      to: "/restaurant/inventory",
      permission: PERMISSIONS.restaurant.restaurantInventory.read,
      bgColor: "bg-blue-500",
    },
    {
      label: "restaurant.serving",
      description: t("restaurant.serving"),
      icon: "i-lucide-utensils",
      to: "/restaurant/serving",
      permission: PERMISSIONS.restaurant.restaurantServing.read,
      bgColor: "bg-green-500",
    },
    {
      label: "restaurant.reports",
      description: t("restaurant.reports"),
      icon: "i-lucide-bar-chart-3",
      to: "/restaurant/reports",
      permission: PERMISSIONS.restaurant.restaurantReports.read,
      bgColor: "bg-purple-500",
    },
  ].filter((item) => useAuthStore().can(item.permission))
);

const { data: stats, pending: statsPending } = useAsyncData<RestaurantStats>(
  "restaurant-stats",
  async () => {
    return await $fetch("/api/restaurant/stats");
  }
);
</script>
