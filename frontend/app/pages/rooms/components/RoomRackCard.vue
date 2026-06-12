<template>
  <div
    class="group relative flex cursor-pointer flex-col rounded-lg border-2 border-gray-200 p-3 transition-all hover:scale-105 hover:shadow-lg dark:border-gray-700"
    :style="{
      backgroundColor: bgColor,
      borderColor: borderColor,
    }"
    @click="$emit('select', room)"
    @contextmenu.prevent="openContextMenu($event)"
  >
    <!-- Status bar -->
    <div
      class="absolute top-0 right-0 left-0 h-1.5 rounded-t-lg"
      :style="{ backgroundColor: barColor }"
    />

    <!-- Status icon top-right -->
    <div class="absolute top-2.5 right-1.5">
      <UIcon :name="statusIcon" class="h-4 w-4 opacity-60" :style="{ color: barColor }" />
    </div>

    <!-- Room number -->
    <span class="mt-2 text-center font-bold text-gray-900 dark:text-white">
      {{ room.roomNumber }}
    </span>

    <!-- Room type -->
    <span class="text-center text-xs text-gray-500 dark:text-gray-400">
      {{ room.roomType?.label }}
    </span>

    <!-- Context Menu -->
    <Teleport to="body">
      <div
        v-if="contextMenuVisible"
        class="fixed z-50 w-64 rounded-lg border border-gray-200 bg-white p-4 shadow-xl dark:border-gray-700 dark:bg-gray-800"
        :style="{ left: `${contextMenuX}px`, top: `${contextMenuY}px` }"
        @click.stop
        @contextmenu.prevent
      >
        <div
          class="mb-3 flex items-center justify-between border-b border-gray-200 pb-2 dark:border-gray-600"
        >
          <h3 class="font-semibold text-gray-900 dark:text-white">
            {{ room.roomNumber }}
          </h3>
          <span class="text-xs text-gray-500">{{ room.roomType?.label }}</span>
        </div>

        <!-- Room info -->
        <div class="mb-4 space-y-1 text-xs text-gray-600 dark:text-gray-400">
          <div class="flex items-center gap-2">
            <UIcon name="i-lucide-layout-grid" class="h-3 w-3" />
            <span>{{ t("rooms.columns.floor") }}: {{ room.floor }}</span>
          </div>
          <div class="flex items-center gap-2">
            <UIcon name="i-lucide-arrow-left-right" class="h-3 w-3" />
            <span>{{ room.roomType?.label }}</span>
          </div>
        </div>

        <!-- Change Status Button -->
        <UButton color="primary" size="lg" class="w-full" @click.stop="handleChangeStatus">
          <UIcon name="i-lucide-rotate-ccw" class="mr-2 h-4 w-4" />
          {{ t("rooms.roomRack.changeStatus") }}
        </UButton>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import type { RoomRack } from "../types";

const props = defineProps<{
  room: RoomRack;
}>();

const emit = defineEmits<{
  select: [room: RoomRack];
  "change-status": [room: RoomRack];
}>();

const { t } = useI18n();

const contextMenuVisible = ref(false);
const contextMenuX = ref(0);
const contextMenuY = ref(0);

const hex = computed(() => `#${props.room.status?.colorHex || "94a3b8"}`);
const bgColor = computed(() => `${hex.value}15`);
const borderColor = computed(() => hex.value);
const barColor = computed(() => hex.value);

const statusIcon = computed(() => {
  const slug = props.room.status?.slug;
  const icons: Record<string, string> = {
    available: "i-lucide-check-circle",
    occupied: "i-lucide-door-open",
    reserved: "i-lucide-calendar-check",
    under_repair: "i-lucide-wrench",
    cleaning: "i-lucide-sparkles",
  };
  return icons[slug as string] || "i-lucide-circle";
});

function openContextMenu(event: MouseEvent) {
  contextMenuX.value = event.clientX;
  contextMenuY.value = event.clientY;
  contextMenuVisible.value = true;

  const close = () => {
    contextMenuVisible.value = false;
    document.removeEventListener("click", close);
  };
  setTimeout(() => document.addEventListener("click", close), 0);
}

function handleChangeStatus() {
  contextMenuVisible.value = false;
  emit("change-status", props.room);
}
</script>
