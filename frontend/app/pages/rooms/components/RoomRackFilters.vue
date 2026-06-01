<template>
  <div class="flex flex-wrap items-center gap-4">
    <USelect
      v-model="localFilters.roomTypeId"
      :items="roomTypeOptions"
      :placeholder="t('rooms.roomRack.allTypes')"
      class="w-full sm:w-40"
      clearable
    />
    <!-- <USelect -->
    <!--   v-model="localFilters.nationalityId" -->
    <!--   :items="nationalityOptions" -->
    <!--   :placeholder="t('rooms.roomRack.allNationalities')" -->
    <!--   class="w-full sm:w-40" -->
    <!--   clearable -->
    <!-- /> -->
    <!-- <USelect -->
    <!--   v-model="localFilters.agencyId" -->
    <!--   :items="agencyOptions" -->
    <!--   :placeholder="t('rooms.roomRack.allAgencies')" -->
    <!--   class="w-full sm:w-44" -->
    <!--   clearable -->
    <!-- /> -->
    <div class="flex gap-2">
      <UInput
        v-model="localFilters.entryDateFrom"
        type="date"
        :placeholder="t('rooms.roomRack.entryFrom')"
        class="w-full sm:w-36"
      />
      <UInput
        v-model="localFilters.entryDateTo"
        type="date"
        :placeholder="t('rooms.roomRack.entryTo')"
        class="w-full sm:w-36"
      />
    </div>
    <div class="flex gap-2">
      <UInput
        v-model="localFilters.departureDateFrom"
        type="date"
        :placeholder="t('rooms.roomRack.departureFrom')"
        class="w-full sm:w-36"
      />
      <UInput
        v-model="localFilters.departureDateTo"
        type="date"
        :placeholder="t('rooms.roomRack.departureTo')"
        class="w-full sm:w-36"
      />
    </div>
    <UButton variant="outline" size="sm" @click="clearFilters">
      {{ t("common.clear") }}
    </UButton>
  </div>
</template>

<script setup lang="ts">
import type { RackFiltersState } from "~/composables/useRoomRackData";
import type { Country, Room, TravelAgency } from "~/utils/client";

const props = defineProps<{
  roomTypes: NonNullable<Room["roomType"]>[];
  countries: Country[];
  agencies: TravelAgency[];
}>();

const modelValue = defineModel<RackFiltersState>("filters", {
  default: (): RackFiltersState => ({
    roomTypeId: null,
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

const nationalityOptions = computed(() =>
  (props.countries || []).map((c: any) => ({
    id: c.id,
    label: c.label,
  }))
);

const agencyOptions = computed(() =>
  (props.agencies || []).map((a: any) => ({
    id: a.id,
    label: a.name,
  }))
);

function clearFilters() {
  modelValue.value = {
    roomTypeId: null,
    nationalityId: null,
    agencyId: null,
    entryDateFrom: "",
    entryDateTo: "",
    departureDateFrom: "",
    departureDateTo: "",
  };
}
</script>
