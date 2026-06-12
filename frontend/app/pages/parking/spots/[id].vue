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
          <UBadge :color="getStatusColor(spot.status?.slug) as any" variant="soft">
            {{ spot.status?.label }}
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
import { getApiParkingLots, getApiParkingSpotsId, putApiParkingSpotsId } from "~/utils/client";
import type { ParkingLot, ParkingSpot } from "~/utils/client";

const { t } = useI18n();

const route = useRoute();
const spotId = Number(route.params.id);

const saving = ref(false);

const form = reactive({
  lotId: 0,
  spotNumber: "",
  floor: "",
  spotType: "standard",
  status: "available",
  isCovered: false,
  description: "",
});

const { data: spotTypeOptions } = useAsyncData("parking-spot-types", async () => {
  const res = await $fetch<{ data: Array<{ label: string; value: string }> }>(
    "/api/parking/spots/types"
  );
  return res.data;
});
const { data: statusOptions } = useAsyncData("parking-spot-statuses", async () => {
  const res = await $fetch<{ data: Array<{ label: string; value: string }> }>(
    "/api/parking/spots/types"
  );
  return res.data;
});

const { data: lotOptions } = useAsyncData("parking-lots", async () => {
  const res = await getApiParkingLots();
  return res.data?.data;
});

const {
  data: spot,
  pending,
  refresh,
} = useAsyncData(async () => {
  const res = await getApiParkingSpotsId({ path: { id: String(spotId) } });

  form.lotId = res.data?.lotId ?? 0;
  form.spotNumber = res.data?.spotNumber || "";
  form.floor = res.data?.floor || "";
  form.spotType = (res.data?.spotType as { slug?: string } | undefined)?.slug ?? "standard";
  form.status = (res.data?.status as { slug?: string } | undefined)?.slug ?? "available";
  form.isCovered = res.data?.isCovered || false;
  form.description = res.data?.description || "";
  return res.data;
});

const updateSpot = async () => {
  saving.value = true;
  try {
    await putApiParkingSpotsId({
      requestValidator: undefined,
      path: { id: String(spotId) },
      body: {
        lotId: parseInt(form.lotId?.toString() ?? ""),
        spotNumber: form.spotNumber,
        floor: form.floor || undefined,
        spotType: form.spotType ? { slug: form.spotType } : undefined,
        status: form.status ? { slug: form.status } : undefined,
        isCovered: form.isCovered,
        description: form.description || undefined,
      },
    });
    refresh();
  } catch (error) {
    console.error("Failed to update spot:", error);
  } finally {
    saving.value = false;
  }
};

const getStatusColor = (status: string | undefined) => {
  const colors: Record<string, string> = {
    available: "success",
    occupied: "warning",
    reserved: "info",
    maintenance: "error",
  };
  return colors[status ?? ""] || "neutral";
};
</script>
