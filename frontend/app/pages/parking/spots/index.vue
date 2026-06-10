<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <UButton to="/parking" variant="ghost" size="sm" class="mb-2">
          <UIcon name="i-lucide-arrow-left" class="mr-1" />{{
            t("parking.back_to_parking")
          }}</UButton
        >
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          {{ t("parking.parking_spots") }}
        </h1>
      </div>
      <UButton to="/parking/spots/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />{{ t("parking.add_spot") }}</UButton
      >
    </div>

    <UCard class="mb-4">
      <div class="flex flex-wrap items-center gap-4">
        <USelect
          v-model="filters.lotId"
          :items="lotOptions"
          :placeholder="t('parking.all_lots')"
          class="w-full sm:w-48"
          @change="fetchSpots"
        />
        <USelect
          v-model="filters.status"
          :items="statusOptions"
          :placeholder="t('common.all_statuses')"
          class="w-full sm:w-40"
          @change="fetchSpots"
        />
        <UButton variant="outline" @click="clearFilters">{{ t("common.clear") }}</UButton>
      </div>
    </UCard>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("parking.parking_spots") }}</span>
          <span class="text-sm text-gray-500">{{
            t("parking.spots_count", { count: pagination.total })
          }}</span>
        </div>
      </template>

      <UTable :data="spots" :columns="columns" :loading="loading" striped>
        <template #spotNumber-cell="{ row }">
          <NuxtLink
            :to="`/parking/spots/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            {{ row.original.spotNumber }}
          </NuxtLink>
        </template>

        <template #lotId-cell="{ row }">
          {{ getLotName(row.original.lotId) }}
        </template>

        <template #spotType-cell="{ row }">
          <UBadge :color="getSpotTypeColor(row.original.spotType?.slug) as any" variant="soft">
            {{ row.original.spotType?.label }}
          </UBadge>
        </template>

        <template #status-cell="{ row }">
          <UBadge :color="getStatusColor(row.original.status?.slug) as any" variant="soft">
            {{ row.original.status?.label }}
          </UBadge>
        </template>

        <template #isCovered-cell="{ row }">
          {{ row.original.isCovered ? t("common.yes") : t("common.no") }}
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/parking/spots/${row.original.id}`">
              <UIcon name="i-lucide-eye" class="h-4 w-4" />
            </UButton>
            <UButton variant="ghost" size="sm" color="error" @click="confirmDelete(row.original)">
              <UIcon name="i-lucide-trash-2" class="h-4 w-4" />
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
          <UPagination
            v-model="page"
            :page-count="pagination.limit"
            :total="pagination.total"
            @change="fetchSpots"
          />
        </div>
      </template>
    </UCard>

    <UModal v-model="deleteModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("actions.confirmDelete") }}</h2>
      </template>
      <template #body>
        <p>{{ t("parking.confirm_delete_spot", { spot: selectedSpot?.spotNumber }) }}</p>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="deleteModalOpen = false">{{
            t("actions.cancel")
          }}</UButton>
          <UButton color="error" :loading="deleting" @click="deleteSpot">{{
            t("actions.delete")
          }}</UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import { getApiParkingLots, getApiParkingSpots, deleteApiParkingSpotsId } from "~/utils/client";
import type { ParkingLot, PaginatedResponseModelsParkingSpot } from "~/utils/client";

interface Spot {
  id?: number;
  lotId?: number | null;
  spotNumber?: string;
  floor?: string | null;
  spotType?: { id?: number; label?: string; slug?: string };
  status?: { colorHex?: string; id?: number; label?: string; slug?: string };
  isCovered?: boolean;
}

const spots = ref<Spot[]>([]);
const lots = ref<any[]>([]);
const { t } = useI18n();
const columns: TableColumn<Spot>[] = [
  { accessorKey: "spotNumber", header: t("parking.spot_number") },
  { accessorKey: "lotId", header: t("parking.parking_lot") },
  { accessorKey: "floor", header: t("common.floor") },
  { accessorKey: "spotType", header: t("parking.spot_type") },
  { accessorKey: "isCovered", header: t("parking.covered") },
  { accessorKey: "status", header: t("common.status") },
  { accessorKey: "actions", header: t("parking.actions") },
];
const loading = ref(false);
const deleting = ref(false);
const deleteModalOpen = ref(false);
const selectedSpot = ref<Spot | null>(null);
const page = ref(1);

const filters = reactive({
  lotId: "",
  status: "",
});

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  totalPages: 0,
});

const statusOptions = [
  { value: "all", label: t("common.all_statuses") },
  { value: "available", label: t("parking.status_available") },
  { value: "occupied", label: t("parking.status_occupied") },
  { value: "reserved", label: t("parking.status_reserved") },
  { value: "maintenance", label: t("parking.status_maintenance") },
];

const lotOptions = ref<{ value: string; label: string }[]>([]);

const fetchLots = async () => {
  try {
    const res = await getApiParkingLots();
    lots.value = res.data?.data ?? [];
    lotOptions.value = [
      { value: "all", label: t("parking.all_lots") },
      ...lots.value.map((l) => ({ value: String(l.id ?? ""), label: l.name ?? "" })),
    ];
  } catch (error) {
    console.error("Failed to fetch lots:", error);
  }
};

const getLotName = (lotId: number | null | undefined) => {
  if (!lotId) return "-";
  const lot = lots.value.find((l) => l.id === lotId);
  return lot?.name || "-";
};

const fetchSpots = async () => {
  loading.value = true;
  try {
    const query: Record<string, any> = {
      page: pagination.page.toString(),
      limit: pagination.limit.toString(),
    };
    if (filters.lotId && filters.lotId !== "all") query.lotId = filters.lotId;
    if (filters.status && filters.status !== "all") query.status = filters.status;

    const response = await getApiParkingSpots({ query });
    spots.value = response.data?.data ?? [];
    pagination.total = response.data?.total ?? 0;
    pagination.totalPages = response.data?.totalPages ?? 0;
  } catch (error) {
    console.error("Failed to fetch spots:", error);
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  filters.lotId = "";
  filters.status = "";
  pagination.page = 1;
  fetchSpots();
};

const confirmDelete = (spot: Spot) => {
  selectedSpot.value = spot;
  deleteModalOpen.value = true;
};

const deleteSpot = async () => {
  if (!selectedSpot.value) return;

  deleting.value = true;
  try {
    await deleteApiParkingSpotsId({ path: { id: String(selectedSpot.value.id) } });
    deleteModalOpen.value = false;
    await fetchSpots();
  } catch (error) {
    console.error("Failed to delete spot:", error);
  } finally {
    deleting.value = false;
  }
};

const getStatusColor = (status: string | undefined) => {
  const colors: Record<string, string> = {
    available: "success",
    occupied: "warning",
    reserved: "info",
    maintenance: "error",
  };
  return colors[status ?? ""] || "neutral";
};

const getSpotTypeColor = (type: string | undefined) => {
  const colors: Record<string, string> = {
    standard: "neutral",
    handicap: "primary",
    electric: "success",
    compact: "info",
    large: "warning",
  };
  return colors[type ?? ""] || "neutral";
};

onMounted(() => {
  fetchLots();
  fetchSpots();
});
</script>
