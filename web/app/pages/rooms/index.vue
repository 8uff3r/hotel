<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("rooms.title") }}</h1>
      <UButton to="/rooms/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />
        {{ t("rooms.addRoom") }}
      </UButton>
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          :placeholder="t('rooms.searchPlaceholder')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.status"
          :items="statusOptions"
          :placeholder="t('rooms.allStatuses')"
          class="w-full sm:w-40"
        />
        <USelect
          v-model="filters.roomType"
          :items="roomTypeOptions"
          :placeholder="t('rooms.allTypes')"
          class="w-full sm:w-40"
        />
        <UButton variant="outline" @click="clearFilters"> {{ t("actions.clear") }} </UButton>
      </div>
    </UCard>

    <!-- Rooms Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("rooms.list") }}</span>
          <span class="text-sm text-gray-500">{{
            t("rooms.count", { count: pagination.total })
          }}</span>
        </div>
      </template>

      <UTable :data="rooms" :columns="columns" :loading="pending" striped>
        <template #roomNumber-cell="{ row }">
          <NuxtLink
            :to="`/rooms/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            {{ row.original.roomNumber }}
          </NuxtLink>
        </template>

        <template #roomType-cell="{ row }">
          <UBadge :color="getRoomTypeColor(row.original.roomTypeId) as any" variant="soft">
            {{ row.original.roomType }}
          </UBadge>
        </template>

        <template #status-cell="{ row }">
          <UBadge :color="getStatusColor(row.original.statusId) as any" variant="soft">
            {{ row.original.status }}
          </UBadge>
        </template>

        <template #basePrice-cell="{ row }"> ${{ row.original.basePrice?.toFixed(2) }} </template>

        <template #floor-cell="{ row }">
          {{ row.original.floor ?? "-" }}
        </template>

        <template #capacity-cell="{ row }">
          {{ t("rooms.capacityValue", { count: row.original.capacity }) }}
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/rooms/${row.original.id}`">
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
          <UPagination v-model="page" :page-count="pagination.limit" :total="pagination.total" />
        </div>
      </template>
    </UCard>

    <!-- Delete Confirmation Modal -->
    <UModal v-model="deleteModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("actions.confirmDelete") }}</h2>
      </template>
      <template #body>
        <p>{{ t("rooms.confirmDelete", { roomNumber: selectedRoom?.roomNumber }) }}</p>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="deleteModalOpen = false">{{
            t("actions.cancel")
          }}</UButton>
          <UButton color="error" :loading="deleting" @click="deleteRoom">{{
            t("actions.delete")
          }}</UButton>
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

const { t } = useI18n();
const columns = computed<TableColumn<Room>[]>(() => [
  { accessorKey: "roomNumber", header: t("rooms.columns.roomNumber") },
  { accessorKey: "roomType", header: t("rooms.columns.type") },
  { accessorKey: "floor", header: t("rooms.columns.floor") },
  { accessorKey: "capacity", header: t("rooms.columns.capacity") },
  { accessorKey: "basePrice", header: t("rooms.columns.price") },
  { accessorKey: "status", header: t("rooms.columns.status") },
  { accessorKey: "actions", header: t("rooms.columns.actions") },
]);

const statusOptions = computed(() => [
  { value: "all", label: t("rooms.allStatuses") },
  { value: "available", label: t("rooms.statuses.available") },
  { value: "occupied", label: t("rooms.statuses.occupied") },
  { value: "maintenance", label: t("rooms.statuses.maintenance") },
  { value: "out_of_order", label: t("rooms.statuses.outOfOrder") },
]);

const roomTypeOptions = computed(() => [
  { value: "all", label: t("rooms.allTypes") },
  { value: "single", label: t("rooms.types.single") },
  { value: "double", label: t("rooms.types.double") },
  { value: "suite", label: t("rooms.types.suite") },
  { value: "deluxe", label: t("rooms.types.deluxe") },
]);

const deleting = ref(false);
const deleteModalOpen = ref(false);
const selectedRoom = ref<Room | null>(null);
const page = ref(1);

const filters = reactive({
  search: "",
  status: "",
  roomType: "",
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
    // fetchRooms();
  }, 300);
};

const { data: rooms, pending } = useAsyncData(async () => {
  const response = await getApiRooms({
    query: computed(() => pagination),
  });
  pagination.total = response.total ?? 0;
  pagination.totalPages = response.totalPages ?? 0;
  return response.data;
});

const clearFilters = () => {
  filters.search = "";
  filters.status = "";
  filters.roomType = "";
  pagination.page = 1;
  // fetchRooms();
};

const confirmDelete = (room: Room) => {
  selectedRoom.value = room;
  deleteModalOpen.value = true;
};

const deleteRoom = async () => {
  if (!selectedRoom.value) return;

  deleting.value = true;
  try {
    await $fetch(`/api/rooms/${selectedRoom.value.id}`, { method: "DELETE" });
    deleteModalOpen.value = false;
    // await fetchRooms();
  } catch (error) {
    console.error("Failed to delete room:", error);
  } finally {
    deleting.value = false;
  }
};

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    available: "success",
    occupied: "warning",
    maintenance: "info",
    out_of_order: "error",
  };
  return colors[status] || "neutral";
};

const getRoomTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    single: "neutral",
    double: "info",
    suite: "primary",
    deluxe: "warning",
  };
  return colors[type] || "neutral";
};
</script>
