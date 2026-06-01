<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/spots" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />{{ t("parking.back_to_spots") }}</UButton
      >
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("parking.parking_spot_details") }}
      </h1>
    </div>

    <div v-if="pending" class="flex justify-center py-12">
      <UIcon name="i-lucide-loader-2" class="h-8 w-8 animate-spin" />
    </div>

    <UCard v-else-if="spot">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-semibold">{{ t("parking.edit_parking_spot") }}</span>
          <UBadge :color="getStatusColor(spot.status) as any" variant="soft">
            {{ spot.status }}
          </UBadge>
        </div>
      </template>

      <form @submit.prevent="updateSpot">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.parking_lot") }} *</label>
            <HSelect v-model="form.lotId" :items="lotOptions" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.spot_number") }} *</label>
            <UInput v-model="form.spotNumber" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("common.floor") }}</label>
            <UInput v-model="form.floor" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.spot_type") }}</label>
            <HSelect v-model="form.spotType" :items="spotTypeOptions" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("common.status") }}</label>
            <HSelect v-model="form.status" :items="statusOptions" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.is_covered") }}</label>
            <UCheckbox v-model="form.isCovered" :label="t('parking.covered_parking')" />
          </div>

          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium">{{ t("common.description") }}</label>
            <UTextarea v-model="form.description" :rows="3" />
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/parking/spots">{{ t("common.cancel") }}</UButton>
          <UButton type="submit" color="primary" :loading="saving">{{
            t("actions.saveChanges")
          }}</UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n();

const route = useRoute();
const spotId = Number(route.params.id);

const saving = ref(false);

const form = reactive<ParkingSpot>({
  lotId: 0,
  spotNumber: "",
  floor: "",
  spotType: "standard",
  status: "available",
  isCovered: false,
  description: "",
});

const { data: spotTypeOptions } = useAsyncData("parking-spot-types", async () => {
  const res = await $fetch<{ data: ParkingSpotType[] }>("/api/parking/spots/types");
  return res.data;
});
const { data: statusOptions } = useAsyncData("parking-spot-types", async () => {
  const res = await $fetch<{ data: ParkingSpotStatus[] }>("/api/parking/spots/types");
  return res.data;
});

const { data: lotOptions } = useAsyncData("parking-lots", async () => {
  const res = await $fetch<{ data: ParkingLot[] }>("/api/parking/lots");
  return res.data;
});

const {
  data: spot,
  pending,
  refresh,
} = useAsyncData(async () => {
  const res = await $fetch<ParkingSpot>(`/api/parking/spots/${spotId}`);

  form.lotId = res.lotId;
  form.spotNumber = res.spotNumber || "";
  form.floor = res.floor || "";
  form.spotType = res.spotType || "standard";
  form.status = res.status || "available";
  form.isCovered = res.isCovered || false;
  form.description = res.description || "";
  return res;
});

const updateSpot = async () => {
  saving.value = true;
  try {
    await $fetch(`/api/parking/spots/${spotId}`, {
      method: "PUT",
      body: {
        lotId: parseInt(form.lotId?.toString() ?? ""),
        spotNumber: form.spotNumber,
        floor: form.floor || null,
        spotType: form.spotType,
        status: form.status,
        isCovered: form.isCovered,
        description: form.description || null,
      },
    });
    refresh();
  } catch (error) {
    console.error("Failed to update spot:", error);
  } finally {
    saving.value = false;
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
</script>
