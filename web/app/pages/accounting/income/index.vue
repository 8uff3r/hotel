hotel/app/pages/accounting/income/index.vue ``` ```vue
<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Income</h1>
      <UButton to="/accounting/income/create" color="success">
        <UIcon name="i-lucide-plus" class="mr-2" />
        Record Income
      </UButton>
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          placeholder="Search by description..."
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.category"
          :items="categoryOptions"
          placeholder="All Categories"
          class="w-full sm:w-40"
          @change="fetchIncome"
        />
        <USelect
          v-model="filters.paymentStatus"
          :items="paymentStatusOptions"
          placeholder="All Payments"
          class="w-full sm:w-40"
          @change="fetchIncome"
        />
        <div class="flex gap-2">
          <UInput
            v-model="filters.dateFrom"
            type="date"
            placeholder="From"
            class="w-full sm:w-40"
            @change="fetchIncome"
          />
          <UInput
            v-model="filters.dateTo"
            type="date"
            placeholder="To"
            class="w-full sm:w-40"
            @change="fetchIncome"
          />
        </div>
        <UButton variant="outline" @click="clearFilters"> Clear </UButton>
      </div>
    </UCard>

    <!-- Income Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">Income Records</span>
          <span class="text-sm text-gray-500">{{ pagination.total }} records</span>
        </div>
      </template>

      <UTable :data="incomeList" :columns="columns" :loading="loading" striped>
        <template #id-cell="{ row }">
          <NuxtLink
            :to="`/accounting/income/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            #{{ row.original.id }}
          </NuxtLink>
        </template>

        <template #description-cell="{ row }">
          <div>
            <p class="font-medium">{{ row.original.description }}</p>
            <p class="text-sm text-gray-500 capitalize">
              {{ formatCategory(row.original.category) }}
            </p>
          </div>
        </template>

        <template #amount-cell="{ row }">
          <span class="font-medium text-green-600"
            >${{ Number(row.original.amount).toFixed(2) }}</span
          >
        </template>

        <template #date-cell="{ row }">
          {{ formatDate(row.original.incomeDate) }}
        </template>

        <template #paymentStatus-cell="{ row }">
          <UBadge :color="getStatusColor(row.original.paymentStatus)" variant="soft">
            {{ row.original.paymentStatus }}
          </UBadge>
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/accounting/income/${row.original.id}`">
              <UIcon name="i-lucide-eye" class="h-4 w-4" />
            </UButton>
          </div>
        </template>
      </UTable>

      <template #footer>
        <div class="flex items-center justify-between">
          <span class="text-sm text-gray-500">
            Page {{ pagination.page }} of {{ pagination.totalPages }}
          </span>
          <UPagination
            v-model="page"
            :page-count="pagination.limit"
            :total="pagination.total"
            @change="fetchIncome"
          />
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";

definePageMeta({
  requiresRole: ["admin", "manager"],
});

interface IncomeRow {
  id: number;
  description: string;
  amount: string | number;
  category: string;
  incomeDate: string;
  paymentStatus: string;
  source: string | null;
}

const columns: TableColumn<IncomeRow>[] = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "description", header: "Description" },
  { accessorKey: "amount", header: "Amount" },
  { accessorKey: "date", header: "Date" },
  { accessorKey: "paymentStatus", header: "Status" },
  { accessorKey: "actions", header: "Actions" },
];

const categoryOptions = [
  { value: "all", label: "All Categories" },
  { value: "room_revenue", label: "Room Revenue" },
  { value: "food_beverage", label: "Food & Beverage" },
  { value: "laundry", label: "Laundry" },
  { value: "spa", label: "Spa" },
  { value: "meeting_rooms", label: "Meeting Rooms" },
  { value: "other", label: "Other" },
];

const paymentStatusOptions = [
  { value: "all", label: "All Payments" },
  { value: "pending", label: "Pending" },
  { value: "received", label: "Received" },
  { value: "refunded", label: "Refunded" },
];

const incomeList = ref<IncomeRow[]>([]);
const loading = ref(false);
const page = ref(1);

const filters = reactive({
  search: "",
  category: "",
  paymentStatus: "",
  dateFrom: "",
  dateTo: "",
});

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  totalPages: 0,
});

let searchTimeout: ReturnType<typeof setTimeout> | null = null;

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    pagination.page = 1;
    fetchIncome();
  }, 300);
};

const fetchIncome = async () => {
  loading.value = true;
  try {
    const params = new URLSearchParams();
    params.append("page", pagination.page.toString());
    params.append("limit", pagination.limit.toString());

    if (filters.search) params.append("search", filters.search);
    if (filters.category && filters.category !== "all") params.append("category", filters.category);
    if (filters.paymentStatus && filters.paymentStatus !== "all")
      params.append("paymentStatus", filters.paymentStatus);
    if (filters.dateFrom) params.append("dateFrom", filters.dateFrom);
    if (filters.dateTo) params.append("dateTo", filters.dateTo);

    const response = await $fetch(`/api/income?${params.toString()}`);
    incomeList.value = response.data;
    pagination.total = response.pagination.total;
    pagination.totalPages = response.pagination.totalPages;
  } catch (error) {
    console.error("Failed to fetch income:", error);
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  filters.search = "";
  filters.category = "";
  filters.paymentStatus = "";
  filters.dateFrom = "";
  filters.dateTo = "";
  pagination.page = 1;
  fetchIncome();
};

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};

const formatCategory = (category: string): string => {
  const categories: Record<string, string> = {
    room_revenue: "Room Revenue",
    food_beverage: "Food & Beverage",
    laundry: "Laundry",
    spa: "Spa",
    meeting_rooms: "Meeting Rooms",
    other: "Other",
  };
  return categories[category] || category;
};

const getStatusColor = (status: string): "success" | "warning" | "error" => {
  const colors: Record<string, "success" | "warning" | "error"> = {
    received: "success",
    pending: "warning",
    refunded: "error",
  };
  return colors[status] || "warning";
};

onMounted(fetchIncome);
</script>
