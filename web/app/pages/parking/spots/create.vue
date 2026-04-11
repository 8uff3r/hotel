<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/spots" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />{{ t('parking.back_to_spots') }}</UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('parking.create_parking_spot') }}</h1>
    </div>

    <UCard>
      <form @submit.prevent="createSpot">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">Parking Lot *</label>
            <HSelect v-model="form.lotId" :items="lots" :placeholder="t('parking.select_lot')" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Spot Number *</label>
            <UInput v-model="form.spotNumber" :placeholder="t('parking.a_101')" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t('common.floor') }}</label>
            <UInput v-model="form.floor" placeholder="1" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t('parking.spot_type') }}</label>
            <USelect v-model="form.spotType" :items="spotTypeOptions" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t('common.status') }}</label>
            <USelect v-model="form.status" :items="statusOptions" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t('parking.is_covered') }}</label>
            <UCheckbox v-model="form.isCovered" :label="t('parking.covered_parking')" />
          </div>

          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium">{{ t('common.description') }}</label>
            <UTextarea v-model="form.description" :placeholder="t('parking.additional_details')" :rows="3" />
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton type="button" variant="outline" to="/parking/spots">{{ t('common.cancel') }}</UButton>
          <UButton type="submit" color="primary" :loading="loading">{{ t('parking.create') }}</UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager"],
});

const form = reactive({
  lotId: 0,
  spotNumber: "",
  floor: "",
  spotType: "standard",
  status: "available",
  isCovered: false,
  description: "",
});

const { t } = useI18n();
const loading = ref(false);
const router = useRouter();

const spotTypeOptions = [
  { value: "standard", label: "Standard" },
  { value: "handicap", label: "Handicap" },
  { value: "electric", label: "Electric" },
  { value: "compact", label: "Compact" },
  { value: "large", label: "Large" },
];

const statusOptions = [
  { value: "available", label: "Available" },
  { value: "occupied", label: "Occupied" },
  { value: "reserved", label: "Reserved" },
  { value: "maintenance", label: "Maintenance" },
];

const { data: lots } = useAsyncData(async () => {
  const res = await $fetch<{ data: ParkingLot[] }>("/api/parking/lots");
  return res.data;
});

const createSpot = async () => {
  loading.value = true;
  try {
    await $fetch("/api/parking/spots", {
      method: "POST",
      body: {
        lotId: form.lotId,
        spotNumber: form.spotNumber,
        floor: form.floor || null,
        spotType: form.spotType,
        status: form.status,
        isCovered: form.isCovered,
        description: form.description || null,
      },
    });
    router.push("/parking/spots");
  } catch (error) {
    console.error("Failed to create spot:", error);
  } finally {
    loading.value = false;
  }
};
</script>
