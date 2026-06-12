<template>
  <div>
    <h1 class="mb-6 text-3xl font-bold text-gray-900 dark:text-white">
      {{ t("dashboard.title") }}
    </h1>

    <div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400">
              {{ t("dashboard.totalRooms") }}
            </span>
            <UIcon name="i-lucide-bed" class="h-5 w-5 text-gray-400" />
          </div>
        </template>
        <div class="text-3xl font-bold">{{ stats?.totalRooms ?? 0 }}</div>
        <p class="mt-1 text-sm text-gray-500">
          {{ t("dashboard.availableRooms", { count: stats?.availableRooms ?? 0 }) }}
        </p>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400">
              {{ t("dashboard.occupancy") }}
            </span>
            <UIcon name="i-lucide-users" class="h-5 w-5 text-gray-400" />
          </div>
        </template>
        <div class="text-3xl font-bold">{{ formatPercent(stats?.occupancyRate ?? 0) }}</div>
        <p class="mt-1 text-sm text-gray-500">
          {{ t("dashboard.roomsOccupied", { count: stats?.occupiedRooms ?? 0 }) }}
        </p>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400">{{
              t("dashboard.todaysRevenue")
            }}</span>
            <UIcon name="i-lucide-dollar-sign" class="h-5 w-5 text-gray-400" />
          </div>
        </template>
        <div class="text-3xl font-bold">{{ formatCurrency(stats?.todaysRevenue ?? 0) }}</div>
        <p class="mt-1 text-sm text-green-500">{{ t("dashboard.revenueDelta") }}</p>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400">{{
              t("dashboard.checkInsToday")
            }}</span>
            <UIcon name="i-lucide-log-in" class="h-5 w-5 text-gray-400" />
          </div>
        </template>
        <div class="text-3xl font-bold">{{ stats?.checkInsToday ?? 0 }}</div>
        <p class="mt-1 text-sm text-gray-500">
          {{ t("dashboard.checkOuts", { count: stats?.checkOutsToday ?? 0 }) }}
        </p>
      </UCard>
    </div>

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold">{{ t("dashboard.recentReservations") }}</h2>
            <UButton variant="ghost" size="sm" to="/reservations">{{
              t("actions.viewAll")
            }}</UButton>
          </div>
        </template>
        <div class="space-y-4">
          <div
            v-for="r in recentReservations"
            :key="r.id"
            class="flex items-center justify-between border-b border-gray-100 py-2 last:border-0 dark:border-gray-800"
          >
            <div>
              <p class="font-medium">{{ r.guestName }}</p>
              <p class="text-sm text-gray-500">
                {{ t("dashboard.roomLabel", { number: r.roomNumber || "-" }) }}
              </p>
            </div>
            <div class="text-right">
              <UBadge color="success" variant="soft">{{ r.status }}</UBadge>
              <p class="mt-1 text-xs text-gray-500">{{ formatDate(r.entryDate) }}</p>
            </div>
          </div>
          <div v-if="!recentReservations?.length" class="py-4 text-center text-gray-500">
            {{ t("dashboard.noReservations") }}
          </div>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold">{{ t("dashboard.quickActions") }}</h2>
          </div>
        </template>
        <div class="grid grid-cols-2 gap-4">
          <UButton to="/guests/create" block size="lg" color="primary" variant="solid">
            <UIcon name="i-lucide-user-plus" class="mr-2" />
            {{ t("dashboard.newGuestCheckIn") }}
          </UButton>
          <UButton to="/reservations/create" block size="lg" color="neutral">
            <UIcon name="i-lucide-calendar-plus" class="mr-2" />
            {{ t("dashboard.newReservation") }}
          </UButton>
          <UButton to="/rooms" block size="lg" color="neutral">
            <UIcon name="i-lucide-bed" class="mr-2" />
            {{ t("dashboard.manageRooms") }}
          </UButton>
          <UButton to="/parking/check-in" block size="lg" color="neutral">
            <UIcon name="i-lucide-car" class="mr-2" />
            {{ t("dashboard.parkingCheckIn") }}
          </UButton>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getApiDashboardStats, getApiDashboardRecentReservations } from "~/utils/client";

const { t } = useI18n();

const { data: stats } = useAsyncData("dashboard-stats", async () => {
  const res = await getApiDashboardStats({ requestValidator: undefined });
  return res.data;
});

interface RecentReservation {
  id?: number;
  reservationCode?: string;
  guestName?: string;
  roomNumber?: string;
  status?: string;
  entryDate?: string;
}

const { data: recentReservations } = useAsyncData<RecentReservation[]>(
  "dashboard-recent-reservations",
  async () => {
    const res = await getApiDashboardRecentReservations({ requestValidator: undefined });
    return (res.data?.data ?? []) as RecentReservation[];
  }
);

const formatPercent = (val: number) => {
  return `${val.toFixed(0)}%`;
};

const formatCurrency = (val: number) => {
  return `$${val.toFixed(2)}`;
};

const { locale } = useI18n();

const formatDate = (date: string | undefined) => {
  if (!date) return "-";
  return new Date(date).toLocaleDateString(locale.value, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};
</script>
