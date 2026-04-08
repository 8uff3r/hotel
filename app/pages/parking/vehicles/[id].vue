<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/vehicles" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />
        Back to Vehicles
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Vehicle Details</h1>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <UIcon name="i-lucide-loader-2" class="h-8 w-8 animate-spin" />
    </div>

    <UCard v-else-if="vehicle">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-semibold">Edit Vehicle</span>
          <UBadge :color="vehicle.isRegistered ? 'success' : 'warning'" variant="soft">
            {{ vehicle.isRegistered ? "Registered" : "Guest" }}
          </UBadge>
        </div>
      </template>

      <form @submit.prevent="updateVehicle">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">License Plate *</label>
            <UInput v-model="form.licensePlate" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Guest</label>
            <USelect
              v-model="form.guestId"
              :items="guestOptions"
              placeholder="Select guest (optional)"
              searchable
              clearable
            />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Vehicle Type</label>
            <USelect v-model="form.vehicleType" :items="typeOptions" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Make</label>
            <UInput v-model="form.make" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Model</label>
            <UInput v-model="form.model" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Color</label>
            <UInput v-model="form.color" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Is Registered?</label>
            <UCheckbox v-model="form.isRegistered" label="Registered vehicle" />
          </div>

          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium">Notes</label>
            <UTextarea v-model="form.notes" :rows="3" />
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/parking/vehicles">Cancel</UButton>
          <UButton type="submit" color="primary" :loading="saving">Save Changes</UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

const route = useRoute();
const vehicleId = Number(route.params.id);

const vehicle = ref<any>(null);
const loading = ref(true);
const saving = ref(false);
const guests = ref<any[]>([]);

const form = reactive({
  licensePlate: "",
  guestId: "",
  vehicleType: "car",
  make: "",
  model: "",
  color: "",
  isRegistered: true,
  notes: "",
});

const guestOptions = ref<{ value: string; label: string }[]>([]);

const typeOptions = [
  { value: "car", label: "Car" },
  { value: "motorcycle", label: "Motorcycle" },
  { value: "truck", label: "Truck" },
  { value: "van", label: "Van" },
  { value: "other", label: "Other" },
];

const fetchData = async () => {
  try {
    const [vehicleRes, guestsRes] = await Promise.all([
      $fetch(`/api/parking/vehicles/${vehicleId}`),
      $fetch("/api/guests"),
    ]);

    vehicle.value = vehicleRes;
    guests.value = guestsRes.data;

    guestOptions.value = guests.value.map((g: any) => ({
      value: g.id.toString(),
      label: `${g.firstName} ${g.lastName}`,
    }));

    form.licensePlate = vehicle.value.licensePlate || "";
    form.guestId = vehicle.value.guestId?.toString() || "";
    form.vehicleType = vehicle.value.vehicleType || "car";
    form.make = vehicle.value.make || "";
    form.model = vehicle.value.model || "";
    form.color = vehicle.value.color || "";
    form.isRegistered = !!vehicle.value.isRegistered;
    form.notes = vehicle.value.notes || "";
  } catch (error) {
    console.error("Failed to fetch vehicle:", error);
  } finally {
    loading.value = false;
  }
};

const updateVehicle = async () => {
  saving.value = true;
  try {
    await $fetch(`/api/parking/vehicles/${vehicleId}`, {
      method: "PUT",
      body: {
        licensePlate: form.licensePlate.toUpperCase(),
        guestId: form.guestId ? parseInt(form.guestId) : null,
        vehicleType: form.vehicleType,
        make: form.make || null,
        model: form.model || null,
        color: form.color || null,
        isRegistered: form.isRegistered,
        notes: form.notes || null,
      },
    });
    await fetchData();
  } catch (error) {
    console.error("Failed to update vehicle:", error);
  } finally {
    saving.value = false;
  }
};

onMounted(fetchData);
</script>
