<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <UButton to="/parking" variant="ghost" size="sm" class="mb-2">
          <UIcon name="i-lucide-arrow-left" class="mr-1" />{{ t('parking.back_to_parking') }}</UButton>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('parking.parking_spots') }}</h1>
      </div>
      <UButton to="/parking/spots/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />{{ t('parking.add_spot') }}</UButton>
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
        <UButton variant="outline" @click="clearFilters">{{ t('common.clear') }}</UButton>
      </div>
    </UCard>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t('parking.parking_spots') }}</span>
          <span class="text-sm text-gray-500">{{ pagination.total }} spots</span>
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
          <UBadge :color="getSpotTypeColor(row.original.spotType) as any" variant="soft">
            {{ row.original.spotType }}
          </UBadge>
        </template>

        <template #status-cell="{ row }">
          <UBadge :color="getStatusColor(row.original.status) as any" variant="soft">
            {{ row.original.status }}
          </UBadge>
        </template>

        <template #isCovered-cell="{ row }">
          {{ row.original.isCovered ? "Yes" : "No" }}
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
            Page {{ pagination.page }} of {{ pagination.totalPages }}
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
        <h2 class="text-lg font-semibold">Confirm Delete</h2>
      </template>
      <template #body>
        <p>Are you sure you want to delete spot "{{ selectedSpot?.spotNumber }}"?</p>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="deleteModalOpen = false">Cancel</UButton>
          <UButton color="error" :loading="deleting" @click="deleteSpot">Delete</UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";

definePageMeta({
  requiresRole: ["admin", "manager"],
});

interface Spot {
  id: number;
  lotId: number | null;
  spotNumber: string;
  floor: string | null;
  spotType: string;
  status: string;
  isCovered: boolean;
}

const columns: TableColumn<Spot>[] = [
  { accessorKey: "spotNumber", header: "Spot #" },
  { accessorKey: "lotId", header: "Lot" },
  { accessorKey: "floor", header: "Floor" },
  { accessorKey: "spotType", header: "Type" },
  { accessorKey: "isCovered", header: "Covered" },
  { accessorKey: "status", header: "Status" },
  { accessorKey: "actions", header: "Actions" },
];

const spots = ref<Spot[]>([]);
const lots = ref<any[]>([]);
const { t } = useI18n();
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
  { value: "all", label: "All Statuses" },
  { value: "available", label: "Available" },
  { value: "occupied", label: "Occupied" },
  { value: "reserved", label: "Reserved" },
  { value: "maintenance", label: "Maintenance" },
];

const lotOptions = ref<{ value: string; label: string }[]>([]);

const fetchLots = async () => {
  try {
    const res = await $fetch("/api/parking/lots");
    lots.value = res.data;
    lotOptions.value = [
      { value: "all", label: "All Lots" },
      ...lots.value.map((l: any) => ({ value: l.id.toString(), label: l.name })),
    ];
  } catch (error) {
    console.error("Failed to fetch lots:", error);
  }
};

const getLotName = (lotId: number | null) => {
  if (!lotId) return "-";
  const lot = lots.value.find((l) => l.id === lotId);
  return lot?.name || "-";
};

const fetchSpots = async () => {
  loading.value = true;
  try {
    const params = new URLSearchParams();
    params.append("page", pagination.page.toString());
    params.append("limit", pagination.limit.toString());

    if (filters.lotId && filters.lotId !== "all") params.append("lotId", filters.lotId);
    if (filters.status && filters.status !== "all") params.append("status", filters.status);

    const response = await $fetch(`/api/parking/spots?${params.toString()}`);
    spots.value = response.data;
    pagination.total = response.pagination.total ?? 0;
    pagination.totalPages = response.pagination.totalPages ?? 0;
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
    await $fetch(`/api/parking/spots/${selectedSpot.value.id}`, { method: "DELETE" });
    deleteModalOpen.value = false;
    await fetchSpots();
  } catch (error) {
    console.error("Failed to delete spot:", error);
  } finally {
    deleting.value = false;
  }
};

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    available: "success",
    occupied: "warning",
    reserved: "info",
    maintenance: "error",
  };
  return colors[status] || "neutral";
};

const getSpotTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    standard: "neutral",
    handicap: "primary",
    electric: "success",
    compact: "info",
    large: "warning",
  };
  return colors[type] || "neutral";
};

onMounted(() => {
  fetchLots();
  fetchSpots();
});
</script>
