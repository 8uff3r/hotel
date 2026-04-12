<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/lots" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />{{
          t("parking.back_to_parking_lots")
        }}</UButton
      >
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("parking.parking_lot_details") }}
      </h1>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <UIcon name="i-lucide-loader-2" class="h-8 w-8 animate-spin" />
    </div>

    <div v-else-if="parkingLot">
      <div class="grid-colsgrid-cols-3 gap-6-1 md: mb-6 grid">
        <UCard>
          <template #header>
            <span class="text-sm text-gray-500">{{ t("parking.total_spots") }}</span>
          </template>
          <div class="text-3xl font-bold">{{ parkingLot.totalSpots }}</div>
        </UCard>

        <UCard>
          <template #header>
            <span class="text-sm text-gray-500">Hourly Rate</span>
          </template>
          <div class="text-3xl font-bold">${{ parkingLot.hourlyRate }}</div>
        </UCard>

        <UCard>
          <template #header>
            <span class="text-sm text-gray-500">Daily Rate</span>
          </template>
          <div class="text-3xl font-bold">${{ parkingLot.dailyRate }}</div>
        </UCard>
      </div>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-semibold">Edit Parking Lot</span>
            <UBadge :color="getStatusColor(parkingLot.status) as any" variant="soft">
              {{ parkingLot.status }}
            </UBadge>
          </div>
        </template>

        <form @submit.prevent="updateParkingLot">
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div>
              <label class="mb-1 block text-sm font-medium">Name *</label>
              <UInput v-model="form.name" required />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">Location</label>
              <UInput v-model="form.location" />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">Total Spots</label>
              <UInput v-model="form.totalSpots" type="number" />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">Hourly Rate ($)</label>
              <UInput v-model="form.hourlyRate" />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">Daily Rate ($)</label>
              <UInput v-model="form.dailyRate" />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">Status</label>
              <USelect v-model="form.status" :items="statusOptions" />
            </div>

            <div class="md:col-span-2">
              <label class="mb-1 block text-sm font-medium">Description</label>
              <UTextarea v-model="form.description" :rows="3" />
            </div>
          </div>

          <div class="mt-6 flex justify-end gap-3">
            <UButton variant="outline" to="/parking/lots">Cancel</UButton>
            <UButton type="submit" color="primary" :loading="saving">Save Changes</UButton>
          </div>
        </form>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager"],
});
const { t } = useI18n();

const route = useRoute();
const parkingLotId = Number(route.params.id);

const parkingLot = ref<any>(null);
const loading = ref(true);
const saving = ref(false);

const form = reactive({
  name: "",
  location: "",
  totalSpots: "",
  hourlyRate: "",
  dailyRate: "",
  status: "active",
  description: "",
});

const statusOptions = [
  { value: "active", label: "Active" },
  { value: "full", label: "Full" },
  { value: "closed", label: "Closed" },
];

const fetchParkingLot = async () => {
  try {
    parkingLot.value = await $fetch(`/api/parking/lots/${parkingLotId}`);
    form.name = parkingLot.value.name;
    form.location = parkingLot.value.location || "";
    form.totalSpots = parkingLot.value.totalSpots?.toString() || "";
    form.hourlyRate = parkingLot.value.hourlyRate || "";
    form.dailyRate = parkingLot.value.dailyRate || "";
    form.status = parkingLot.value.status;
    form.description = parkingLot.value.description || "";
  } catch (error) {
    console.error("Failed to fetch parking lot:", error);
  } finally {
    loading.value = false;
  }
};

const updateParkingLot = async () => {
  saving.value = true;
  try {
    await $fetch(`/api/parking/lots/${parkingLotId}`, {
      method: "PUT",
      body: {
        name: form.name,
        location: form.location || null,
        totalSpots: parseInt(form.totalSpots) || 0,
        hourlyRate: form.hourlyRate,
        dailyRate: form.dailyRate,
        status: form.status,
        description: form.description || null,
      },
    });
    await fetchParkingLot();
  } catch (error) {
    console.error("Failed to update parking lot:", error);
  } finally {
    saving.value = false;
  }
};

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    active: "success",
    full: "warning",
    closed: "error",
  };
  return colors[status] || "neutral";
};

onMounted(fetchParkingLot);
</script>
