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
          {{ t("parking.parking_lots") }}
        </h1>
      </div>
      <UButton to="/parking/lots/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />{{ t("parking.add_parking_lot") }}</UButton
      >
    </div>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("parking.parking_lots") }}</span>
          <span class="text-sm text-gray-500">{{ t("parking.lots_count", { count: pagination.total }) }}</span>
        </div>
      </template>

      <UTable :data="parkingLots" :columns="columns" :loading="loading" striped>
        <template #name-cell="{ row }">
          <NuxtLink
            :to="`/parking/lots/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            {{ row.original.name }}
          </NuxtLink>
        </template>

        <template #totalSpots-cell="{ row }">
          {{ t("parking.spots_count", { count: row.original.totalSpots }) }}
        </template>

        <template #hourlyRate-cell="{ row }"> ${{ row.original.hourlyRate }}/hr </template>

        <template #dailyRate-cell="{ row }"> ${{ row.original.dailyRate }}/day </template>

        <template #status-cell="{ row }">
          <UBadge :color="getStatusColor(row.original.status) as any" variant="soft">
            {{ row.original.status }}
          </UBadge>
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/parking/lots/${row.original.id}`">
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
            {{ t("pagination.pageOf", { page: pagination.page, totalPages: pagination.totalPages }) }}
          </span>
          <UPagination
            v-model="page"
            :page-count="pagination.limit"
            :total="pagination.total"
            @change="fetchParkingLots"
          />
        </div>
      </template>
    </UCard>

    <UModal v-model="deleteModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("actions.confirmDelete") }}</h2>
      </template>
      <template #body>
        <p>{{ t("parking.confirm_delete_lot", { name: selectedLot?.name }) }}</p>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="deleteModalOpen = false">{{ t("actions.cancel") }}</UButton>
          <UButton color="error" :loading="deleting" @click="deleteLot">{{ t("actions.delete") }}</UButton>
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

interface ParkingLot {
  id: number;
  name: string;
  location: string | null;
  totalSpots: number;
  hourlyRate: string;
  dailyRate: string;
  status: string;
}

const parkingLots = ref<ParkingLot[]>([]);
const { t } = useI18n();
const columns: TableColumn<ParkingLot>[] = [
  { accessorKey: "name", header: t("parking.name") },
  { accessorKey: "location", header: t("parking.location") },
  { accessorKey: "totalSpots", header: t("parking.capacity") },
  { accessorKey: "hourlyRate", header: t("parking.hourly_rate") },
  { accessorKey: "dailyRate", header: t("parking.daily_rate") },
  { accessorKey: "status", header: t("common.status") },
  { accessorKey: "actions", header: t("parking.actions") },
];
const loading = ref(false);
const deleting = ref(false);
const deleteModalOpen = ref(false);
const selectedLot = ref<ParkingLot | null>(null);
const page = ref(1);

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  totalPages: 0,
});

const fetchParkingLots = async () => {
  loading.value = true;
  try {
    const params = new URLSearchParams();
    params.append("page", pagination.page.toString());
    params.append("limit", pagination.limit.toString());

    const response = await $fetch(`/api/parking/lots?${params.toString()}`);
    parkingLots.value = response.data;
    pagination.total = response.pagination.total ?? 0;
    pagination.totalPages = response.pagination.totalPages ?? 0;
  } catch (error) {
    console.error("Failed to fetch parking lots:", error);
  } finally {
    loading.value = false;
  }
};

const confirmDelete = (lot: ParkingLot) => {
  selectedLot.value = lot;
  deleteModalOpen.value = true;
};

const deleteLot = async () => {
  if (!selectedLot.value) return;

  deleting.value = true;
  try {
    await $fetch(`/api/parking/lots/${selectedLot.value.id}`, { method: "DELETE" });
    deleteModalOpen.value = false;
    await fetchParkingLots();
  } catch (error) {
    console.error("Failed to delete parking lot:", error);
  } finally {
    deleting.value = false;
  }
};

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    active: "success",
    full: "warning",
    closed: "error",
  };
  return colors[status] || "neutral";
};

onMounted(fetchParkingLots);
</script>
