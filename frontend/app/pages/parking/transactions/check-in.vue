<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/transactions" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />{{
          t("parking.back_to_transactions")
        }}</UButton
      >
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("parking.check_in_vehicle") }}
      </h1>
    </div>

    <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
      <UCard>
        <template #header>
          <span class="font-semibold">{{ t("parking.vehicle_check_in") }}</span>
        </template>

        <form @submit.prevent="checkIn">
          <div class="space-y-4">
            <div>
              <label class="mb-1 block text-sm font-medium"
                >{{ t("parking.license_plate") }} *</label
              >
              <UInput v-model="form.licensePlate" :placeholder="t('parking.abc_1234')" required />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">{{ t("parking.parking_lot") }} *</label>
              <USelect
                v-model="form.lotId"
                :items="lotOptions"
                :placeholder="t('parking.select_lot')"
                required
                @change="loadSpots"
              />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">{{
                t("parking.parking_spot_optional")
              }}</label>
              <USelect
                v-model="form.spotId"
                :items="spotOptions"
                :placeholder="t('parking.select_spot_auto_assign_if_empty')"
              />
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
              <label class="mb-1 block text-sm font-medium">{{ t("parking.reservation") }}</label>
              <USelect
                v-model="form.reservationId"
                :items="reservationOptions"
                :placeholder="t('parking.link_to_reservation_optional')"
                clearable
              />
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium">{{ t("parking.rate_type") }}</label>
              <USelect v-model="form.rateApplied" :items="rateOptions" />
            </div>
          </div>

          <div class="mt-6 flex justify-end gap-3">
            <UButton type="button" variant="outline" to="/parking/transactions">{{
              t("common.cancel")
            }}</UButton>
            <UButton type="submit" color="primary" :loading="loading">
              <UIcon name="i-lucide-car" class="mr-2" />
              {{ t("parking.check_in") }}
            </UButton>
          </div>
        </form>
      </UCard>

      <UCard>
        <template #header>
          <span class="font-semibold">{{ t("parking.quick_stats") }}</span>
        </template>
        <div class="space-y-4">
          <div class="flex items-center justify-between rounded-lg bg-gray-50 p-3 dark:bg-gray-800">
            <span class="text-gray-500">{{ t("parking.available_spots") }}</span>
            <span class="text-2xl font-bold text-green-600">{{ stats.available }}</span>
          </div>
          <div class="flex items-center justify-between rounded-lg bg-gray-50 p-3 dark:bg-gray-800">
            <span class="text-gray-500">{{ t("parking.occupied_spots") }}</span>
            <span class="text-2xl font-bold text-warning">{{ stats.occupied }}</span>
          </div>
          <div class="flex items-center justify-between rounded-lg bg-gray-50 p-3 dark:bg-gray-800">
            <span class="text-gray-500">{{ t("parking.active_vehicles") }}</span>
            <span class="text-2xl font-bold text-info">{{ stats.active }}</span>
          </div>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

const form = reactive({
  licensePlate: "",
  lotId: "",
  spotId: "",
  guestId: "",
  reservationId: "",
  rateApplied: "hourly",
});

const lotOptions = ref<{ value: string; label: string }[]>([]);
const spotOptions = ref<{ value: string; label: string }[]>([]);
const guestOptions = ref<{ value: string; label: string }[]>([]);
const reservationOptions = ref<{ value: string; label: string }[]>([]);
const { t } = useI18n();
const loading = ref(false);
const router = useRouter();

const stats = reactive({
  available: 0,
  occupied: 0,
  active: 0,
});

const rateOptions = [
  { value: "hourly", label: t("parking.hourly") },
  { value: "daily", label: t("parking.daily") },
];

const fetchInitialData = async () => {
  try {
    const [lotsRes, guestsRes, resRes] = await Promise.all([
      $fetch("/api/parking/lots"),
      $fetch("/api/guests"),
      $fetch("/api/reservations"),
    ]);

    lotOptions.value = (lotsRes.data as any[]).map((l) => ({
      value: l.id.toString(),
      label: l.name,
    }));

    guestOptions.value = (guestsRes.data as any[]).map((g: any) => ({
      value: g.id.toString(),
      label: `${g.firstName} ${g.lastName}`,
    }));

    reservationOptions.value = (resRes.data as any[]).map((r: any) => ({
      value: r.id.toString(),
      label: `Res #${r.id} - ${r.guestId}`,
    }));

    await fetchStats();
  } catch (error) {
    console.error("Failed to fetch initial data:", error);
  }
};

const loadSpots = async () => {
  if (!form.lotId) return;

  try {
    const res = await $fetch(`/api/parking/spots?lotId=${form.lotId}&status=available`);
    spotOptions.value = [
      { value: "", label: t("parking.auto_assign") },
      ...(res.data as any[]).map((s) => ({
        value: s.id.toString(),
        label: `${s.spotNumber} (${s.spotType})`,
      })),
    ];
  } catch (error) {
    console.error("Failed to fetch spots:", error);
  }
};

const fetchStats = async () => {
  try {
    const [spotsRes, txRes] = await Promise.all([
      $fetch("/api/parking/spots"),
      $fetch("/api/parking/transactions"),
    ]);

    const spots = spotsRes.data as any[];
    const txs = txRes.data as any[];

    stats.available = spots.filter((s: any) => s.status === "available").length;
    stats.occupied = spots.filter((s: any) => s.status === "occupied").length;
    stats.active = txs.filter((t: any) => t.status === "active").length;
  } catch (error) {
    console.error("Failed to fetch stats:", error);
  }
};

const checkIn = async () => {
  loading.value = true;
  try {
    await $fetch("/api/parking/transactions", {
      method: "POST",
      body: {
        licensePlate: form.licensePlate.toUpperCase(),
        lotId: parseInt(form.lotId),
        spotId: form.spotId ? parseInt(form.spotId) : null,
        guestId: form.guestId ? parseInt(form.guestId) : null,
        reservationId: form.reservationId ? parseInt(form.reservationId) : null,
        rateApplied: form.rateApplied,
      },
    });
    router.push("/parking/transactions");
  } catch (error) {
    console.error("Failed to check in:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchInitialData);
</script>
