<template>
  <div class="space-y-6">
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <UCard>
        <div class="flex items-center gap-4">
          <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-blue-500">
            <UIcon name="i-lucide-receipt" class="h-6 w-6 text-white" />
          </div>
          <div>
            <div class="text-2xl font-bold">{{ stats?.totalBills ?? 0 }}</div>
            <div class="text-sm text-gray-500">{{ t("restaurant.totalBills") }}</div>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center gap-4">
          <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-green-500">
            <UIcon name="i-lucide-utensils" class="h-6 w-6 text-white" />
          </div>
          <div>
            <div class="text-2xl font-bold">{{ stats?.totalMeals ?? 0 }}</div>
            <div class="text-sm text-gray-500">{{ t("restaurant.totalMeals") }}</div>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center gap-4">
          <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-primary-500">
            <UIcon name="i-lucide-building" class="h-6 w-6 text-white" />
          </div>
          <div>
            <div class="text-2xl font-bold text-primary-600">
              ${{ stats?.internalRevenue?.toFixed(2) ?? "0.00" }}
            </div>
            <div class="text-sm text-gray-500">{{ t("restaurant.internalRevenue") }}</div>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center gap-4">
          <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-warning-500">
            <UIcon name="i-lucide-shop" class="h-6 w-6 text-white" />
          </div>
          <div>
            <div class="text-2xl font-bold text-warning-600">
              ${{ stats?.externalRevenue?.toFixed(2) ?? "0.00" }}
            </div>
            <div class="text-sm text-gray-500">{{ t("restaurant.externalRevenue") }}</div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Revenue Chart Section -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("restaurant.totalRevenue") }}</span>
          <span class="text-2xl font-bold text-primary-600">
            ${{ stats?.totalRevenue?.toFixed(2) ?? "0.00" }}
          </span>
        </div>
      </template>
      <div class="flex items-center justify-center py-8">
        <div class="text-center">
          <div
            class="mb-4 flex h-32 w-32 items-center justify-center rounded-full border-8 border-success-500"
          >
            <div class="text-center">
              <div class="text-xl font-bold">{{ internalPercentage }}%</div>
              <div class="text-xs text-gray-500">{{ t("restaurant.internal") }}</div>
            </div>
          </div>
          <div class="flex items-center justify-center gap-4">
            <div class="flex items-center gap-2">
              <div class="h-3 w-3 rounded-full bg-success-500" />
              <span class="text-sm">{{ t("restaurant.internal") }}</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="h-3 w-3 rounded-full bg-warning-500" />
              <span class="text-sm">{{ t("restaurant.external") }}</span>
            </div>
          </div>
        </div>
      </div>
    </UCard>

    <!-- Recent Bills -->
    <UCard>
      <template #header>
        <span class="text-lg font-semibold">{{ t("restaurant.transactions") }}</span>
      </template>

      <div class="mb-4 flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          :placeholder="t('restaurant.searchPlaceholder')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.status"
          :items="statusOptions"
          :placeholder="t('common.all_statuses')"
          class="w-full sm:w-40"
        />
        <UButton variant="outline" @click="clearFilters">{{ t("common.clear") }}</UButton>
      </div>

      <UTable :data="bills" :columns="columns" :loading="pending" striped>
        <template #billDate-cell="{ row }">
          {{ formatDate(row.original.billDate) }}
        </template>

        <template #guestId-cell="{ row }">
          <NuxtLink
            v-if="row.original.guestId"
            :to="`/guests/${row.original.guestId}`"
            class="text-primary hover:underline"
          >
            #{{ row.original.guestId }}
          </NuxtLink>
          <span v-else>-</span>
        </template>

        <template #totalAmount-cell="{ row }">
          ${{ row.original.totalAmount?.toFixed(2) }}
        </template>

        <template #isExternal-cell="{ row }">
          <UBadge :color="row.original.isExternal ? 'warning' : 'success'" variant="soft">
            {{ row.original.isExternal ? t("restaurant.external") : t("restaurant.internal") }}
          </UBadge>
        </template>

        <template #status-cell="{ row }">
          <UBadge
            :style="{
              backgroundColor: row.original.status?.colorHex,
            }"
            variant="soft"
          >
            {{ row.original.status?.label }}
          </UBadge>
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" @click="viewBill(row.original)">
              <UIcon name="i-lucide-eye" class="h-4 w-4" />
            </UButton>
            <UButton
              v-if="!row.original.settled"
              variant="ghost"
              size="sm"
              color="success"
              @click="settleBill(row.original)"
            >
              <UIcon name="i-lucide-check-circle" class="h-4 w-4" />
            </UButton>
          </div>
        </template>
      </UTable>

      <template #footer>
        <div class="flex items-center justify-between">
          <span class="text-sm text-gray-500">
            {{
              t("pagination.pageOf", { page: pagination.page, totalPages: pagination.totalPages })
            }}
          </span>
          <UPagination v-model="page" :page-count="pagination.limit" :total="pagination.total" />
        </div>
      </template>
    </UCard>

    <!-- Bill Detail Modal -->
    <UModal v-model="detailModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("restaurant.billDetails") }}</h2>
      </template>

      <template #body>
        <div v-if="selectedBill" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <div class="text-sm text-gray-500">{{ t("restaurant.billDetails") }}</div>
              <div class="font-semibold">#{{ selectedBill.id }}</div>
            </div>
            <div>
              <div class="text-sm text-gray-500">{{ t("common.status") }}</div>
              <UBadge
                :style="{
                  backgroundColor: selectedBill.status?.colorHex,
                }"
                variant="soft"
              >
                {{ selectedBill.status?.label }}
              </UBadge>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <div class="text-sm text-gray-500">{{ t("forms.firstName") }}</div>
              <div>{{ selectedBill.guestId ? `#${selectedBill.guestId}` : "-" }}</div>
            </div>
            <div>
              <div class="text-sm text-gray-500">{{ t("forms.phone") }}</div>
              <div>{{ selectedBill.roomId ? `#${selectedBill.roomId}` : "-" }}</div>
            </div>
          </div>

          <div class="border-t pt-4">
            <div class="flex justify-between">
              <span>{{ t("restaurant.subtotal") }}</span>
              <span>${{ selectedBill.subtotal?.toFixed(2) }}</span>
            </div>
            <div class="flex justify-between">
              <span>{{ t("restaurant.tax") }}</span>
              <span>${{ selectedBill.taxAmount?.toFixed(2) }}</span>
            </div>
            <div class="flex justify-between">
              <span>{{ t("restaurant.discount") }}</span>
              <span>-${{ selectedBill.discountAmount?.toFixed(2) }}</span>
            </div>
            <div class="flex justify-between border-t pt-2 font-bold">
              <span>{{ t("restaurant.total") }}</span>
              <span>${{ selectedBill.totalAmount?.toFixed(2) }}</span>
            </div>
          </div>

          <div
            v-if="selectedBill.isExternal"
            class="rounded-lg bg-warning-50 p-2 dark:bg-warning-900"
          >
            <div class="flex items-center gap-2 text-warning-700 dark:text-warning-300">
              <UIcon name="i-lucide-shop" class="h-4 w-4" />
              <span
                >{{ t("restaurant.externalRestaurant") }}:
                {{ selectedBill.externalRestaurant }}</span
              >
            </div>
          </div>

          <div v-if="selectedBill.notes">
            <div class="text-sm text-gray-500">{{ t("forms.notes") }}</div>
            <div>{{ selectedBill.notes }}</div>
          </div>
        </div>
      </template>

      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="detailModalOpen = false">{{
            t("actions.back")
          }}</UButton>
          <UButton
            v-if="selectedBill && !selectedBill.settled"
            color="success"
            :loading="settling"
            @click="confirmSettle"
          >
            {{ t("restaurant.settleBill") }}
          </UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import type {
  PaginatedResponseModelsRestaurantBill,
  RestaurantBill,
  RestaurantStats,
} from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.restaurant.restaurantReports.read,
});

const { t } = useI18n();

const columns = computed<TableColumn<RestaurantBill>[]>(() => [
  { accessorKey: "id", header: "#" },
  { accessorKey: "billDate", header: t("accounting.income_date") },
  { accessorKey: "guestId", header: t("common.guest") },
  { accessorKey: "totalAmount", header: t("restaurant.total") },
  { accessorKey: "isExternal", header: t("restaurant.category") },
  { accessorKey: "status", header: t("common.status") },
  { accessorKey: "actions", header: t("restaurant.columns.actions") },
]);

const { data: statusOptions } = useAsyncData(() => getApiRestaurantBillsStatuses({}), {
  transform: (response) => response.data?.data,
});

const pagination = reactive({ page: 1, limit: 20, total: 0, totalPages: 0 });
const page = computed({
  get: () => pagination.page,
  set: (val) => {
    pagination.page = val;
    fetchBills();
  },
});

const filters = reactive({ search: "", status: "" });
let searchTimeout: ReturnType<typeof setTimeout> | null = null;

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    pagination.page = 1;
    fetchBills();
  }, 300);
};

const clearFilters = () => {
  filters.search = "";
  filters.status = "";
  pagination.page = 1;
  fetchBills();
};

const {
  data: billsData,
  pending,
  refresh: fetchBills,
} = useAsyncData(
  "restaurant-bills",
  async () => {
    const params: any = { page: pagination.page, limit: pagination.limit };
    if (filters.search) params.search = filters.search;
    if (filters.status) params.status = filters.status;

    const response = await $fetch<PaginatedResponseModelsRestaurantBill>("/api/restaurant/bills", {
      query: params,
    });
    pagination.total = response.data?.data?.total ?? 0;
    pagination.totalPages = response.data?.data?.totalPages ?? 0;
    return response.data?.data;
  },
  { watch: [() => pagination.page] }
);

const bills = computed(() => billsData.value ?? []);

const { data: stats } = useAsyncData<RestaurantStats>("restaurant-stats", () =>
  getApiRestaurantStats({})
);

const internalPercentage = computed(() => {
  if (!stats.value || stats.value.totalRevenue === 0) return 0;
  return Math.round(((stats.value.internalRevenue ?? 0) / (stats.value.totalRevenue ?? 1)) * 100);
});

const detailModalOpen = ref(false);
const selectedBill = ref<RestaurantBill | null>(null);
const settling = ref(false);

const viewBill = (bill: RestaurantBill) => {
  selectedBill.value = bill;
  detailModalOpen.value = true;
};

const confirmSettle = async () => {
  if (!selectedBill.value) return;
  settling.value = true;
  try {
    await $fetch(`/api/restaurant/bills/${selectedBill.value.id}/settle`, { method: "POST" });
    detailModalOpen.value = false;
    await fetchBills();
  } catch (error) {
    console.error("Failed to settle bill:", error);
  } finally {
    settling.value = false;
  }
};

const settleBill = async (bill: RestaurantBill) => {
  selectedBill.value = bill;
  await confirmSettle();
};

const formatDate = (date: string | Date | undefined) => {
  if (!date) return "-";
  return new Date(date).toLocaleDateString();
};

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    open: "info",
    settled: "success",
    cancelled: "error",
  };
  return colors[status] || "neutral";
};
</script>
