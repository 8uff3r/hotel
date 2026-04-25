<template>
  <div>
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("parking.parking_management") }}
      </h1>
      <p class="mt-1 text-gray-500 dark:text-gray-400">
        {{ t("parking.manage_parking_lots_spots_and_transactions") }}
      </p>
    </div>

    <div class="mb-6 grid grid-cols-1 gap-4 md:grid-cols-4">
      <UCard>
        <template #header>
          <span class="text-sm text-gray-500">{{ t("parking.total_lots") }}</span>
        </template>
        <div class="text-3xl font-bold">{{ stats?.lots }}</div>
      </UCard>

      <UCard>
        <template #header>
          <span class="text-sm text-gray-500">Total Spots</span>
        </template>
        <div class="text-3xl font-bold">{{ stats?.spots }}</div>
      </UCard>

      <UCard>
        <template #header>
          <span class="text-sm text-gray-500">Available</span>
        </template>
        <div class="text-3xl font-bold text-green-600">{{ stats?.availableSpots }}</div>
      </UCard>

      <UCard>
        <template #header>
          <span class="text-sm text-gray-500">Active Vehicles</span>
        </template>
        <div class="text-3xl font-bold text-blue-600">{{ 0 }}</div>
      </UCard>
    </div>

    <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-semibold">Parking Lots</span>
            <UButton size="sm" to="/parking/lots/create" color="primary">
              <UIcon name="i-lucide-plus" class="mr-1" />
              Add Lot
            </UButton>
          </div>
        </template>
        <div class="space-y-3">
          <div
            v-for="lot in lots"
            :key="lot.id"
            class="flex items-center justify-between rounded-lg bg-gray-50 p-3 dark:bg-gray-800"
          >
            <div>
              <div class="font-medium">{{ lot.name }}</div>
              <div class="text-sm text-gray-500">
                {{ lot.totalSpots }} spots - {{ lot.hourlyRate }}/hr
              </div>
            </div>
            <UBadge :color="lot.status === 'active' ? 'success' : 'warning'" variant="soft">
              {{ lot.status }}
            </UBadge>
          </div>
          <div v-if="lots.length === 0" class="py-4 text-center text-gray-500">
            No parking lots found
          </div>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-semibold">Recent Transactions</span>
            <UButton size="sm" to="/parking/transactions" variant="outline"> View All </UButton>
          </div>
        </template>
        <div class="space-y-3">
          <div
            v-for="tx in recentTransactions"
            :key="tx.id"
            class="flex items-center justify-between rounded-lg bg-gray-50 p-3 dark:bg-gray-800"
          >
            <div>
              <div class="font-medium">{{ tx.licensePlate }}</div>
              <div class="text-sm text-gray-500">{{ formatDate(tx.entryTime) }}</div>
            </div>
            <UBadge :color="tx.status === 'active' ? 'info' : 'success'" variant="soft">
              {{ tx.status }}
            </UBadge>
          </div>
          <div v-if="recentTransactions.length === 0" class="py-4 text-center text-gray-500">
            No recent transactions
          </div>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});
const { t } = useI18n();

// const stats = reactive({
//   totalLots: 0,
//   totalSpots: 0,
//   availableSpots: 0,
//   activeTransactions: 0,
// });

const lots = ref<any[]>([]);
const recentTransactions = ref<any[]>([]);

const { data: stats } = useAsyncData(async () => {
  const res = await $api<ParkingStats>("/api/parking/stats");
  return res;
});

const formatDate = (date: string) => {
  return new Date(date).toLocaleString();
};
</script>
