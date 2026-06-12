<script setup lang="ts">
import type { CreateRequest } from "../utils";
defineProps<{
  loading: boolean;
}>();

const { t } = useI18n();
const form = defineModel<CreateRequest>({ default: {} });

const { data: rooms } = useAsyncData(async () => {
  const res = await getApiRooms({});
  return res.data?.data?.map((v) => ({
    id: v.id,
    label: `${v.name ?? v.roomNumber}`,
  }));
});
</script>
<template>
  <template v-if="form.stay">
    <UFormField
      :label="t('guest.reservationCode')"
      name="reservation.reservationCode"
      class="md:col-span-3"
    >
      <UInput v-model="form.stay.stayCode" :disabled="true" />
    </UFormField>

    <UFormField :label="t('guest.entryDate')" name="reservation.entryDate" required>
      <HDate v-model="form.stay.entryDate" />
    </UFormField>

    <UFormField :label="t('guest.departureDate')" name="reservation.departureDate">
      <HDate v-model="form.stay.departureDate" />
    </UFormField>

    <UFormField :label="t('guest.numberOfPeople')" name="reservation.numberOfPeople">
      <UInput type="number" min="1" v-model.number="form.stay.numberOfPeople" :disabled="loading" />
    </UFormField>

    <UFormField :label="t('guest.durationOfStay')" name="reservation.durationOfStay">
      <UInput type="number" min="1" v-model.number="form.stay.durationOfStay" :disabled="loading" />
    </UFormField>

    <UFormField :label="t('guest.origin')" name="reservation.origin">
      <UInput v-model="form.stay.origin" :disabled="loading" />
    </UFormField>

    <UFormField :label="t('guest.destination')" name="reservation.destination">
      <UInput
        v-model="form.stay.destination"
        :disabled="loading"
        :placeholder="t('guest.defaultDestination')"
      />
    </UFormField>

    <UFormField :label="t('guest.purposeOfTravel')" name="reservation.purposeOfTravel">
      <UInput
        v-model="form.stay.purposeOfTravel"
        :disabled="loading"
        :placeholder="t('guest.defaultPurposeOfTravel')"
      />
    </UFormField>

    <UFormField :label="t('guest.roomPrice')" name="reservation.roomPrice">
      <UInput type="number" min="0" v-model.number="form.stay.roomPrice" :disabled="loading" />
    </UFormField>

    <UFormField :label="t('guest.selectRoom')" name="reservation.rooms">
      <HSelect v-model="form.stay.rooms" :items="rooms ?? []" multiple :disabled="loading" />
    </UFormField>

    <UFormField :label="t('guest.breakfast')" name="reservation.breakfast">
      <UCheckbox v-model="form.stay.breakfast" :disabled="loading" />
    </UFormField>

    <UFormField :label="t('guest.fullBoard')" name="reservation.fullBoard">
      <UCheckbox v-model="form.stay.fullBoard" :disabled="loading" />
    </UFormField>

    <UFormField :label="t('guest.parking')" name="reservation.parking">
      <UCheckbox v-model="form.stay.parking" :disabled="loading" />
    </UFormField>

    <UFormField :label="t('guest.notes')" name="reservation.notes" class="md:col-span-3">
      <UTextarea v-model="form.stay.notes" :rows="3" class="w-full" :disabled="loading" />
    </UFormField>
  </template>
</template>
