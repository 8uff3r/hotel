<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("reservations.reservations") }}
      </h1>
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
          <span class="text-sm text-gray-500">{{
            t("reservations.n_reservations", { count: pagination.total })
          }}</span>
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
            {{ row.original.status?.label }}
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
              v-if="['awaiting_payment', 'verified'].includes(row.original.status?.slug ?? '')"
              variant="ghost"
              size="sm"
              color="primary"
              @click="acceptReservation(row.original)"
            >
              {{ t("reservations.accept") }}
            </UButton>
            <UButton
              v-if="row.original.status?.slug === 'accepted'"
              variant="ghost"
              size="sm"
              color="success"
              @click="checkIn(row.original)"
            >
              <UIcon name="i-lucide-log-in" class="h-4 w-4" />
            </UButton>
            <UButton
              v-if="row.original.status?.slug === 'accepted'"
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
            {{
              t("reservations.page_of", {
                page: pagination.page,
                totalPages: pagination.totalPages,
              })
            }}
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
import {
  getApiReservation,
  getApiAccountingPaymentStatuses,
  postApiReservationIdAccept,
  postApiReservationIdCheckIn,
  postApiReservationIdCheckOut,
  putApiReservationId,
} from "~/utils/client";

type ReservationRow = Reservation;

const columns: TableColumn<ReservationRow>[] = [
  { accessorKey: "id", header: () => t("reservations.id") },
  { accessorKey: "guest", header: () => t("reservations.guest") },
  { accessorKey: "room", header: () => t("reservations.room") },
  { accessorKey: "dates", header: () => t("reservations.dates") },
  { accessorKey: "status", header: () => t("reservations.status") },
  { accessorKey: "paymentStatus", header: () => t("reservations.payment") },
  { accessorKey: "totalAmount", header: () => t("reservations.total") },
  { accessorKey: "actions", header: () => t("reservations.actions") },
];

const { t } = useI18n();
const page = ref(1);

const filters = reactive<{
  search: string;
  status: string | undefined;
  paymentStatus: string | undefined;
  entryDate: string;
  departureDate: string;
}>({
  search: "",
  status: undefined,
  paymentStatus: undefined,
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

const { data: reservationStatuses } = useQuery({
  key: ["reservation", "statuses"],
  query: async () => {
    const res = await getApiReservationStatuses({ query: { limit: -1 } });
    return res.data?.data;
  },
});

const { data: paymentStatuses } = useQuery({
  key: ["payment", "statuses"],
  query: async () => {
    const res = await getApiAccountingPaymentStatuses();
    return res.data?.data;
  },
});

const statusOptions = computed(() => [
  { value: "all", label: t("common.all_statuses") },
  ...(reservationStatuses.value?.map((s) => ({ value: s.slug, label: s.label })) ?? []),
]);

const paymentStatusOptions = computed(() => [
  { value: "all", label: t("common.all_payments") },
  ...(paymentStatuses.value?.map((s) => ({ value: s.slug, label: s.label })) ?? []),
]);

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
          payment_status: {
            op: "eq",
            value: filters.paymentStatus != "all" ? filters.paymentStatus : undefined,
          },
          departure_date: { op: "eq", value: filters.departureDate },
          status: { op: "eq", value: filters.status != "all" ? filters.status : undefined },
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
  filters.status = undefined;
  filters.paymentStatus = undefined;
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

const acceptReservation = async (reservation: ReservationRow) => {
  if (!reservation.id) return;
  try {
    await postApiReservationIdAccept({ path: { id: String(reservation.id) } });
    refetch();
  } catch (error) {
    console.error("Failed to accept reservation:", error);
  }
};

const checkIn = async (reservation: ReservationRow) => {
  if (!reservation.id) return;
  try {
    await postApiReservationIdCheckIn({ path: { id: String(reservation.id) } });
    refetch();
  } catch (error) {
    console.error("Failed to check in:", error);
  }
};

const checkOut = async (reservation: ReservationRow) => {
  if (!reservation.id) return;
  try {
    await postApiReservationIdCheckOut({ path: { id: String(reservation.id) } });
    refetch();
  } catch (error) {
    console.error("Failed to check out:", error);
  }
};

const cancelReservation = async (reservation: ReservationRow) => {
  if (!reservation.id) return;
  try {
    await putApiReservationId({
      path: { id: String(reservation.id) },
      body: { status: { slug: "cancelled" } },
    });
    refetch();
  } catch (error) {
    console.error("Failed to cancel:", error);
  }
};
</script>
