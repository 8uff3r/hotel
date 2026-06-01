<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("accounting.accounting") }}
      </h1>
    </div>

    <!-- Summary Cards -->
    <div class="mb-6 grid grid-cols-1 gap-4 md:grid-cols-4">
      <UCard>
        <div class="flex items-center">
          <div
            class="flex h-12 w-12 items-center justify-center rounded-full bg-green-100 dark:bg-green-900"
          >
            <UIcon name="i-lucide-trending-up" class="h-6 w-6 text-green-600 dark:text-green-400" />
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-500">{{ t("accounting.total_income") }}</p>
            <p class="text-2xl font-bold text-green-600">${{ totalIncome.toFixed(2) }}</p>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div
            class="flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900"
          >
            <UIcon name="i-lucide-trending-down" class="h-6 w-6 text-red-600 dark:text-red-400" />
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-500">{{ t("accounting.total_expenses") }}</p>
            <p class="text-2xl font-bold text-red-600">${{ totalExpenses.toFixed(2) }}</p>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div
            class="flex h-12 w-12 items-center justify-center rounded-full bg-blue-100 dark:bg-blue-900"
          >
            <UIcon name="i-lucide-wallet" class="h-6 w-6 text-blue-600 dark:text-blue-400" />
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-500">{{ t("accounting.net_balance") }}</p>
            <p :class="['text-2xl font-bold', netBalance >= 0 ? 'text-green-600' : 'text-red-600']">
              ${{ netBalance.toFixed(2) }}
            </p>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div
            class="flex h-12 w-12 items-center justify-center rounded-full bg-purple-100 dark:bg-purple-900"
          >
            <UIcon name="i-lucide-receipt" class="h-6 w-6 text-purple-600 dark:text-purple-400" />
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-500">{{ t("accounting.pending_payments") }}</p>
            <p class="text-2xl font-bold text-purple-600">${{ pendingPayments.toFixed(2) }}</p>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Quick Actions -->
    <div class="mb-6 grid grid-cols-1 gap-4 md:grid-cols-3">
      <UButton to="/accounting/income/create" color="success" block size="lg">
        <UIcon name="i-lucide-plus" class="mr-2" />{{ t("accounting.record_income") }}</UButton
      >
      <UButton to="/accounting/expenses/create" color="error" variant="outline" block size="lg">
        <UIcon name="i-lucide-minus" class="mr-2" />{{ t("accounting.record_expense") }}</UButton
      >
      <UButton to="/accounting/accounts" color="info" variant="outline" block size="lg">
        <UIcon name="i-lucide-book-open" class="mr-2" />
        {{ t("accounting.chart_of_accounts") }}
      </UButton>
    </div>

    <!-- Recent Transactions -->
    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <!-- Recent Income -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">{{ t("accounting.recent_income") }}</h3>
            <UButton variant="ghost" size="sm" to="/accounting/income">
              {{ t("accounting.view_all") }}
            </UButton>
          </div>
        </template>

        <div v-if="recentIncome.length === 0" class="py-8 text-center text-gray-500">
          No recent income recorded
        </div>
        <div v-else class="space-y-3">
          <div
            v-for="item in recentIncome"
            :key="item.id"
            class="flex items-center justify-between border-b pb-3 last:border-0"
          >
            <div>
              <p class="font-medium">{{ item.description }}</p>
              <p class="text-sm text-gray-500">{{ formatDate(item.incomeDate) }}</p>
            </div>
            <div class="text-right">
              <p class="font-medium text-green-600">+${{ Number(item.amount).toFixed(2) }}</p>
              <UBadge :color="getPaymentStatusColor(item.paymentStatus)" variant="soft" size="sm">
                {{ item.paymentStatus }}
              </UBadge>
            </div>
          </div>
        </div>
      </UCard>

      <!-- Recent Expenses -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Recent Expenses</h3>
            <UButton variant="ghost" size="sm" to="/accounting/expenses">View All</UButton>
          </div>
        </template>

        <div v-if="recentExpenses.length === 0" class="py-8 text-center text-gray-500">
          No recent expenses recorded
        </div>
        <div v-else class="space-y-3">
          <div
            v-for="item in recentExpenses"
            :key="item.id"
            class="flex items-center justify-between border-b pb-3 last:border-0"
          >
            <div>
              <p class="font-medium">{{ item.description }}</p>
              <p class="text-sm text-gray-500">{{ formatDate(item.expenseDate) }}</p>
            </div>
            <div class="text-right">
              <p class="font-medium text-red-600">-${{ Number(item.amount).toFixed(2) }}</p>
              <UBadge :color="getExpenseStatusColor(item.paymentStatus)" variant="soft" size="sm">
                {{ item.paymentStatus }}
              </UBadge>
            </div>
          </div>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getApiAccountingIncome, getApiAccountingExpenses } from "~/utils/client";

const { t } = useI18n();

const totalIncome = ref(0);
const totalExpenses = ref(0);
const netBalance = ref(0);
const pendingPayments = ref(0);
const recentIncome = ref<any[]>([]);
const recentExpenses = ref<any[]>([]);

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
  });
};

const getPaymentStatusColor = (status: string): "success" | "warning" | "error" => {
  const colors: Record<string, "success" | "warning" | "error"> = {
    received: "success",
    pending: "warning",
    refunded: "error",
  };
  return colors[status] || "warning";
};

const getExpenseStatusColor = (status: string): "success" | "warning" | "error" => {
  const colors: Record<string, "success" | "warning" | "error"> = {
    paid: "success",
    pending: "warning",
    cancelled: "error",
  };
  return colors[status] || "warning";
};

const fetchDashboardData = async () => {
  try {
    const incomeResponse = await getApiAccountingIncome({ query: { limit: 5 } });
    recentIncome.value = incomeResponse.data?.data ?? [];

    const expenseResponse = await getApiAccountingExpenses({ query: { limit: 5 } });
    recentExpenses.value = expenseResponse.data?.data ?? [];

    const allIncome = await getApiAccountingIncome({ query: { limit: 1000 } });
    const allExpenses = await getApiAccountingExpenses({ query: { limit: 1000 } });

    const incomeData = allIncome.data?.data ?? [];
    const expenseData = allExpenses.data?.data ?? [];

    totalIncome.value = incomeData.reduce((sum, item) => sum + Number(item.amount), 0);
    totalExpenses.value = expenseData.reduce((sum, item) => sum + Number(item.amount), 0);
    netBalance.value = totalIncome.value - totalExpenses.value;

    const pendingIncome = incomeData
      .filter((i) => i.paymentStatus === "pending")
      .reduce((sum, item) => sum + Number(item.amount), 0);
    const pendingExpenses = expenseData
      .filter((e) => e.paymentStatus === "pending")
      .reduce((sum, item) => sum + Number(item.amount), 0);
    pendingPayments.value = pendingIncome + pendingExpenses;
  } catch (error) {
    console.error("Failed to fetch dashboard data:", error);
  }
};

onMounted(fetchDashboardData);
</script>
