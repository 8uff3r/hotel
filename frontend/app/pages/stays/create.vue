<template>
  <div class="p-6 max-w-2xl mx-auto">
    <h1 class="text-2xl font-bold mb-6">{{ t("stays.createTitle") }}</h1>

    <UAlert v-if="warnings.length > 0" color="warning" class="mb-4">
      <ul class="list-disc pr-4">
        <li v-for="(w, i) in warnings" :key="i">{{ w }}</li>
      </ul>
    </UAlert>

    <UForm :state="state" @submit="onSubmit">
      <div class="grid grid-cols-2 gap-4">
        <UFormGroup :label="t('stays.guest')" name="guestId">
          <USelect v-model="state.guestId" :options="guestOptions" />
        </UFormGroup>
        <UFormGroup :label="t('stays.room')" name="roomId">
          <USelect v-model="state.roomId" :options="roomOptions" @change="onRoomChange" />
        </UFormGroup>
        <UFormGroup :label="t('stays.entryDate')" name="entryDate">
          <UInput v-model="state.entryDate" type="datetime-local" />
        </UFormGroup>
        <UFormGroup :label="t('stays.departureDate')" name="departureDate">
          <UInput v-model="state.departureDate" type="datetime-local" />
        </UFormGroup>
        <UFormGroup :label="t('stays.numberOfPeople')" name="numberOfPeople">
          <UInput v-model.number="state.numberOfPeople" type="number" />
        </UFormGroup>
        <UFormGroup :label="t('stays.roomPrice')" name="roomPrice">
          <UInput v-model.number="state.roomPrice" type="number" />
        </UFormGroup>
      </div>
      <div class="grid grid-cols-4 gap-4 mt-4">
        <UFormGroup :label="t('stays.breakfast')" name="breakfast">
          <UToggle v-model="state.breakfast" />
        </UFormGroup>
        <UFormGroup :label="t('stays.halfBoard')" name="halfBoard">
          <UToggle v-model="state.halfBoard" />
        </UFormGroup>
        <UFormGroup :label="t('stays.fullBoard')" name="fullBoard">
          <UToggle v-model="state.fullBoard" />
        </UFormGroup>
        <UFormGroup :label="t('stays.parking')" name="parking">
          <UToggle v-model="state.parking" />
        </UFormGroup>
      </div>
      <UFormGroup :label="t('stays.notes')" name="notes" class="mt-4">
        <UTextarea v-model="state.notes" />
      </UFormGroup>
      <div class="mt-6 flex gap-3">
        <UButton type="submit" :loading="submitting">{{ t("stays.submit") }}</UButton>
        <UButton variant="outline" to="/stays">{{ t("actions.cancel") }}</UButton>
      </div>
    </UForm>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { getApiStaysCheckAvailability, getApiStaysCheckCapacity, postApiStays } from "~/utils/client";
import type { Guest, Room } from "~/utils/client";
import { useAuthStore } from "~/stores/auth";

definePageMeta({
  requiresPermission: "guests:create",
});

const { t } = useI18n();

const state = reactive({
  hotelId: "",
  guestId: undefined as number | undefined,
  roomId: undefined as number | undefined,
  entryDate: "",
  departureDate: "",
  numberOfPeople: 1,
  roomPrice: 0,
  breakfast: false,
  halfBoard: false,
  fullBoard: false,
  parking: false,
  notes: "",
});

const { data: guests } = useFetch("/api/guests", {
  key: "stays-guests",
  transform: (res) => (res as { data?: Guest[] })?.data ?? [],
});
const { data: rooms } = useFetch("/api/rooms", {
  key: "stays-rooms",
  transform: (res) => (res as { data?: Room[] })?.data ?? [],
});

const guestOptions = computed(() =>
  (guests.value ?? []).map((g) => ({ label: `${g.firstName} ${g.lastName}`, value: g.id }))
);
const roomOptions = computed(() =>
  (rooms.value ?? []).map((r) => ({ label: r.roomNumber, value: r.id }))
);

const authStore = useAuthStore();
watchEffect(() => {
  if (authStore.currentHotelId) {
    state.hotelId = authStore.currentHotelId;
  }
});

const warnings = ref<string[]>([]);

async function onRoomChange() {
  warnings.value = [];
  if (!state.roomId || !state.entryDate || !state.departureDate) return;

  // Check availability
  try {
    const res = await getApiStaysCheckAvailability({
      query: {
        roomId: String(state.roomId),
        entryDate: state.entryDate,
        departureDate: state.departureDate,
      },
      requestValidator: undefined,
    } as Parameters<typeof getApiStaysCheckAvailability>[0] & { query: Record<string, string> });
    const av = res.data;
    if (!av?.available) {
      warnings.value.push(t("stays.roomNotAvailable"));
    }
  } catch {
    // ignore
  }

  // Check capacity
  try {
    const res = await getApiStaysCheckCapacity({
      query: {
        roomId: String(state.roomId),
        guests: String(state.numberOfPeople),
      },
      requestValidator: undefined,
    } as Parameters<typeof getApiStaysCheckCapacity>[0] & { query: Record<string, string> });
    const cap = res.data;
    if (!cap?.ok) {
      warnings.value.push(t("stays.capacityWarning"));
    }
  } catch {
    // ignore
  }
}

watch(() => [state.entryDate, state.departureDate], () => {
  if (state.roomId) onRoomChange();
});

const submitting = ref(false);
async function onSubmit() {
  submitting.value = true;
  try {
    const response = await postApiStays({ body: state, requestValidator: undefined });
    const res = response.data;
    if (res?.warnings && res.warnings.length > 0) {
      warnings.value = res.warnings;
    }
    navigateTo("/stays");
  } catch (e: unknown) {
    console.error(e);
    const err = e as { data?: { error?: string } };
    if (err.data?.error === "room_occupied") {
      warnings.value.push(t("stays.roomNotAvailable"));
    } else {
      warnings.value.push(err.data?.error || t("stays.createError"));
    }
  } finally {
    submitting.value = false;
  }
}
</script>
