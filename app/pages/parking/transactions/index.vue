<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <UButton to="/parking" variant="ghost" size="sm" class="mb-2">
          <UIcon name="i-lucide-arrow-left" class="mr-1" />
          Back to Parking
        </UButton>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Parking Transactions</h1>
      </div>
      <UButton to="/parking/transactions/check-in" color="primary">
        <UIcon name="i-lucide-car" class="mr-2" />
        Check In
      </UButton>
    </div>

    <UCard class="mb-4">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          placeholder="Search by license plate..."
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.status"
          :items="statusOptions"
          placeholder="All Statuses"
          class="w-full sm:w-40"
          @change="fetchTransactions"
        />
        <USelect
          v-model="filters.paymentStatus"
          :items="paymentStatusOptions"
          placeholder="All Payments"
          class="w-full sm:w-40"
          @change="fetchTransactions"
        />
        <UButton variant="outline" @click="clearFilters"> Clear </UButton>
      </div>
    </UCard>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">Transactions</span>
          <span class="text-sm text-gray-500">{{ pagination.total }} transactions</span>
        </div>
      </template>

      <UTable :data="transactions" :columns="columns" :loading="loading" striped>
        <template #licensePlate-cell="{ row }">
          <NuxtLink
            :to="`/parking/transactions/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            {{ row.original.licensePlate }}
          </NuxtLink>
        </template>

        <template #entryTime-cell="{ row }">
          {{ formatDate(row.original.entryTime) }}
        </template>

        <template #exitTime-cell="{ row }">
          {{ row.original.exitTime ? formatDate(row.original.exitTime) : "-" }}
        </template>

        <template #hoursParked-cell="{ row }">
          {{ row.original.hoursParked ? row.original.hoursParked.toFixed(1) + " hrs" : "-" }}
        </template>

        <template #amountDue-cell="{ row }">
          ${{ row.original.amountDue?.toFixed(2) || "0.00" }}
        </template>

        <template #status-cell="{ row }">
          <UBadge :color="getStatusColor(row.original.status) as any" variant="soft">
            {{ row.original.status }}
          </UBadge>
        </template>

        <template #paymentStatus-cell="{ row }">
          <UBadge :color="getPaymentColor(row.original.paymentStatus) as any" variant="soft">
            {{ row.original.paymentStatus }}
          </UBadge>
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton
              v-if="row.original.status === 'active'"
              variant="ghost"
              size="sm"
              color="warning"
              @click="checkOut(row.original)"
            >
              <UIcon name="i-lucide-log-out" class="h-4 w-4" />
            </UButton>
            <UButton variant="ghost" size="sm" :to="`/parking/transactions/${row.original.id}`">
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
            @change="fetchTransactions"
          />
        </div>
      </template>
    </UCard>

    <UModal v-model="checkoutModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">Check Out Vehicle</h2>
      </template>
      <template #body>
        <div v-if="selectedTx" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <span class="text-sm text-gray-500">License Plate</span>
              <div class="font-medium">{{ selectedTx.licensePlate }}</div>
            </div>
            <div>
              <span class="text-sm text-gray-500">Entry Time</span>
              <div class="font-medium">{{ formatDate(selectedTx.entryTime) }}</div>
            </div>
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Payment Method</label>
            <USelect v-model="checkoutForm.paymentMethod" :items="paymentMethods" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Amount Paid</label>
            <UInput v-model="checkoutForm.amountPaid" type="number" step="0.01" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Notes</label>
            <UTextarea v-model="checkoutForm.notes" :rows="2" />
          </div>
        </div>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="checkoutModalOpen = false">Cancel</UButton>
          <UButton color="warning" :loading="checkingOut" @click="confirmCheckout"
            >Check Out</UButton
          >
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";

definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

interface Transaction {
  id: number;
  licensePlate: string;
  entryTime: string;
  exitTime: string | null;
  hoursParked: number | null;
  amountDue: number;
  amountPaid: number;
  status: string;
  paymentStatus: string;
}

const columns: TableColumn<Transaction>[] = [
  { accessorKey: "licensePlate", header: "License Plate" },
  { accessorKey: "entryTime", header: "Entry Time" },
  { accessorKey: "exitTime", header: "Exit Time" },
  { accessorKey: "hoursParked", header: "Duration" },
  { accessorKey: "amountDue", header: "Amount Due" },
  { accessorKey: "status", header: "Status" },
  { accessorKey: "paymentStatus", header: "Payment" },
  { accessorKey: "actions", header: "Actions" },
];

const transactions = ref<Transaction[]>([]);
const loading = ref(false);
const checkoutModalOpen = ref(false);
const selectedTx = ref<Transaction | null>(null);
const checkingOut = ref(false);
const page = ref(1);

const filters = reactive({
  search: "",
  status: "",
  paymentStatus: "",
});

const checkoutForm = reactive({
  paymentMethod: "cash",
  amountPaid: "",
  notes: "",
});

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  totalPages: 0,
});

const statusOptions = [
  { value: "all", label: "All Statuses" },
  { value: "active", label: "Active" },
  { value: "completed", label: "Completed" },
  { value: "cancelled", label: "Cancelled" },
];

const paymentStatusOptions = [
  { value: "all", label: "All Payments" },
  { value: "pending", label: "Pending" },
  { value: "paid", label: "Paid" },
  { value: "waived", label: "Waived" },
];

const paymentMethods = [
  { value: "cash", label: "Cash" },
  { value: "card", label: "Card" },
  { value: "guest_account", label: "Guest Account" },
];

let searchTimeout: ReturnType<typeof setTimeout> | null = null;

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    pagination.page = 1;
    fetchTransactions();
  }, 300);
};

const fetchTransactions = async () => {
  loading.value = true;
  try {
    const params = new URLSearchParams();
    params.append("page", pagination.page.toString());
    params.append("limit", pagination.limit.toString());

    if (filters.search) params.append("search", filters.search);
    if (filters.status && filters.status !== "all") params.append("status", filters.status);
    if (filters.paymentStatus && filters.paymentStatus !== "all")
      params.append("paymentStatus", filters.paymentStatus);

    const response = await $fetch(`/api/parking/transactions?${params.toString()}`);
    transactions.value = response.data;
    pagination.total = response.pagination.total ?? 0;
    pagination.totalPages = response.pagination.totalPages ?? 0;
  } catch (error) {
    console.error("Failed to fetch transactions:", error);
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  filters.search = "";
  filters.status = "";
  filters.paymentStatus = "";
  pagination.page = 1;
  fetchTransactions();
};

const checkOut = (tx: Transaction) => {
  selectedTx.value = tx;
  checkoutForm.amountPaid = tx.amountDue?.toString() || "0";
  checkoutForm.paymentMethod = "cash";
  checkoutForm.notes = "";
  checkoutModalOpen.value = true;
};

const confirmCheckout = async () => {
  if (!selectedTx.value) return;

  checkingOut.value = true;
  try {
    await $fetch(`/api/parking/transactions/${selectedTx.value.id}/check-out`, {
      method: "POST",
      body: {
        rateType: "hourly",
        amountPaid: parseFloat(checkoutForm.amountPaid) || 0,
        paymentMethod: checkoutForm.paymentMethod,
        notes: checkoutForm.notes || null,
      },
    });
    checkoutModalOpen.value = false;
    await fetchTransactions();
  } catch (error) {
    console.error("Failed to check out:", error);
  } finally {
    checkingOut.value = false;
  }
};

const formatDate = (date: string) => {
  return new Date(date).toLocaleString();
};

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    active: "info",
    completed: "success",
    cancelled: "error",
  };
  return colors[status] || "neutral";
};

const getPaymentColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: "warning",
    paid: "success",
    waived: "info",
  };
  return colors[status] || "neutral";
};

onMounted(fetchTransactions);
</script>
