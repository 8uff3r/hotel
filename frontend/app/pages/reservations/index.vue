<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("reservations.reservations") }}
      </h1>
      <UButton to="/reservations/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />{{ t("reservations.new_reservation") }}</UButton
      >
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
        />
        <USelect
          v-model="filters.paymentStatus"
          :items="paymentStatusOptions"
          :placeholder="t('common.all_payments')"
          class="w-full sm:w-40"
        />
        <div class="flex gap-2">
          <UInput
            v-model="filters.entryDate"
            type="date"
            :placeholder="t('reservations.check_in_from')"
            class="w-full sm:w-40"
          />
          <UInput
            v-model="filters.departureDate"
            type="date"
            :placeholder="t('reservations.check_in_to')"
            class="w-full sm:w-40"
          />
        </div>
        <UButton variant="outline" @click="clearFilters">{{ t("common.clear") }}</UButton>
      </div>
    </UCard>

    <!-- Reservations Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("reservations.reservation_list") }}</span>
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
            <p class="font-medium">{{ row.original.rooms?.map((v) => v.roomNumber).join("،") }}</p>
          </div>
        </template>

        <template #dates-cell="{ row }">
          <div>
            <p class="text-sm">
              <UIcon name="i-lucide-log-in" class="mr-1 inline h-3 w-3" />
              {{ formatDate(row.original.entryDate) }}
            </p>
            <p class="text-sm">
              <UIcon name="i-lucide-log-out" class="mr-1 inline h-3 w-3" />
              {{ formatDate(row.original.departureDate) }}
            </p>
          </div>
        </template>

        <template #status-cell="{ row }">
          <UBadge :style="{ backgroundColor: `#${row.original.status?.colorHex}` }" variant="soft">
            {{ row.original.status }}
          </UBadge>
        </template>

        <template #paymentStatus-cell="{ row }">
          <UBadge
            :style="{ backgroundColor: `#${row.original.payment?.status?.colorHex}` }"
            variant="soft"
          >
            {{ row.original.payment?.status?.label }}
          </UBadge>
        </template>

        <template #totalAmount-cell="{ row }">
          ${{ row.original.payment?.amount?.toFixed(2) }}
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/reservations/${row.original.id}`">
              <UIcon name="i-lucide-eye" class="h-4 w-4" />
            </UButton>
            <UButton
              v-if="row.original.status?.slug === 'confirmed'"
              variant="ghost"
              size="sm"
              color="success"
              @click="checkIn(row.original)"
            >
              <UIcon name="i-lucide-log-in" class="h-4 w-4" />
            </UButton>
            <UButton
              v-if="row.original.status?.slug === 'checked_in'"
              variant="ghost"
              size="sm"
              color="warning"
              @click="checkOut(row.original)"
            >
              <UIcon name="i-lucide-log-out" class="h-4 w-4" />
            </UButton>
            <UButton
              v-if="row.original.status?.slug === 'confirmed'"
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
          <UPagination v-model="page" :page-count="pagination.limit" :total="pagination.total" />
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import type { Reservation } from "~/utils/client";

definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

type ReservationRow = Reservation;

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

const { t } = useI18n();
const page = ref(1);

const filters = reactive({
  search: "",
  status: "",
  paymentStatus: "",
  entryDate: "",
  departureDate: "",
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
    refetch();
  }, 300);
};

const {
  data: reservations,
  isLoading: loading,
  refetch,
} = useQuery({
  key: ["reservations", "list", pagination],
  query: async () => {
    const response = await getApiReservation({
      query: {
        page: pagination.page,
        limit: pagination.limit,
        filters: buildFilters({
          entry_date: { op: "eq", value: filters.entryDate },
          payment_status: { op: "eq", value: filters.paymentStatus },
          departure_date: { op: "eq", value: filters.departureDate },
          status: { op: "eq", value: filters.status },
        }),
      },
    });
    reservations.value = response.data?.data as any;
    pagination.total = response.data?.total ?? 0;
    pagination.totalPages = response?.data?.totalPages ?? 0;
    return response.data?.data;
  },
});

const clearFilters = () => {
  filters.search = "";
  filters.status = "";
  filters.paymentStatus = "";
  filters.entryDate = "";
  filters.departureDate = "";
  pagination.page = 1;
  refetch();
};

const formatDate = (date: Date | string | undefined) => {
  if (!date) return "";
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
    refetch();
  } catch (error) {
    console.error("Failed to check in:", error);
  }
};

const checkOut = async (reservation: ReservationRow) => {
  try {
    await $fetch(`/api/reservations/${reservation.id}/check-out`, { method: "POST" });
    refetch();
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
    refetch();
  } catch (error) {
    console.error("Failed to cancel reservation:", error);
  }
};
</script>
