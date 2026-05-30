<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("accounting.chart_of_accounts") }}
      </h1>
      <UButton to="/accounting/accounts/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />
        {{ t("accounting.add_account") }}
      </UButton>
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          :placeholder="t('accounting.search_by_account_name_or_code')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.accountType"
          :items="accountTypeOptions"
          :placeholder="t('common.all_types')"
          class="w-full sm:w-40"
          @change="fetchAccounts"
        />
        <UButton variant="outline" @click="clearFilters">{{ t("common.clear") }}</UButton>
      </div>
    </UCard>

    <!-- Accounts Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("accounting.account_list") }}</span>
          <span class="text-sm text-gray-500">{{ pagination.total }} accounts</span>
        </div>
      </template>

      <UTable :data="accounts" :columns="columns" :loading="loading" striped>
        <template #accountCode-cell="{ row }">
          <span class="font-mono font-medium">{{ row.original.accountCode }}</span>
        </template>

        <template #accountName-cell="{ row }">
          <div>
            <p class="font-medium">{{ row.original.accountName }}</p>
            <p v-if="row.original.description" class="text-sm text-gray-500">
              {{ row.original.description }}
            </p>
          </div>
        </template>

        <template #accountType-cell="{ row }">
          <UBadge :color="getTypeColor(row.original.accountType)" variant="soft">
            {{ formatAccountType(row.original.accountType) }}
          </UBadge>
        </template>

        <template #normalBalance-cell="{ row }">
          <span
            :class="row.original.normalBalance === 'debit' ? 'text-blue-600' : 'text-orange-600'"
          >
            {{ row.original.normalBalance === "debit" ? "Debit" : "Credit" }}
          </span>
        </template>

        <template #isActive-cell="{ row }">
          <UBadge :color="row.original.isActive ? 'success' : 'error'" variant="soft">
            {{ row.original.isActive ? "Active" : "Inactive" }}
          </UBadge>
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/accounting/accounts/${row.original.id}`">
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
            @change="fetchAccounts"
          />
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { PaginatedResponseModelsAccount } from "~/utils/client";
import type { TableColumn } from "@nuxt/ui";

definePageMeta({
  requiresPermission: ["admin", "manager"],
});

type Account = NonNullable<PaginatedResponseModelsAccount["data"]>[0];
const columns: TableColumn<Account>[] = [
  { accessorKey: "accountCode", header: "Code" },
  { accessorKey: "accountName", header: "Account Name" },
  { accessorKey: "accountType", header: "Type" },
  { accessorKey: "normalBalance", header: "Normal Balance" },
  { accessorKey: "isActive", header: "Status" },
  { accessorKey: "actions", header: "Actions" },
];

const accountTypeOptions = [
  { value: "all", label: "All Types" },
  { value: "asset", label: "Asset" },
  { value: "liability", label: "Liability" },
  { value: "equity", label: "Equity" },
  { value: "revenue", label: "Revenue" },
  { value: "expense", label: "Expense" },
];

const accounts = ref<Account[]>([]);
const { t } = useI18n();
const loading = ref(false);
const page = ref(1);

const filters = reactive({
  search: "",
  accountType: "",
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
    fetchAccounts();
  }, 300);
};

const fetchAccounts = async () => {
  loading.value = true;
  try {
    const query: Record<string, any> = {
      page: pagination.page,
      limit: pagination.limit,
    };

    if (filters.search) query["search"] = filters.search;
    if (filters.accountType && filters.accountType !== "all")
      query["accountType"] = filters.accountType;

    const response = await getApiAccountingAccounts({});
    accounts.value = response.data?.data ?? [];
    pagination.total = response.data?.data?.total ?? 0;
    pagination.totalPages = response.data?.data?.totalPages ?? 1;
  } catch (error) {
    console.error("Failed to fetch accounts:", error);
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  filters.search = "";
  filters.accountType = "";
  pagination.page = 1;
  fetchAccounts();
};

const formatAccountType = (type: string | undefined): string => {
  if (!type) return "";
  const types: Record<string, string> = {
    asset: "Asset",
    liability: "Liability",
    equity: "Equity",
    revenue: "Revenue",
    expense: "Expense",
  };
  return types[type] || type;
};

const getTypeColor = (
  type: string | undefined
): "info" | "warning" | "success" | "error" | "neutral" => {
  if (!type) return "neutral";
  const colors: Record<string, "info" | "warning" | "success" | "error" | "neutral"> = {
    asset: "info",
    liability: "warning",
    equity: "success",
    revenue: "success",
    expense: "error",
  };
  return colors[type] || "neutral";
};

onMounted(fetchAccounts);
</script>
