<template>
  <div class="mx-auto max-w-2xl p-6">
    <h1 class="mb-6 text-2xl font-bold">{{ t("stays.editTitle") }}</h1>
    <UForm :state="state" @submit="onSubmit">
      <div class="grid grid-cols-2 gap-4">
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
      <div class="mt-4 grid grid-cols-4 gap-4">
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
        <UButton type="submit" :loading="submitting">{{ t("actions.save") }}</UButton>
        <UButton variant="outline" to="/stays">{{ t("actions.cancel") }}</UButton>
      </div>
    </UForm>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { putApiStaysId } from "~/utils/client";
import type { Stay } from "~/utils/client";

const route = useRoute();
const stayId = route.params.id as string;
const { t } = useI18n();

definePageMeta({
  requiresPermission: "guests:update",
});

const state = reactive({
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

const { data: stay } = useFetch(`/api/stays/${stayId}`, {
  key: `stay-edit-${stayId}`,
  transform: (res) => (res as { data?: Stay })?.data,
});

watchEffect(() => {
  if (stay.value) {
    state.entryDate = stay.value.entryDate ?? "";
    state.departureDate = stay.value.departureDate ?? "";
    state.numberOfPeople = stay.value.numberOfPeople ?? 1;
    state.roomPrice = stay.value.roomPrice ?? 0;
    state.breakfast = stay.value.breakfast ?? false;
    state.halfBoard = stay.value.halfBoard ?? false;
    state.fullBoard = stay.value.fullBoard ?? false;
    state.parking = stay.value.parking ?? false;
    state.notes = stay.value.notes ?? "";
  }
});

const submitting = ref(false);
async function onSubmit() {
  submitting.value = true;
  try {
    await putApiStaysId({ path: { id: stayId }, body: state, requestValidator: undefined });
    navigateTo("/stays");
  } catch (e) {
    console.error(e);
  } finally {
    submitting.value = false;
  }
}
</script>
