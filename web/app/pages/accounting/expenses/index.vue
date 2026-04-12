<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("accounting.expenses") }}
      </h1>
      <UButton to="/accounting/expenses/create" color="error">
        <UIcon name="i-lucide-plus" class="mr-2" />{{ t("accounting.record_expense") }}</UButton
      >
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          :placeholder="t('accounting.search_by_description')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.category"
          :items="categoryOptions"
          :placeholder="t('accounting.all_categories')"
          class="w-full sm:w-40"
          @change="fetchExpenses"
        />
        <USelect
          v-model="filters.paymentStatus"
          :items="paymentStatusOptions"
          :placeholder="t('common.all_payments')"
          class="w-full sm:w-40"
          @change="fetchExpenses"
        />
        <div class="flex gap-2">
          <UInput
            v-model="filters.dateFrom"
            type="date"
            :placeholder="t('accounting.from')"
            class="w-full sm:w-40"
            @change="fetchExpenses"
          />
          <UInput
            v-model="filters.dateTo"
            type="date"
            :placeholder="t('accounting.to')"
            class="w-full sm:w-40"
            @change="fetchExpenses"
          />
        </div>
        <UButton variant="outline" @click="clearFilters">{{ t("common.clear") }}</UButton>
      </div>
    </UCard>

    <!-- Expenses Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("accounting.expense_records") }}</span>
          <span class="text-sm text-gray-500">{{ pagination.total }} records</span>
        </div>
      </template>

      <UTable :data="expensesList" :columns="columns" :loading="loading" striped>
        <template #id-cell="{ row }">
          <NuxtLink
            :to="`/accounting/expenses/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            #{{ row.original.id }}
          </NuxtLink>
        </template>

        <template #description-cell="{ row }">
          <div>
            <p class="font-medium">{{ row.original.description }}</p>
            <p class="text-sm text-gray-500">
              {{ formatCategory(row.original.category) }}
              <span v-if="row.original.vendor"> - {{ row.original.vendor }}</span>
            </p>
          </div>
        </template>

        <template #amount-cell="{ row }">
          <span class="font-medium text-red-600"
            >${{ Number(row.original.amount).toFixed(2) }}</span
          >
        </template>

        <template #date-cell="{ row }">
          {{ formatDate(row.original.expenseDate) }}
        </template>

        <template #paymentStatus-cell="{ row }">
          <UBadge :color="getStatusColor(row.original.paymentStatus)" variant="soft">
            {{ row.original.paymentStatus }}
          </UBadge>
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/accounting/expenses/${row.original.id}`">
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
            @change="fetchExpenses"
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

interface ExpenseRow {
  id: number;
  description: string;
  amount: string | number;
  category: string;
  expenseDate: string;
  paymentStatus: string;
  vendor: string | null;
}

const columns: TableColumn<ExpenseRow>[] = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "description", header: "Description" },
  { accessorKey: "amount", header: "Amount" },
  { accessorKey: "date", header: "Date" },
  { accessorKey: "paymentStatus", header: "Status" },
  { accessorKey: "actions", header: "Actions" },
];

const categoryOptions = [
  { value: "all", label: "All Categories" },
  { value: "food_beverage", label: "Food & Beverage" },
  { value: "housekeeping", label: "Housekeeping" },
  { value: "maintenance", label: "Maintenance" },
  { value: "utilities", label: "Utilities" },
  { value: "salaries", label: "Salaries" },
  { value: "marketing", label: "Marketing" },
  { value: "supplies", label: "Supplies" },
  { value: "insurance", label: "Insurance" },
  { value: "taxes", label: "Taxes" },
  { value: "rent", label: "Rent" },
  { value: "other", label: "Other" },
];

const paymentStatusOptions = [
  { value: "all", label: "All Payments" },
  { value: "pending", label: "Pending" },
  { value: "paid", label: "Paid" },
  { value: "cancelled", label: "Cancelled" },
];

const expensesList = ref<ExpenseRow[]>([]);
const { t } = useI18n();
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
    fetchExpenses();
  }, 300);
};

const fetchExpenses = async () => {
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

    const response = await $fetch(`/api/expenses?${params.toString()}`);
    expensesList.value = response.data;
    pagination.total = response.pagination.total;
    pagination.totalPages = response.pagination.totalPages;
  } catch (error) {
    console.error("Failed to fetch expenses:", error);
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
  fetchExpenses();
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
    food_beverage: "Food & Beverage",
    housekeeping: "Housekeeping",
    maintenance: "Maintenance",
    utilities: "Utilities",
    salaries: "Salaries",
    marketing: "Marketing",
    supplies: "Supplies",
    insurance: "Insurance",
    taxes: "Taxes",
    rent: "Rent",
    other: "Other",
  };
  return categories[category] || category;
};

const getStatusColor = (status: string): "success" | "warning" | "error" => {
  const colors: Record<string, "success" | "warning" | "error"> = {
    paid: "success",
    pending: "warning",
    cancelled: "error",
  };
  return colors[status] || "warning";
};

onMounted(fetchExpenses);
</script>
