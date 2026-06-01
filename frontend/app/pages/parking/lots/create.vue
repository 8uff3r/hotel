<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/lots" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />{{
          t("parking.back_to_parking_lots")
        }}</UButton
      >
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("parking.create_parking_lot") }}
      </h1>
    </div>

    <UCard>
      <form @submit.prevent="createParkingLot">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.name") }} *</label>
            <UInput v-model="form.name" :placeholder="t('parking.main_parking')" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.location") }}</label>
            <UInput v-model="form.location" :placeholder="t('parking.ground_floor')" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.total_spots") }}</label>
            <UInput v-model="form.totalSpots" type="number" placeholder="50" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{
              t("parking.hourly_rate_dollar")
            }}</label>
            <UInput v-model="form.hourlyRate" type="number" placeholder="5.00" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{
              t("parking.daily_rate_dollar")
            }}</label>
            <UInput v-model="form.dailyRate" type="number" placeholder="25.00" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("common.status") }}</label>
            <USelect
              v-model="form.status"
              :items="statusOptions"
              :placeholder="t('parking.select_status')"
            />
          </div>

          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium">{{ t("common.description") }}</label>
            <UTextarea
              v-model="form.description"
              :placeholder="t('parking.additional_details')"
              :rows="3"
            />
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton type="button" variant="outline" to="/parking/lots">{{
            t("common.cancel")
          }}</UButton>
          <UButton type="submit" color="primary" :loading="loading">{{
            t("parking.create")
          }}</UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
const form = reactive({
  name: "",
  location: "",
  totalSpots: "",
  hourlyRate: "",
  dailyRate: "",
  status: "active",
  description: "",
});

const { t } = useI18n();
const statusOptions = [
  { value: "active", label: t("parking.status_active") },
  { value: "full", label: t("parking.status_full") },
  { value: "closed", label: t("parking.status_closed") },
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
