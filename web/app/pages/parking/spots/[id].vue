<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/spots" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />
        Back to Spots
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Parking Spot Details</h1>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <UIcon name="i-lucide-loader-2" class="h-8 w-8 animate-spin" />
    </div>

    <UCard v-else-if="spot">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-semibold">Edit Parking Spot</span>
          <UBadge :color="getStatusColor(spot.status) as any" variant="soft">
            {{ spot.status }}
          </UBadge>
        </div>
      </template>

      <form @submit.prevent="updateSpot">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">Parking Lot *</label>
            <USelect v-model="form.lotId" :items="lotOptions" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Spot Number *</label>
            <UInput v-model="form.spotNumber" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Floor</label>
            <UInput v-model="form.floor" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Spot Type</label>
            <USelect v-model="form.spotType" :items="spotTypeOptions" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Status</label>
            <USelect v-model="form.status" :items="statusOptions" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Is Covered?</label>
            <UCheckbox v-model="form.isCovered" label="Covered parking" />
          </div>

          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium">Description</label>
            <UTextarea v-model="form.description" :rows="3" />
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/parking/spots">Cancel</UButton>
          <UButton type="submit" color="primary" :loading="saving">Save Changes</UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager"],
});

const route = useRoute();
const spotId = Number(route.params.id);

const spot = ref<any>(null);
const loading = ref(true);
const saving = ref(false);
const lots = ref<any[]>([]);

const form = reactive({
  lotId: "",
  spotNumber: "",
  floor: "",
  spotType: "standard",
  status: "available",
  isCovered: false,
  description: "",
});

const lotOptions = ref<{ value: string; label: string }[]>([]);

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

const fetchData = async () => {
  try {
    const [spotRes, lotsRes] = await Promise.all([
      $fetch(`/api/parking/spots/${spotId}`),
      $fetch("/api/parking/lots"),
    ]);

    spot.value = spotRes;
    lots.value = lotsRes.data;

    lotOptions.value = lots.value.map((l: any) => ({
      value: l.id.toString(),
      label: l.name,
    }));

    form.lotId = spot.value.lotId?.toString() || "";
    form.spotNumber = spot.value.spotNumber || "";
    form.floor = spot.value.floor || "";
    form.spotType = spot.value.spotType || "standard";
    form.status = spot.value.status || "available";
    form.isCovered = spot.value.isCovered || false;
    form.description = spot.value.description || "";
  } catch (error) {
    console.error("Failed to fetch spot:", error);
  } finally {
    loading.value = false;
  }
};

const updateSpot = async () => {
  saving.value = true;
  try {
    await $fetch(`/api/parking/spots/${spotId}`, {
      method: "PUT",
      body: {
        lotId: parseInt(form.lotId),
        spotNumber: form.spotNumber,
        floor: form.floor || null,
        spotType: form.spotType,
        status: form.status,
        isCovered: form.isCovered,
        description: form.description || null,
      },
    });
    await fetchData();
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

onMounted(fetchData);
</script>
