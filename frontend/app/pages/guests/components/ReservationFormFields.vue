<script setup lang="ts">
import type { CreateRequest } from "../utils";
defineProps<{
  loading: boolean;
}>();

const { t } = useI18n();
const form = defineModel<Required<CreateRequest>>({ default: {} });

const { data: rooms } = useAsyncData(async () => {
  const res = await getApiRooms({});
  return res.data?.data?.map((v) => ({
    id: v.id,
    label: `${v.name ?? v.roomNumber}`,
  }));
});
</script>
<template>
  <UFormField
    :label="t('guest.reservationCode')"
    name="reservation.reservationCode"
    class="md:col-span-3"
  >
    <UInput v-model="form.reservation.reservationCode" :disabled="true" />
  </UFormField>

  <UFormField :label="t('guest.entryDate')" name="reservation.entryDate" required>
    <HDate v-model="form.reservation.entryDate" />
  </UFormField>

  <UFormField :label="t('guest.departureDate')" name="reservation.departureDate">
    <HDate v-model="form.reservation.departureDate" />
  </UFormField>

  <UFormField :label="t('guest.numberOfPeople')" name="reservation.numberOfPeople">
    <UInput
      type="number"
      min="1"
      v-model.number="form.reservation.numberOfPeople"
      :disabled="loading"
    />
  </UFormField>

  <UFormField :label="t('guest.durationOfStay')" name="reservation.durationOfStay">
    <UInput
      type="number"
      min="1"
      v-model.number="form.reservation.durationOfStay"
      :disabled="loading"
    />
  </UFormField>

  <UFormField :label="t('guest.origin')" name="reservation.origin">
    <UInput v-model="form.reservation.origin" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.destination')" name="reservation.destination">
    <UInput
      v-model="form.reservation.destination"
      :disabled="loading"
      :placeholder="t('guest.defaultDestination')"
    />
  </UFormField>

  <UFormField :label="t('guest.purposeOfTravel')" name="reservation.purposeOfTravel">
    <UInput
      v-model="form.reservation.purposeOfTravel"
      :disabled="loading"
      :placeholder="t('guest.defaultPurposeOfTravel')"
    />
  </UFormField>

  <UFormField :label="t('guest.roomPrice')" name="reservation.roomPrice">
    <UInput type="number" min="0" v-model.number="form.reservation.roomPrice" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.selectRoom')" name="reservation.rooms">
    <HSelect v-model="form.reservation.rooms" :items="rooms ?? []" multiple :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.breakfast')" name="reservation.breakfast">
    <UCheckbox v-model="form.reservation.breakfast" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.fullBoard')" name="reservation.fullBoard">
    <UCheckbox v-model="form.reservation.fullBoard" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.parking')" name="reservation.parking">
    <UCheckbox v-model="form.reservation.parking" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.notes')" name="reservation.notes" class="md:col-span-3">
    <UTextarea v-model="form.reservation.notes" :rows="3" class="w-full" :disabled="loading" />
  </UFormField>
</template>
