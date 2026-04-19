<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/guests" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        {{ t("actions.backToGuests") }}
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("guests.guestNumber", { id: guest?.id }) }}
      </h1>
    </div>

    <div v-if="pending" class="flex justify-center py-12">
      <ULoader size="lg" />
    </div>

    <UCard v-else-if="guest">
      <UForm @submit="handleSubmit" :state="form" :schema>
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <UFormField label="Room Number" name="roomNumber"
            ><UInput v-model="form.roomNumber" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Referrer" name="referrer"
            ><UInput v-model="form.referrer" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Room Type" name="roomType"
            ><UInput v-model="form.roomType" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Register Card" name="registerCard"
            ><UInput v-model="form.registerCard" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Nationality" name="nationality"
            ><UInput v-model="form.nationality" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Room Price" name="roomPrice"
            ><UInput
              v-model.number="form.roomPrice"
              type="number"
              min="0"
              :disabled="loading || !editing"
          /></UFormField>

          <UFormField :label="t('forms.firstName')" name="firstName" required
            ><UInput v-model="form.firstName" :disabled="loading || !editing"
          /></UFormField>
          <UFormField :label="t('forms.lastName')" name="lastName" required
            ><UInput v-model="form.lastName" :disabled="loading || !editing"
          /></UFormField>

          <UFormField label="Origin" name="origin"
            ><UInput v-model="form.origin" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Destination" name="destination"
            ><UInput v-model="form.destination" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Reservation Code" name="reservationCode"
            ><UInput v-model="form.reservationCode" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Father's Name" name="fatherName"
            ><UInput v-model="form.fatherName" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Purpose of Travel / Border Entry" name="purposeOfTravel"
            ><UInput v-model="form.purposeOfTravel" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="ID/Passport Number" name="idNumber"
            ><UInput v-model="form.idNumber" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Gender / Visa Validity" name="genderVisaValidity"
            ><UInput v-model="form.genderVisaValidity" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Date of Birth" name="dateOfBirth"
            ><UInput v-model="form.dateOfBirth" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Number of People" name="numberOfPeople"
            ><UInput
              v-model.number="form.numberOfPeople"
              type="number"
              min="1"
              :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Guest Type" name="guestType"
            ><UInput v-model="form.guestType" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Contract Type" name="contractType"
            ><UInput v-model="form.contractType" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Place of Birth / Stay Duration" name="placeOfBirthStayDuration"
            ><UInput v-model="form.placeOfBirthStayDuration" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="National ID" name="nationalId"
            ><UInput v-model="form.nationalId" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Car License Plate" name="carLicensePlate"
            ><UInput v-model="form.carLicensePlate" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Occupation / Visa Number" name="occupationVisaNumber"
            ><UInput v-model="form.occupationVisaNumber" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Extra Person" name="extraPerson"
            ><UInput
              v-model.number="form.extraPerson"
              type="number"
              min="0"
              :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Duration of Stay" name="durationOfStay"
            ><UInput v-model="form.durationOfStay" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Entry Date" name="entryDate"
            ><UInput v-model="form.entryDate" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Departure Date" name="departureDate"
            ><UInput v-model="form.departureDate" :disabled="loading || !editing"
          /></UFormField>
          <UFormField :label="t('forms.address')" name="address" class="md:col-span-2"
            ><UInput v-model="form.address" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Postal Code" name="postalCode"
            ><UInput v-model="form.postalCode" :disabled="loading || !editing"
          /></UFormField>
          <UFormField :label="t('forms.phone')" name="phone"
            ><UInput v-model="form.phone" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="User - Check-in" name="userCheckIn"
            ><UInput v-model="form.userCheckIn" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="User - Check-out" name="userCheckOut"
            ><UInput v-model="form.userCheckOut" :disabled="loading || !editing"
          /></UFormField>

          <UFormField label="Full Board" name="fullBoard"
            ><UCheckbox v-model="form.fullBoard" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Guide" name="guide"
            ><UCheckbox v-model="form.guide" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Breakfast" name="breakfast"
            ><UCheckbox v-model="form.breakfast" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Cash" name="cash"
            ><UCheckbox v-model="form.cash" :disabled="loading || !editing"
          /></UFormField>
          <UFormField label="Agency" name="agency"
            ><UCheckbox v-model="form.agency" :disabled="loading || !editing"
          /></UFormField>

          <UFormField :label="t('forms.notes')" name="notes" class="md:col-span-2"
            ><UTextarea v-model="form.notes" :rows="3" :disabled="loading || !editing"
          /></UFormField>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" :disabled="loading">{{ t("actions.editGuest") }}</UButton>
          <UButton type="submit" color="primary" :loading="loading">{{
            t("actions.saveChanges")
          }}</UButton>
        </div>
      </UForm>
    </UCard>

    <UCard v-else>
      <div class="py-12 text-center">
        <p class="text-gray-500">{{ t("guests.notFound") }}</p>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { FormSubmitEvent } from "@nuxt/ui";
import { z as zod } from "zod/v4";
import { zGuest } from "~/utils/client/zod.gen";

definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

const { t } = useI18n();

const route = useRoute();
const guestId = route.params.id as string;

const loading = ref(false);

const schema = zGuest;
type Schema = zod.output<typeof schema>;

const form = ref<Schema>({} as any);

const { data: guest, pending } = useAsyncData(async () => {
  const response = await getApiGuestsId({
    path: {
      id: guestId,
    },
  });

  form.value = response;
  return response;
});

const cancelEdit = () => {
  if (guest.value) {
    form.value = { ...guest.value };
  }
};

const handleSubmit = async (event: FormSubmitEvent<Schema>) => {
  loading.value = true;
  try {
    await putApiGuestsId({
      path: {
        id: guestId,
      },
      body: event.data,
    });

    navigateTo("/guests");
  } catch (error) {
    console.error("Failed to update guest:", error);
  } finally {
    loading.value = false;
  }
};
</script>
