<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Rooms</h1>
      <UButton to="/rooms/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />
        Add Room
      </UButton>
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          placeholder="Search rooms..."
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.status"
          :items="statusOptions"
          placeholder="All Statuses"
          class="w-full sm:w-40"
          @change="fetchRooms"
        />
        <USelect
          v-model="filters.roomType"
          :items="roomTypeOptions"
          placeholder="All Types"
          class="w-full sm:w-40"
          @change="fetchRooms"
        />
        <UButton variant="outline" @click="clearFilters"> Clear </UButton>
      </div>
    </UCard>

    <!-- Rooms Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">Room List</span>
          <span class="text-sm text-gray-500">{{ pagination.total }} rooms</span>
        </div>
      </template>

      <UTable :data="rooms" :columns="columns" :loading="loading" striped>
        <template #roomNumber-cell="{ row }">
          <NuxtLink
            :to="`/rooms/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            {{ row.original.roomNumber }}
          </NuxtLink>
        </template>

        <template #roomType-cell="{ row }">
          <UBadge :color="getRoomTypeColor(row.original.roomType) as any" variant="soft">
            {{ row.original.roomType }}
          </UBadge>
        </template>

        <template #status-cell="{ row }">
          <UBadge :color="getStatusColor(row.original.status) as any" variant="soft">
            {{ row.original.status }}
          </UBadge>
        </template>

        <template #basePrice-cell="{ row }"> ${{ row.original.basePrice?.toFixed(2) }} </template>

        <template #floor-cell="{ row }">
          {{ row.original.floor ?? "-" }}
        </template>

        <template #capacity-cell="{ row }"> {{ row.original.capacity }} guests </template>

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
            Page {{ pagination.page }} of {{ pagination.totalPages }}
          </span>
          <UPagination
            v-model="page"
            :page-count="pagination.limit"
            :total="pagination.total"
            @change="fetchRooms"
          />
        </div>
      </template>
    </UCard>

    <!-- Delete Confirmation Modal -->
    <UModal v-model="deleteModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">Confirm Delete</h2>
      </template>
      <template #body>
        <p>Are you sure you want to delete room {{ selectedRoom?.roomNumber }}?</p>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="deleteModalOpen = false">Cancel</UButton>
          <UButton color="error" :loading="deleting" @click="deleteRoom">Delete</UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import type { Room } from "~~/server/db/schema";

definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

interface RoomRow {
  id: number;
  roomNumber: string;
  roomType: string;
  floor: number | null;
  capacity: number;
  basePrice: number;
  status: string;
}

const columns: TableColumn<RoomRow>[] = [
  { accessorKey: "roomNumber", header: "Room #" },
  { accessorKey: "roomType", header: "Type" },
  { accessorKey: "floor", header: "Floor" },
  { accessorKey: "capacity", header: "Capacity" },
  { accessorKey: "basePrice", header: "Price" },
  { accessorKey: "status", header: "Status" },
  { accessorKey: "actions", header: "Actions" },
];

const statusOptions = [
  { value: "all", label: "All Statuses" },
  { value: "available", label: "Available" },
  { value: "occupied", label: "Occupied" },
  { value: "maintenance", label: "Maintenance" },
  { value: "out_of_order", label: "Out of Order" },
];

const roomTypeOptions = [
  { value: "all", label: "All Types" },
  { value: "single", label: "Single" },
  { value: "double", label: "Double" },
  { value: "suite", label: "Suite" },
  { value: "deluxe", label: "Deluxe" },
];

const rooms = ref<RoomRow[]>([]);
const loading = ref(false);
const deleting = ref(false);
const deleteModalOpen = ref(false);
const selectedRoom = ref<RoomRow | null>(null);
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
    fetchRooms();
  }, 300);
};

const fetchRooms = async () => {
  loading.value = true;
  try {
    const params = new URLSearchParams();
    params.append("page", pagination.page.toString());
    params.append("limit", pagination.limit.toString());

    if (filters.search) params.append("search", filters.search);
    if (filters.status) params.append("status", filters.status);
    if (filters.roomType) params.append("roomType", filters.roomType);

    const response = await $fetch(`/api/rooms?${params.toString()}`);
    rooms.value = response.data;
    pagination.total = response.pagination.total ?? 0;
    pagination.totalPages = response.pagination.totalPages ?? 0;
  } catch (error) {
    console.error("Failed to fetch rooms:", error);
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  filters.search = "";
  filters.status = "";
  filters.roomType = "";
  pagination.page = 1;
  fetchRooms();
};

const confirmDelete = (room: RoomRow) => {
  selectedRoom.value = room;
  deleteModalOpen.value = true;
};

const deleteRoom = async () => {
  if (!selectedRoom.value) return;

  deleting.value = true;
  try {
    await $fetch(`/api/rooms/${selectedRoom.value.id}`, { method: "DELETE" });
    deleteModalOpen.value = false;
    await fetchRooms();
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

onMounted(fetchRooms);
</script>
