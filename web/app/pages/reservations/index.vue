<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('reservations.reservations') }}</h1>
      <UButton to="/reservations/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />{{ t('reservations.new_reservation') }}</UButton>
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          :placeholder="t('reservations.search_by_guest_name')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.status"
          :items="statusOptions"
          :placeholder="t('common.all_statuses')"
          class="w-full sm:w-40"
          @change="fetchReservations"
        />
        <USelect
          v-model="filters.paymentStatus"
          :items="paymentStatusOptions"
          :placeholder="t('common.all_payments')"
          class="w-full sm:w-40"
          @change="fetchReservations"
        />
        <div class="flex gap-2">
          <UInput
            v-model="filters.checkInFrom"
            type="date"
            :placeholder="t('reservations.check_in_from')"
            class="w-full sm:w-40"
            @change="fetchReservations"
          />
          <UInput
            v-model="filters.checkInTo"
            type="date"
            :placeholder="t('reservations.check_in_to')"
            class="w-full sm:w-40"
            @change="fetchReservations"
          />
        </div>
        <UButton variant="outline" @click="clearFilters">{{ t('common.clear') }}</UButton>
      </div>
    </UCard>

    <!-- Reservations Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t('reservations.reservation_list') }}</span>
          <span class="text-sm text-gray-500">{{ pagination.total }} reservations</span>
        </div>
      </template>

      <UTable :data="reservations" :columns="columns" :loading="loading" striped>
        <template #id-cell="{ row }">
          <NuxtLink
            :to="`/reservations/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            #{{ row.original.id }}
          </NuxtLink>
        </template>

        <template #guest-cell="{ row }">
          <div>
            <p class="font-medium">
              {{ row.original.guest?.firstName }} {{ row.original.guest?.lastName }}
            </p>
            <p class="text-sm text-gray-500">{{ row.original.guest?.email }}</p>
          </div>
        </template>

        <template #room-cell="{ row }">
          <div>
            <p class="font-medium">{{ row.original.roomNumber }}</p>
            <p class="text-sm text-gray-500 capitalize">{{ row.original.roomType }}</p>
          </div>
        </template>

        <template #dates-cell="{ row }">
          <div>
            <p class="text-sm">
              <UIcon name="i-lucide-log-in" class="mr-1 inline h-3 w-3" />
              {{ formatDate(row.original.checkInDate) }}
            </p>
            <p class="text-sm">
              <UIcon name="i-lucide-log-out" class="mr-1 inline h-3 w-3" />
              {{ formatDate(row.original.checkOutDate) }}
            </p>
          </div>
        </template>

        <template #status-cell="{ row }">
          <UBadge :color="getStatusColor(row.original.status)" variant="soft">
            {{ row.original.status }}
          </UBadge>
        </template>

        <template #paymentStatus-cell="{ row }">
          <UBadge :color="getPaymentColor(row.original.paymentStatus)" variant="soft">
            {{ row.original.paymentStatus }}
          </UBadge>
        </template>

        <template #totalAmount-cell="{ row }">
          ${{ row.original.totalAmount?.toFixed(2) }}
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/reservations/${row.original.id}`">
              <UIcon name="i-lucide-eye" class="h-4 w-4" />
            </UButton>
            <UButton
              v-if="row.original.status === 'confirmed'"
              variant="ghost"
              size="sm"
              color="success"
              @click="checkIn(row.original)"
            >
              <UIcon name="i-lucide-log-in" class="h-4 w-4" />
            </UButton>
            <UButton
              v-if="row.original.status === 'checked_in'"
              variant="ghost"
              size="sm"
              color="warning"
              @click="checkOut(row.original)"
            >
              <UIcon name="i-lucide-log-out" class="h-4 w-4" />
            </UButton>
            <UButton
              v-if="row.original.status === 'confirmed'"
              variant="ghost"
              size="sm"
              color="error"
              @click="cancelReservation(row.original)"
            >
              <UIcon name="i-lucide-x" class="h-4 w-4" />
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
            @change="fetchReservations"
          />
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";

definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

interface ReservationRow {
  id: number;
  guestId: number;
  roomId: number;
  roomNumber: string | null;
  roomType: string | null;
  checkInDate: Date | string;
  checkOutDate: Date | string;
  status: string;
  totalAmount: number;
  paidAmount: number;
  paymentStatus: string;
  guest?: {
    firstName: string;
    lastName: string;
    email: string;
  };
}

const columns: TableColumn<ReservationRow>[] = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "guest", header: "Guest" },
  { accessorKey: "room", header: "Room" },
  { accessorKey: "dates", header: "Dates" },
  { accessorKey: "status", header: "Status" },
  { accessorKey: "paymentStatus", header: "Payment" },
  { accessorKey: "totalAmount", header: "Total" },
  { accessorKey: "actions", header: "Actions" },
];

const statusOptions = [
  { value: "all", label: "All Statuses" },
  { value: "pending", label: "Pending" },
  { value: "confirmed", label: "Confirmed" },
  { value: "checked_in", label: "Checked In" },
  { value: "checked_out", label: "Checked Out" },
  { value: "cancelled", label: "Cancelled" },
  { value: "no_show", label: "No Show" },
];

const paymentStatusOptions = [
  { value: "all", label: "All Payments" },
  { value: "pending", label: "Pending" },
  { value: "partial", label: "Partial" },
  { value: "paid", label: "Paid" },
  { value: "refunded", label: "Refunded" },
];

const reservations = ref<ReservationRow[]>([]);
const { t } = useI18n();
const loading = ref(false);
const page = ref(1);

const filters = reactive({
  search: "",
  status: "",
  paymentStatus: "",
  checkInFrom: "",
  checkInTo: "",
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
    fetchReservations();
  }, 300);
};

const fetchReservations = async () => {
  loading.value = true;
  try {
    const params: Record<string, any> = {
      page: pagination.page,
      limit: pagination.limit,
    };

    if (filters.status) params["status"] = filters.status;
    if (filters.paymentStatus) params["paymentStatus"] = filters.paymentStatus;
    if (filters.checkInFrom) params["checkInFrom"] = filters.checkInFrom;
    if (filters.checkInTo) params["checkInTo"] = filters.checkInTo;

    const response = await $fetch(`/api/reservations`, {
      query: params,
    });
    reservations.value = response.data as any;
    pagination.total = response.pagination.total;
    pagination.totalPages = response.pagination.totalPages;
  } catch (error) {
    console.error("Failed to fetch reservations:", error);
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  filters.search = "";
  filters.status = "";
  filters.paymentStatus = "";
  filters.checkInFrom = "";
  filters.checkInTo = "";
  pagination.page = 1;
  fetchReservations();
};

const formatDate = (date: Date | string) => {
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};

const getStatusColor = (status: string): "success" | "warning" | "info" | "error" | "neutral" => {
  const colors: Record<string, "success" | "warning" | "info" | "error" | "neutral"> = {
    pending: "warning",
    confirmed: "info",
    checked_in: "success",
    checked_out: "neutral",
    cancelled: "error",
    no_show: "error",
  };
  return colors[status] || "neutral";
};

const getPaymentColor = (status: string): "success" | "warning" | "info" | "error" | "neutral" => {
  const colors: Record<string, "success" | "warning" | "info" | "error" | "neutral"> = {
    pending: "warning",
    partial: "info",
    paid: "success",
    refunded: "error",
  };
  return colors[status] || "neutral";
};

const checkIn = async (reservation: ReservationRow) => {
  try {
    await $fetch(`/api/reservations/${reservation.id}/check-in`, { method: "POST" });
    await fetchReservations();
  } catch (error) {
    console.error("Failed to check in:", error);
  }
};

const checkOut = async (reservation: ReservationRow) => {
  try {
    await $fetch(`/api/reservations/${reservation.id}/check-out`, { method: "POST" });
    await fetchReservations();
  } catch (error) {
    console.error("Failed to check out:", error);
  }
};

const cancelReservation = async (reservation: ReservationRow) => {
  try {
    await $fetch(`/api/reservations/${reservation.id}`, {
      method: "PUT" as any,
      body: { status: "cancelled" },
    });
    await fetchReservations();
  } catch (error) {
    console.error("Failed to cancel reservation:", error);
  }
};

onMounted(fetchReservations);
</script>
