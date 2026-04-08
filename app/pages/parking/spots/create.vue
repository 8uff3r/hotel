<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/spots" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />
        Back to Spots
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Create Parking Spot</h1>
    </div>

    <UCard>
      <form @submit.prevent="createSpot">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">Parking Lot *</label>
            <USelect v-model="form.lotId" :items="lotOptions" placeholder="Select lot" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Spot Number *</label>
            <UInput v-model="form.spotNumber" placeholder="A-101" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Floor</label>
            <UInput v-model="form.floor" placeholder="1" />
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
            <UTextarea v-model="form.description" placeholder="Additional details..." :rows="3" />
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton type="button" variant="outline" to="/parking/spots">Cancel</UButton>
          <UButton type="submit" color="primary" :loading="loading">Create</UButton>
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
  lotId: "",
  spotNumber: "",
  floor: "",
  spotType: "standard",
  status: "available",
  isCovered: false,
  description: "",
});

const lotOptions = ref<{ value: string; label: string }[]>([]);
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

const fetchLots = async () => {
  try {
    const res = await $fetch("/api/parking/lots");
    lotOptions.value = (res.data as any[]).map((l) => ({
      value: l.id.toString(),
      label: l.name,
    }));
  } catch (error) {
    console.error("Failed to fetch lots:", error);
  }
};

const createSpot = async () => {
  loading.value = true;
  try {
    await $fetch("/api/parking/spots", {
      method: "POST",
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
    router.push("/parking/spots");
  } catch (error) {
    console.error("Failed to create spot:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchLots);
</script>
