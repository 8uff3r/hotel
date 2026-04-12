<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/vehicles" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />{{
          t("parking.back_to_vehicles")
        }}</UButton
      >
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("parking.register_vehicle") }}
      </h1>
    </div>

    <UCard>
      <form @submit.prevent="createVehicle">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">License Plate *</label>
            <UInput v-model="form.licensePlate" :placeholder="t('parking.abc_1234')" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("common.guest") }}</label>
            <USelect
              v-model="form.guestId"
              :items="guestOptions"
              :placeholder="t('parking.select_guest_optional')"
              searchable
              clearable
            />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.vehicle_type") }}</label>
            <USelect v-model="form.vehicleType" :items="typeOptions" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.make") }}</label>
            <UInput v-model="form.make" :placeholder="t('parking.toyota')" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.model") }}</label>
            <UInput v-model="form.model" :placeholder="t('parking.camry')" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.color") }}</label>
            <UInput v-model="form.color" :placeholder="t('parking.silver')" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.is_registered") }}</label>
            <UCheckbox v-model="form.isRegistered" :label="t('parking.registered_vehicle')" />
          </div>

          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium">{{ t("common.notes") }}</label>
            <UTextarea
              v-model="form.notes"
              :placeholder="t('parking.additional_details')"
              :rows="3"
            />
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton type="button" variant="outline" to="/parking/vehicles">{{
            t("common.cancel")
          }}</UButton>
          <UButton type="submit" color="primary" :loading="loading">{{
            t("parking.register")
          }}</UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

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
const { t } = useI18n();
const loading = ref(false);
const router = useRouter();

const typeOptions = [
  { value: "car", label: "Car" },
  { value: "motorcycle", label: "Motorcycle" },
  { value: "truck", label: "Truck" },
  { value: "van", label: "Van" },
  { value: "other", label: "Other" },
];

const fetchGuests = async () => {
  try {
    const res = await $fetch("/api/guests");
    guestOptions.value = (res.data as any[]).map((g) => ({
      value: g.id.toString(),
      label: `${g.firstName} ${g.lastName}`,
    }));
  } catch (error) {
    console.error("Failed to fetch guests:", error);
  }
};

const createVehicle = async () => {
  loading.value = true;
  try {
    await $fetch("/api/parking/vehicles", {
      method: "POST",
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
    router.push("/parking/vehicles");
  } catch (error) {
    console.error("Failed to create vehicle:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchGuests);
</script>
