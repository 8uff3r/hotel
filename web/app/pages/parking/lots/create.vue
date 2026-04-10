<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/lots" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />
        Back to Parking Lots
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Create Parking Lot</h1>
    </div>

    <UCard>
      <form @submit.prevent="createParkingLot">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">Name *</label>
            <UInput v-model="form.name" placeholder="Main Parking" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Location</label>
            <UInput v-model="form.location" placeholder="Ground floor" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Total Spots</label>
            <UInput v-model="form.totalSpots" type="number" placeholder="50" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Hourly Rate ($)</label>
            <UInput v-model="form.hourlyRate" type="number" placeholder="5.00" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Daily Rate ($)</label>
            <UInput v-model="form.dailyRate" type="number" placeholder="25.00" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Status</label>
            <USelect v-model="form.status" :items="statusOptions" placeholder="Select status" />
          </div>

          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium">Description</label>
            <UTextarea v-model="form.description" placeholder="Additional details..." :rows="3" />
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton type="button" variant="outline" to="/parking/lots">Cancel</UButton>
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

const loading = ref(false);
const router = useRouter();

const createParkingLot = async () => {
  loading.value = true;
  try {
    await $fetch("/api/parking/lots", {
      method: "POST",
      body: {
        name: form.name,
        location: form.location || null,
        totalSpots: Number.parseInt(form.totalSpots) || 0,
        hourlyRate: Number.parseInt(form.hourlyRate) || 0,
        dailyRate: Number.parseInt(form.dailyRate) || 0,
        status: form.status,
        description: form.description || null,
      },
    });
    router.push("/parking/lots");
  } catch (error) {
    console.error("Failed to create parking lot:", error);
  } finally {
    loading.value = false;
  }
};
</script>
