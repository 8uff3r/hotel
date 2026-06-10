<template>
  <UModal v-model:open="open">
    <template #title>
      {{ t("rooms.roomRack.changeStatusFor", { room: room?.roomNumber }) }}
    </template>

    <template #content>
      <div class="space-y-4 p-4">
        <!-- Guest Info -->
        <div
          v-if="room.currentReservation?.guest"
          class="rounded-lg border border-primary/20 bg-primary/5 p-3"
        >
          <h4 class="mb-2 text-sm font-semibold text-gray-900 dark:text-white">
            <UIcon name="i-lucide-user" class="mr-1 inline h-4 w-4" />
            {{ t("common.guest") }}
          </h4>
          <div class="space-y-1 text-sm text-gray-700 dark:text-gray-300">
            <p class="font-medium">
              {{ room.currentReservation.guest.firstName }}
              {{ room.currentReservation.guest.lastName }}
            </p>
            <div class="flex items-center gap-1 text-xs text-gray-500">
              <UIcon name="i-lucide-globe" class="h-3 w-3" />
              <span>{{ room.currentReservation.guest.nationality?.label || "—" }}</span>
            </div>
            <div class="flex items-center gap-1 text-xs text-gray-500">
              <UIcon name="i-lucide-phone" class="h-3 w-3" />
              <span>{{ room.currentReservation.guest.phone || "—" }}</span>
            </div>
            <div class="flex items-center gap-1 text-xs text-gray-500">
              <UIcon name="i-lucide-calendar" class="h-3 w-3" />
              <span
                >{{ formatDate(room.currentReservation.entryDate) }} →
                {{ formatDate(room.currentReservation.departureDate) }}</span
              >
            </div>
          </div>
        </div>

        <p class="text-sm text-gray-600 dark:text-gray-400">
          {{ t("rooms.roomRack.selectNewStatus") }}
        </p>
        <div class="grid grid-cols-2 gap-3">
          <UButton
            v-for="status in statuses"
            :key="status.id"
            :style="{
              backgroundColor: status.colorHex ? `#${status.colorHex}20` : undefined,
              borderColor: status.colorHex ? `#${status.colorHex}` : undefined,
              color: status.colorHex ? `#${status.colorHex}` : undefined,
            }"
            :class="status.id === room.statusId ? 'ring-2 ring-offset-2' : ''"
            variant="outline"
            size="xl"
            class="flex flex-col items-center gap-1 py-4"
            @click="changeStatus(status)"
          >
            <UIcon :name="getStatusIcon(status.slug)" class="h-6 w-6" />
            <span class="text-sm">{{ status.label }}</span>
          </UButton>
        </div>

        <!-- Housekeeping review for checked-out (cleaning) rooms -->
        <div v-if="room.status?.slug === 'cleaning'" class="border-t pt-4 mt-4">
          <p class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
            بازبینی خدمات
          </p>
          <div class="flex gap-3">
            <UButton color="success" variant="soft" block @click="confirmCleaning">
              <UIcon name="i-lucide-check" class="mr-2" />
              تایید نظافت
            </UButton>
            <UButton color="warning" variant="soft" block @click="markRepairNeeded">
              <UIcon name="i-lucide-wrench" class="mr-2" />
              نیاز به تعمیر
            </UButton>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { putApiRoomsId } from "~/utils/client";
import type { RoomRack, RoomStatus } from "../types";

const props = defineProps<{
  room: RoomRack;
  statuses: RoomStatus[];
}>();

const open = defineModel<boolean>("open", { default: false });

const emit = defineEmits<{
  "status-changed": [roomId: number, statusId: number];
}>();

const { t } = useI18n();
const updating = ref(false);

function formatDate(date: string | undefined) {
  if (!date) return "—";
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

function getStatusIcon(slug: string | undefined) {
  if (!slug) return "";
  const icons: Record<string, string> = {
    available: "i-lucide-check-circle",
    occupied: "i-lucide-door-open",
    maintenance: "i-lucide-wrench",
    out_of_order: "i-lucide-x-circle",
    cleaning: "i-lucide-sparkles",
  };
  return icons[slug] || "i-lucide-circle";
}

async function changeStatus(status: RoomStatus) {
  if (updating.value || !props.room.id || !status.id) return;
  updating.value = true;
  try {
    await putApiRoomsId({ path: { id: String(props.room.id) }, body: { statusId: status.id } });
    emit("status-changed", props.room.id, status.id);
    open.value = false;
  } catch {
    // error handled by UI
  } finally {
    updating.value = false;
  }
}

async function confirmCleaning() {
  const availableStatus = props.statuses.find((s) => s.slug === "available");
  if (availableStatus) {
    await changeStatus(availableStatus);
  }
}

async function markRepairNeeded() {
  const repairStatus = props.statuses.find((s) => s.slug === "under_repair");
  if (repairStatus) {
    await changeStatus(repairStatus);
  }
}
</script>
