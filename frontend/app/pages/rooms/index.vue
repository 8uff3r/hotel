<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("rooms.title") }}</h1>
      <div class="flex items-center gap-2">
        <AddFloorModal
          v-model="floorForm"
          v-model:open="addFloorModalOpen"
          :loading="addingFloor"
          @confirm="addFloor"
          @cancel="resetFloor"
        >
          <UButton variant="outline">
            <UIcon name="i-lucide-plus" />
            {{ t("rooms.addFloor") }}
          </UButton>
        </AddFloorModal>
        <UButton to="/rooms/create" color="primary">
          <UIcon name="i-lucide-plus" />
          {{ t("rooms.addRoom") }}
        </UButton>
      </div>
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
          <UBadge
            :style="{ backgroundColor: `#${row.original.roomType?.colorHex}` }"
            variant="soft"
          >
            {{ row.original.roomType?.label }}
          </UBadge>
        </template>

        <template #status-cell="{ row }">
          <UBadge :style="{ backgroundColor: `#${row.original.status?.colorHex}` }" variant="soft">
            {{ row.original.status?.label }}
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

    <!-- Add Floor Modal -->
    <UModal v-model="addFloorModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("rooms.addFloor") }}</h2>
      </template>
      <template #body>
        <div class="space-y-4">
          <UFormField :label="t('rooms.floor')" required>
            <UInput
              v-model="floorForm.number"
              type="number"
              :placeholder="t('rooms.floorPlaceholder')"
            />
          </UFormField>
          <UFormField :label="t('common.description')">
            <UInput
              v-model="floorForm.description"
              :placeholder="t('rooms.descriptionPlaceholder')"
            />
          </UFormField>
        </div>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="addFloorModalOpen = false">
            {{ t("actions.cancel") }}
          </UButton>
          <UButton color="primary" :loading="addingFloor" @click="addFloor">
            {{ t("actions.add") }}
          </UButton>
        </div>
      </template>
    </UModal>

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
import { deleteApiRoomsId, postApiRoomsFloors } from "~/utils/client";
import type { TableColumn } from "@nuxt/ui";
import type { PaginatedResponseModelsRoom } from "~/utils/client";
import AddFloorModal from "./components/AddFloorModal.vue";

type Room = NonNullable<PaginatedResponseModelsRoom["data"]>[0];
definePageMeta({
  requiresPermission: PERMISSIONS.rooms.rooms.read,
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

const addingFloor = ref(false);
const addFloorModalOpen = ref(false);
const floorForm = ref({
  number: "",
  description: "",
});

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

const { data: rooms, isPending: pending } = useQuery({
  key: () => ["rooms", "list", pagination],
  query: async () => {
    const response = await getApiRooms({
      query: pagination,
    });
    pagination.total = response.data?.total ?? 0;
    pagination.totalPages = response.data?.totalPages ?? 0;
    return response.data?.data;
  },
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
    await deleteApiRoomsId({ path: { id: String(selectedRoom.value.id) } });
    deleteModalOpen.value = false;
    // await fetchRooms();
  } catch (error) {
    console.error("Failed to delete room:", error);
  } finally {
    deleting.value = false;
  }
};

const resetFloor = () => (floorForm.value = { number: "", description: "" });
const addFloor = async () => {
  if (!floorForm.value.number) return;
  addingFloor.value = true;
  try {
    await postApiRoomsFloors({
      body: {
        number: parseInt(floorForm.value.number),
        description: floorForm.value.description || undefined,
      },
    });
    addFloorModalOpen.value = false;
    resetFloor();
  } catch (error) {
    console.error("Failed to add floor:", error);
  } finally {
    addingFloor.value = false;
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

const getRoomTypeColor = (type: string | undefined) => {
  const colors: Record<string, string> = {
    single: "neutral",
    double: "info",
    suite: "primary",
    deluxe: "warning",
  };
  return type && colors[type] ? colors[type] : "neutral";
};
</script>
