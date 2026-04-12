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

    <div v-if="loading" class="flex justify-center py-12">
      <ULoader size="lg" />
    </div>

    <UCard v-else-if="guest">
      <form @submit.prevent="handleSubmit">
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
          <UButton v-if="!editing" variant="outline" @click="editing = true">{{
            t("actions.editGuest")
          }}</UButton>
          <template v-else>
            <UButton variant="outline" :disabled="loading" @click="cancelEdit">{{
              t("actions.cancel")
            }}</UButton>
            <UButton type="submit" color="primary" :loading="loading">{{
              t("actions.saveChanges")
            }}</UButton>
          </template>
        </div>
      </form>
    </UCard>

    <UCard v-else>
      <div class="py-12 text-center">
        <p class="text-gray-500">{{ t("guests.notFound") }}</p>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

const { t } = useI18n();

interface Guest {
  id: number;
  roomNumber: string | null;
  referrer: string | null;
  roomType: string | null;
  registerCard: string | null;
  nationality: string | null;
  roomPrice: number | null;
  firstName: string;
  lastName: string;
  origin: string | null;
  fullBoard: boolean;
  destination: string | null;
  reservationCode: string | null;
  fatherName: string | null;
  purposeOfTravel: string | null;
  guide: boolean;
  idNumber: string | null;
  genderVisaValidity: string | null;
  breakfast: boolean;
  cash: boolean;
  agency: boolean;
  dateOfBirth: string | null;
  numberOfPeople: number | null;
  guestType: string | null;
  contractType: string | null;
  placeOfBirthStayDuration: string | null;
  nationalId: string | null;
  carLicensePlate: string | null;
  occupationVisaNumber: string | null;
  extraPerson: number | null;
  durationOfStay: string | null;
  entryDate: string | null;
  departureDate: string | null;
  address: string | null;
  postalCode: string | null;
  phone: string | null;
  userCheckIn: string | null;
  userCheckOut: string | null;
  notes: string | null;
}

const route = useRoute();
const guestId = Number(route.params.id);

const loading = ref(false);
const editing = ref(false);
const guest = ref<Guest | null>(null);

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

const fetchGuest = async () => {
  loading.value = true;
  try {
    const response = await $fetch<{ data: Guest }>(`/api/guests/${guestId}`);
    guest.value = response.data;
    resetForm();
  } catch (error) {
    console.error("Failed to fetch guest:", error);
    guest.value = null;
  } finally {
    loading.value = false;
  }
};

const resetForm = () => {
  if (!guest.value) return;
  form.roomNumber = guest.value.roomNumber || "";
  form.referrer = guest.value.referrer || "";
  form.roomType = guest.value.roomType || "";
  form.registerCard = guest.value.registerCard || "";
  form.nationality = guest.value.nationality || "";
  form.roomPrice = guest.value.roomPrice || 0;
  form.firstName = guest.value.firstName || "";
  form.lastName = guest.value.lastName || "";
  form.origin = guest.value.origin || "";
  form.fullBoard = guest.value.fullBoard || false;
  form.destination = guest.value.destination || "";
  form.reservationCode = guest.value.reservationCode || "";
  form.fatherName = guest.value.fatherName || "";
  form.purposeOfTravel = guest.value.purposeOfTravel || "";
  form.guide = guest.value.guide || false;
  form.idNumber = guest.value.idNumber || "";
  form.genderVisaValidity = guest.value.genderVisaValidity || "";
  form.breakfast = guest.value.breakfast || false;
  form.cash = guest.value.cash || false;
  form.agency = guest.value.agency || false;
  form.dateOfBirth = guest.value.dateOfBirth || "";
  form.numberOfPeople = guest.value.numberOfPeople || 1;
  form.guestType = guest.value.guestType || "";
  form.contractType = guest.value.contractType || "";
  form.placeOfBirthStayDuration = guest.value.placeOfBirthStayDuration || "";
  form.nationalId = guest.value.nationalId || "";
  form.carLicensePlate = guest.value.carLicensePlate || "";
  form.occupationVisaNumber = guest.value.occupationVisaNumber || "";
  form.extraPerson = guest.value.extraPerson || 0;
  form.durationOfStay = guest.value.durationOfStay || "";
  form.entryDate = guest.value.entryDate || "";
  form.departureDate = guest.value.departureDate || "";
  form.address = guest.value.address || "";
  form.postalCode = guest.value.postalCode || "";
  form.phone = guest.value.phone || "";
  form.userCheckIn = guest.value.userCheckIn || "";
  form.userCheckOut = guest.value.userCheckOut || "";
  form.notes = guest.value.notes || "";
};

const cancelEdit = () => {
  editing.value = false;
  resetForm();
};

const handleSubmit = async () => {
  loading.value = true;
  try {
    await $fetch(`/api/guests/${guestId}`, {
      method: "put",
      body: form,
    });

    editing.value = false;
    await fetchGuest();
  } catch (error) {
    console.error("Failed to update guest:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchGuest);
</script>
