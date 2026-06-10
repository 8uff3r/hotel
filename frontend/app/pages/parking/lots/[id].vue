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
            <span class="text-sm text-gray-500">{{ t("parking.hourly_rate") }}</span>
          </template>
          <div class="text-3xl font-bold">${{ parkingLot.hourlyRate }}</div>
        </UCard>

        <UCard>
          <template #header>
            <span class="text-sm text-gray-500">{{ t("parking.daily_rate") }}</span>
          </template>
          <div class="text-3xl font-bold">${{ parkingLot.dailyRate }}</div>
        </UCard>
      </div>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-semibold">{{ t("parking.edit_parking_lot") }}</span>
            <UBadge :color="getStatusColor(parkingLot.status) as any" variant="soft">
              {{ parkingLot.status }}
            </UBadge>
          </div>
        </template>

        <form @submit.prevent="updateParkingLot">
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div>
              <label class="mb-1 block text-sm font-medium">{{ t("parking.name") }} *</label>
              <UInput v-model="form.name" required />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">{{ t("parking.location") }}</label>
              <UInput v-model="form.location" />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">{{ t("parking.total_spots") }}</label>
              <UInput v-model="form.totalSpots" type="number" />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">{{
                t("parking.hourly_rate_dollar")
              }}</label>
              <UInput v-model="form.hourlyRate" />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">{{
                t("parking.daily_rate_dollar")
              }}</label>
              <UInput v-model="form.dailyRate" />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">{{ t("common.status") }}</label>
              <USelect v-model="form.status" :items="statusOptions" />
            </div>

            <div class="md:col-span-2">
              <label class="mb-1 block text-sm font-medium">{{ t("common.description") }}</label>
              <UTextarea v-model="form.description" :rows="3" />
            </div>
          </div>

          <div class="mt-6 flex justify-end gap-3">
            <UButton variant="outline" to="/parking/lots">{{ t("common.cancel") }}</UButton>
            <UButton type="submit" color="primary" :loading="saving">{{
              t("actions.saveChanges")
            }}</UButton>
          </div>
        </form>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getApiParkingLotsId, putApiParkingLotsId } from "~/utils/client";

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
  { value: "active", label: t("parking.status_active") },
  { value: "full", label: t("parking.status_full") },
  { value: "closed", label: t("parking.status_closed") },
];

const fetchParkingLot = async () => {
  try {
    const res = await getApiParkingLotsId({ path: { id: String(parkingLotId) } });
    parkingLot.value = res.data;
    form.name = res.data?.name || "";
    form.location = res.data?.location || "";
    form.totalSpots = String(res.data?.totalSpots ?? "");
    form.hourlyRate = String(res.data?.hourlyRate ?? "");
    form.dailyRate = String(res.data?.dailyRate ?? "");
    form.status = (res.data?.status as { slug?: string } | undefined)?.slug ?? "active";
    form.description = res.data?.description || "";
  } catch (error) {
    console.error("Failed to fetch parking lot:", error);
  } finally {
    loading.value = false;
  }
};

const updateParkingLot = async () => {
  saving.value = true;
  try {
    await putApiParkingLotsId({
      requestValidator: undefined,
      path: { id: String(parkingLotId) },
      body: {
        name: form.name,
        location: form.location || undefined,
        totalSpots: parseInt(form.totalSpots) || 0,
        hourlyRate: parseFloat(form.hourlyRate) || 0,
        dailyRate: parseFloat(form.dailyRate) || 0,
        status: { slug: form.status },
        description: form.description || undefined,
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
