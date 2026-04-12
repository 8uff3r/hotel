<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/guests" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        {{ t("actions.backToGuests") }}
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("guests.addNew") }}</h1>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
          <UFormField label="Room Number" name="roomNumber">
            <UInput v-model="form.roomNumber" :disabled="loading" />
          </UFormField>
          <UFormField label="Referrer" name="referrer">
            <UInput v-model="form.referrer" :disabled="loading" />
          </UFormField>
          <UFormField label="Room Type" name="roomType">
            <UInput v-model="form.roomType" :disabled="loading" />
          </UFormField>
          <UFormField label="Register Card" name="registerCard">
            <UInput v-model="form.registerCard" :disabled="loading" />
          </UFormField>
          <UFormField label="Nationality" name="nationality">
            <UInput v-model="form.nationality" :disabled="loading" />
          </UFormField>
          <UFormField label="Room Price" name="roomPrice">
            <UInput v-model.number="form.roomPrice" type="number" min="0" :disabled="loading" />
          </UFormField>

          <UFormField :label="t('forms.firstName')" name="firstName" required>
            <UInput v-model="form.firstName" :disabled="loading" />
          </UFormField>
          <UFormField :label="t('forms.lastName')" name="lastName" required>
            <UInput v-model="form.lastName" :disabled="loading" />
          </UFormField>

          <UFormField label="Origin" name="origin">
            <UInput v-model="form.origin" :disabled="loading" />
          </UFormField>
          <UFormField label="Destination" name="destination">
            <UInput v-model="form.destination" :disabled="loading" />
          </UFormField>
          <UFormField label="Reservation Code" name="reservationCode">
            <UInput v-model="form.reservationCode" :disabled="loading" />
          </UFormField>
          <UFormField label="Father's Name" name="fatherName">
            <UInput v-model="form.fatherName" :disabled="loading" />
          </UFormField>
          <UFormField label="Purpose of Travel / Border Entry" name="purposeOfTravel">
            <UInput v-model="form.purposeOfTravel" :disabled="loading" />
          </UFormField>
          <UFormField label="ID/Passport Number" name="idNumber">
            <UInput v-model="form.idNumber" :disabled="loading" />
          </UFormField>
          <UFormField label="Gender / Visa Validity" name="genderVisaValidity">
            <UInput v-model="form.genderVisaValidity" :disabled="loading" />
          </UFormField>
          <UFormField label="Date of Birth" name="dateOfBirth">
            <UInput v-model="form.dateOfBirth" :disabled="loading" />
          </UFormField>
          <UFormField label="Number of People" name="numberOfPeople">
            <UInput
              v-model.number="form.numberOfPeople"
              type="number"
              min="1"
              :disabled="loading"
            />
          </UFormField>
          <UFormField label="Guest Type" name="guestType">
            <UInput v-model="form.guestType" :disabled="loading" />
          </UFormField>
          <UFormField label="Contract Type" name="contractType">
            <UInput v-model="form.contractType" :disabled="loading" />
          </UFormField>
          <UFormField label="Place of Birth / Stay Duration" name="placeOfBirthStayDuration">
            <UInput v-model="form.placeOfBirthStayDuration" :disabled="loading" />
          </UFormField>
          <UFormField label="National ID" name="nationalId">
            <UInput v-model="form.nationalId" :disabled="loading" />
          </UFormField>
          <UFormField label="Car License Plate" name="carLicensePlate">
            <UInput v-model="form.carLicensePlate" :disabled="loading" />
          </UFormField>
          <UFormField label="Occupation / Visa Number" name="occupationVisaNumber">
            <UInput v-model="form.occupationVisaNumber" :disabled="loading" />
          </UFormField>
          <UFormField label="Extra Person" name="extraPerson">
            <UInput v-model.number="form.extraPerson" type="number" min="0" :disabled="loading" />
          </UFormField>
          <UFormField label="Duration of Stay" name="durationOfStay">
            <UInput v-model="form.durationOfStay" :disabled="loading" />
          </UFormField>
          <UFormField label="Entry Date" name="entryDate">
            <UInput v-model="form.entryDate" :disabled="loading" />
          </UFormField>
          <UFormField label="Departure Date" name="departureDate">
            <UInput v-model="form.departureDate" :disabled="loading" />
          </UFormField>

          <UFormField :label="t('forms.address')" name="address" class="md:col-span-2">
            <UInput v-model="form.address" :disabled="loading" />
          </UFormField>
          <UFormField label="Postal Code" name="postalCode">
            <UInput v-model="form.postalCode" :disabled="loading" />
          </UFormField>
          <UFormField :label="t('forms.phone')" name="phone">
            <UInput v-model="form.phone" :disabled="loading" />
          </UFormField>
          <UFormField label="User - Check-in" name="userCheckIn">
            <UInput v-model="form.userCheckIn" :disabled="loading" />
          </UFormField>
          <UFormField label="User - Check-out" name="userCheckOut">
            <UInput v-model="form.userCheckOut" :disabled="loading" />
          </UFormField>

          <UFormField label="Full Board" name="fullBoard">
            <UCheckbox v-model="form.fullBoard" :disabled="loading" />
          </UFormField>
          <UFormField label="Guide" name="guide">
            <UCheckbox v-model="form.guide" :disabled="loading" />
          </UFormField>
          <UFormField label="Breakfast" name="breakfast">
            <UCheckbox v-model="form.breakfast" :disabled="loading" />
          </UFormField>
          <UFormField label="Cash" name="cash">
            <UCheckbox v-model="form.cash" :disabled="loading" />
          </UFormField>
          <UFormField label="Agency" name="agency">
            <UCheckbox v-model="form.agency" :disabled="loading" />
          </UFormField>

          <UFormField :label="t('forms.notes')" name="notes" class="md:col-span-2">
            <UTextarea v-model="form.notes" :rows="3" :disabled="loading" />
          </UFormField>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/guests" :disabled="loading">{{
            t("actions.cancel")
          }}</UButton>
          <UButton type="submit" color="primary" :loading="loading">{{
            t("guests.createGuest")
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

const { t } = useI18n();
const loading = ref(false);

const form = reactive({
  roomNumber: "",
  referrer: "",
  roomType: "",
  registerCard: "",
  nationality: "",
  roomPrice: 0,
  firstName: "",
  lastName: "",
  origin: "",
  fullBoard: false,
  destination: "",
  reservationCode: "",
  fatherName: "",
  purposeOfTravel: "",
  guide: false,
  idNumber: "",
  genderVisaValidity: "",
  breakfast: false,
  cash: false,
  agency: false,
  dateOfBirth: "",
  numberOfPeople: 1,
  guestType: "",
  contractType: "",
  placeOfBirthStayDuration: "",
  nationalId: "",
  carLicensePlate: "",
  occupationVisaNumber: "",
  extraPerson: 0,
  durationOfStay: "",
  entryDate: "",
  departureDate: "",
  address: "",
  postalCode: "",
  phone: "",
  userCheckIn: "",
  userCheckOut: "",
  notes: "",
});

const handleSubmit = async () => {
  loading.value = true;
  try {
    await $fetch("/api/guests", {
      method: "POST",
      body: form,
    });

    await navigateTo("/guests");
  } catch (error) {
    console.error("Failed to create guest:", error);
  } finally {
    loading.value = false;
  }
};
</script>
