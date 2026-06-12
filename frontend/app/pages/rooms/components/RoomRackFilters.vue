<template>
  <div class="flex flex-wrap items-center gap-4">
    <USelect
      v-model="localFilters.roomTypeId"
      :items="roomTypeOptions"
      :placeholder="t('rooms.roomRack.allTypes')"
      class="w-full sm:w-40"
      clearable
    />
    <USelect
      v-model="localFilters.statusId"
      :items="statusOptions"
      :placeholder="t('rooms.roomRack.allStatuses')"
      class="w-full sm:w-40"
      clearable
    />
    <UButton variant="outline" size="sm" @click="clearFilters">
      {{ t("common.clear") }}
    </UButton>
  </div>
</template>

<script setup lang="ts">
import type { RackFiltersState } from "~/composables/useRoomRackData";
import type { Room, RoomStatus } from "~/utils/client";

const props = defineProps<{
  roomTypes: NonNullable<Room["roomType"]>[];
  statuses: RoomStatus[];
}>();

const modelValue = defineModel<RackFiltersState>("filters", {
  default: (): RackFiltersState => ({
    roomTypeId: null,
    statusId: null,
    nationalityId: null,
    agencyId: null,
    entryDateFrom: "",
    entryDateTo: "",
    departureDateFrom: "",
    departureDateTo: "",
  }),
});

const { t } = useI18n();

const localFilters = computed({
  get: () => modelValue.value,
  set: (val) => {
    modelValue.value = val;
  },
});

const roomTypeOptions = computed(() =>
  (props.roomTypes || []).map((rt: any) => ({
    id: rt.id,
    label: rt.label,
  }))
);

const statusOptions = computed(() =>
  (props.statuses || []).map((s: any) => ({
    id: s.id,
    label: s.label,
  }))
);

function clearFilters() {
  modelValue.value = {
    roomTypeId: null,
    statusId: null,
    nationalityId: null,
    agencyId: null,
    entryDateFrom: "",
    entryDateTo: "",
    departureDateFrom: "",
    departureDateTo: "",
  };
}
</script>
