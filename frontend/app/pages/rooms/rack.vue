<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("rooms.roomRack.title") }}
      </h1>
    </div>

    <!-- Status Legend -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <div class="flex items-center gap-2">
          <span class="text-sm font-medium text-gray-600 dark:text-gray-400"
            >{{ t("rooms.roomRack.legend") }}:</span
          >
        </div>
        <div class="flex flex-wrap gap-3">
          <div
            v-for="status in statuses"
            :key="status.id"
            class="flex items-center gap-2 rounded-md border border-primary/10 p-0.5"
          >
            <div
              class="h-4 w-4 rounded border border-gray-200 dark:border-gray-700"
              :style="{ backgroundColor: `#${status.colorHex || '94a3b8'}` }"
            />
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ status.label }}</span>
          </div>
        </div>

        <div class="ml-auto flex items-center gap-4">
          <!-- Display Mode Toggle -->
          <div class="flex items-center gap-2">
            <span class="text-sm text-gray-600 dark:text-gray-400"
              >{{ t("rooms.roomRack.displayMode") }}:</span
            >
            <USelect
              v-model="displayMode"
              :items="displayModeOptions"
              value-key="value"
              label-key="label"
              class="w-40"
            />
          </div>

          <!-- Floors Per Page (for paginated mode) -->
          <div v-if="displayMode === 'paginated'" class="flex items-center gap-2">
            <span class="text-sm text-gray-600 dark:text-gray-400"
              >{{ t("rooms.roomRack.floorsPerPage") }}:</span
            >
            <USelectMenu
              v-model="floorsPerPage"
              :items="floorsPerPageOptions"
              class="w-20"
              value-key="value"
              label-key="label"
            />
          </div>

          <!-- Sort Option -->
          <div class="flex items-center gap-2">
            <span class="text-sm text-gray-600 dark:text-gray-400"
              >{{ t("rooms.roomRack.sortBy") }}:</span
            >
            <USelect
              v-model="sortBy"
              :items="sortOptions"
              value-key="value"
              label-key="label"
              class="w-36"
            />
          </div>
        </div>
      </div>
    </UCard>

    <!-- Filters -->
    <UCard v-if="agencies?.length || roomTypes?.length" class="mb-6">
      <RoomRackFilters
        v-model:filters="filters"
        :room-types="roomTypes ?? []"
        :countries="countries ?? []"
        :agencies="agencies ?? []"
      />
    </UCard>

    <!-- Room Rack Grid -->
    <UCard>
      <!-- Summary Stats -->
      <template #header>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-6">
            <span class="text-lg font-semibold">{{
              t("rooms.roomRack.currentFloor", { floor: currentFloor })
            }}</span>
            <div class="flex items-center gap-4">
              <span
                v-for="status in currentFloorStatusCounts"
                :key="status.id"
                class="flex items-center gap-1.5"
              >
                <div
                  class="h-3 w-3 rounded-full"
                  :style="{ backgroundColor: `#${status.colorHex || '94a3b8'}` }"
                />
                <span class="text-sm text-gray-600 dark:text-gray-400">
                  {{ status.count }}
                </span>
              </span>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <UButton
              variant="ghost"
              size="sm"
              :disabled="displayMode === 'paginated' && currentFloor <= minFloor"
              @click="prevFloor"
            >
              <UIcon name="i-lucide-chevron-right" class="h-4 w-4 rtl:rotate-180" />
            </UButton>
            <span class="text-sm text-gray-500"> {{ currentFloor }} / {{ maxFloor }} </span>
            <UButton
              variant="ghost"
              size="sm"
              :disabled="displayMode === 'paginated' && currentFloor >= maxFloor"
              @click="nextFloor"
            >
              <UIcon name="i-lucide-chevron-left" class="h-4 w-4 rtl:rotate-180" />
            </UButton>
          </div>
        </div>
      </template>

      <!-- Paginated Mode: One Floor -->
      <div v-if="displayMode === 'paginated'" class="p-4">
        <div v-if="pending" class="flex h-64 items-center justify-center">
          <UIcon name="i-lucide-loader" class="h-8 w-8 animate-spin text-primary" />
        </div>
        <div
          v-else-if="(currentFloorRooms?.length ?? 0) === 0"
          class="flex h-64 items-center justify-center"
        >
          <p class="text-gray-500">{{ t("rooms.roomRack.noRoomsOnFloor") }}</p>
        </div>
        <div
          v-else
          class="grid grid-cols-4 gap-4 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 xl:grid-cols-12"
        >
          <RoomRackCard
            v-for="room in currentFloorRooms"
            :key="room.id"
            :room="room"
            @select="openDetail"
            @change-status="openStatusChange"
          />
        </div>
      </div>

      <!-- Compact Mode: All Floors with Infinite Scroll -->
      <div v-else-if="displayMode === 'compact'" class="p-4">
        <div v-if="pending" class="flex h-64 items-center justify-center">
          <UIcon name="i-lucide-loader" class="h-8 w-8 animate-spin text-primary" />
        </div>
        <div v-else-if="filteredRooms?.length === 0" class="flex h-64 items-center justify-center">
          <p class="text-gray-500">{{ t("rooms.roomRack.noRooms") }}</p>
        </div>
        <div v-else ref="scrollContainer" class="max-h-150 overflow-y-auto">
          <div v-for="floor in sortedFloors" :key="floor" class="mb-6">
            <h3 class="mb-3 text-lg font-semibold text-gray-700 dark:text-gray-300">
              {{ t("rooms.roomRack.floorNumber", { floor }) }}
            </h3>
            <div
              class="grid grid-cols-4 gap-2 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 xl:grid-cols-12"
            >
              <RoomRackCard
                v-for="room in getRoomsByFloor(floor)"
                :key="room.id"
                :room="room"
                @select="openDetail"
                @change-status="openStatusChange"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Grid Mode: Floors with Rooms -->
      <div v-else-if="displayMode === 'grid'" class="p-4">
        <div v-if="pending" class="flex h-64 items-center justify-center">
          <UIcon name="i-lucide-loader" class="h-8 w-8 animate-spin text-primary" />
        </div>
        <div v-else-if="filteredRooms?.length === 0" class="flex h-64 items-center justify-center">
          <p class="text-gray-500">{{ t("rooms.roomRack.noRooms") }}</p>
        </div>
        <div v-else ref="scrollContainer" class="max-h-150 overflow-y-auto">
          <div v-for="floor in sortedFloors" :key="floor" class="mb-8">
            <h3
              class="mb-3 flex items-center gap-2 text-lg font-semibold text-gray-700 dark:text-gray-300"
            >
              {{ t("rooms.roomRack.floorNumber", { floor }) }}
              <span class="text-sm font-normal text-gray-500">
                ({{ getRoomsByFloor(floor)?.length ?? 0 }} {{ t("rooms.roomRack.rooms") }})
              </span>
            </h3>
            <div
              class="grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6 xl:grid-cols-8"
            >
              <RoomRackCard
                v-for="room in getRoomsByFloor(floor)"
                :key="room.id"
                :room="room"
                @select="openDetail"
                @change-status="openStatusChange"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Floor Navigation for Paginated Mode -->
      <template v-if="displayMode === 'paginated'" #footer>
        <div class="flex items-center justify-center gap-2">
          <UButton
            v-for="floor in paginatedFloorList"
            :key="floor"
            :variant="floor === currentFloor ? 'solid' : 'ghost'"
            :color="floor === currentFloor ? 'primary' : 'neutral'"
            size="xs"
            @click="currentFloor = floor"
          >
            {{ floor }}
          </UButton>
        </div>
      </template>
    </UCard>

    <!-- Room Detail Modal -->
    <RoomDetailModal
      v-if="selectedRoom"
      :room="selectedRoom"
      :open="detailModalOpen"
      @update:open="detailModalOpen = $event"
      @change-status="openStatusChange"
    />

    <!-- Status Change Modal -->
    <RoomStatusChangeModal
      v-if="selectedRoom"
      :room="selectedRoom"
      :statuses="statuses ?? []"
      :open="statusModalOpen"
      @update:open="statusModalOpen = $event"
      @status-changed="onStatusChanged"
    />
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresPermission: PERMISSIONS.rooms.rooms.read,
});

import type { RackFiltersState } from "~/composables/useRoomRackData";
import RoomRackCard from "./components/RoomRackCard.vue";
import RoomRackFilters from "./components/RoomRackFilters.vue";
import RoomDetailModal from "./components/RoomDetailModal.vue";
import RoomStatusChangeModal from "./components/RoomStatusChangeModal.vue";
import type { GetApiRoomsRackResponse } from "~/utils/client";
import type { RoomRack } from "./types";

const { t } = useI18n();

const displayMode = ref<"paginated" | "compact" | "grid">("paginated");
const displayModeOptions: { value: "paginated" | "compact" | "grid"; label: string }[] = [
  { value: "paginated", label: t("rooms.roomRack.modes.paginated") },
  { value: "compact", label: t("rooms.roomRack.modes.compact") },
  { value: "grid", label: t("rooms.roomRack.modes.grid") },
];

const floorsPerPage = ref(1);
const floorsPerPageOptions = [
  { value: 1, label: "1" },
  { value: 2, label: "2" },
  { value: 3, label: "3" },
  { value: 4, label: "4" },
  { value: 5, label: "5" },
];

const sortBy = ref<"roomNumber" | "status" | "type">("roomNumber");
const sortOptions = [
  { value: "roomNumber", label: t("rooms.roomRack.sortOptions.roomNumber") },
  { value: "status", label: t("rooms.roomRack.sortOptions.status") },
  { value: "type", label: t("rooms.roomRack.sortOptions.type") },
];

const currentFloor = ref(1);

const { allRooms, pending, statuses, roomTypes, agencies, countries } = useRoomRackData();

const filters = ref<RackFiltersState>({
  roomTypeId: null,
  nationalityId: null,
  agencyId: null,
  entryDateFrom: "",
  entryDateTo: "",
  departureDateFrom: "",
  departureDateTo: "",
});

const filteredRooms = computed(() => {
  let rooms = allRooms.value || [];

  if (filters.value.roomTypeId) {
    rooms = rooms.filter((r) => r.roomType?.id === filters.value.roomTypeId);
  }

  if (filters.value.nationalityId) {
    rooms = rooms.filter(
      (r) => r.currentReservation?.guest?.nationality?.id === filters.value.nationalityId
    );
  }

  if (filters.value.agencyId) {
    const agency = agencies.value?.find((a) => a.id === filters.value.agencyId);
    if (agency) {
      rooms = rooms.filter((r) => r.currentReservation?.origin === agency.name);
    }
  }

  if (filters.value.entryDateFrom) {
    rooms = rooms.filter((r) => {
      const d = r.currentReservation?.entryDate;
      return d && d >= filters.value.entryDateFrom;
    });
  }

  if (filters.value.entryDateTo) {
    rooms = rooms.filter((r) => {
      const d = r.currentReservation?.entryDate;
      return d && d <= filters.value.entryDateTo;
    });
  }

  if (filters.value.departureDateFrom) {
    rooms = rooms.filter((r) => {
      const d = r.currentReservation?.departureDate;
      return d && d >= filters.value.departureDateFrom;
    });
  }

  if (filters.value.departureDateTo) {
    rooms = rooms.filter((r) => {
      const d = r.currentReservation?.departureDate;
      return d && d <= filters.value.departureDateTo;
    });
  }

  return rooms;
});

const floors = computed(() => {
  const floorSet = new Set<number>();
  filteredRooms.value?.forEach((room) => {
    if (room.floor) floorSet.add(room.floor);
  });
  return Array.from(floorSet).sort((a, b) => a - b);
});

const minFloor = computed(() => floors.value[0] || 1);
const maxFloor = computed(() => floors.value[floors.value.length - 1] || 1);

const sortedFloors = computed(() => {
  return [...floors.value];
});

const currentFloorRooms = computed(() => {
  return filteredRooms.value
    ?.filter((room) => room.floor === currentFloor.value)
    .sort((a, b) => {
      if (sortBy.value === "status") {
        return (a.status?.label || "").localeCompare(b.status?.label || "");
      }
      if (sortBy.value === "type") {
        return (a.roomType?.label || "").localeCompare(b.roomType?.label || "");
      }
      return (a.roomNumber || "").localeCompare(b.roomNumber || "");
    });
});

const currentFloorStatusCounts = computed(() => {
  const counts: Record<number, { id: number; name: string; colorHex?: string; count: number }> = {};
  currentFloorRooms.value?.forEach((room) => {
    if (room.status) {
      const statusId = room.status.id || 0;
      if (!counts[statusId]) {
        counts[statusId] = {
          id: statusId,
          name: room.status.label || "",
          colorHex: room.status.colorHex,
          count: 0,
        };
      }
      counts[statusId].count++;
    }
  });
  return Object.values(counts);
});

const paginatedFloorList = computed(() => floors.value);

const getRoomsByFloor = (floor: number) => {
  return filteredRooms.value
    ?.filter((room) => room.floor === floor)
    .sort((a, b) => {
      if (sortBy.value === "status") {
        return (a.status?.label || "").localeCompare(b.status?.label || "");
      }
      if (sortBy.value === "type") {
        return (a.roomType?.label || "").localeCompare(b.roomType?.label || "");
      }
      return (a.roomNumber || "").localeCompare(b.roomNumber || "");
    });
};

const nextFloor = () => {
  const step = floorsPerPage.value;
  const currentIndex = floors.value.indexOf(currentFloor.value);
  const nextIndex = Math.min(currentIndex + step, floors.value.length - 1);
  currentFloor.value = floors.value[nextIndex]!;
};

const prevFloor = () => {
  const step = floorsPerPage.value;
  const currentIndex = floors.value.indexOf(currentFloor.value);
  const prevIndex = Math.max(currentIndex - step, 0);
  currentFloor.value = floors.value[prevIndex]!;
};

watch(
  () => minFloor.value,
  (newFloor) => {
    if (currentFloor.value < newFloor) {
      currentFloor.value = newFloor;
    }
  }
);

watch(
  filteredRooms,
  (rooms) => {
    if ((rooms?.length ?? 0) > 0 && currentFloor.value < minFloor.value) {
      currentFloor.value = minFloor.value;
    }
  },
  { immediate: true }
);

const selectedRoom = ref<RoomRack>();
const detailModalOpen = ref(false);
const statusModalOpen = ref(false);

function openDetail(room: RoomRack) {
  selectedRoom.value = room;
  detailModalOpen.value = true;
}

function openStatusChange(room: RoomRack) {
  selectedRoom.value = room;
  detailModalOpen.value = false;
  statusModalOpen.value = true;
}

function onStatusChanged(roomId: number, statusId: number) {
  const room = allRooms.value?.find((r) => r.id === roomId);
  if (room) {
    room.statusId = statusId;
    const newStatus = statuses.value?.find((s) => s.id === statusId);
    if (newStatus) {
      room.status = newStatus;
    }
  }
}
</script>
