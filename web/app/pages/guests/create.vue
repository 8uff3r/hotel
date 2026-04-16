<template>
  <div class="flex flex-col gap-6 lg:flex-row">
    <!-- Left: Main Form -->
    <div class="flex-1">
      <div class="mb-6">
        <UButton variant="ghost" to="/guests" class="mb-4">
          <UIcon name="i-lucide-arrow-left" class="mr-2" />
          {{ t("actions.backToGuests") }}
        </UButton>
        <h1 class="text-3xl font-bold">{{ t("guests.addNew") }}</h1>
      </div>

      <UCard>
        <form @submit.prevent="handleSubmit" class="space-y-8">
          <!-- SECTION: Guest Information -->
          <UCollapsible default-open>
            <template #default="{ open }">
              <UButton
                :label="t('guest.personalInfo')"
                color="neutral"
                variant="outline"
                :trailing-icon="open ? 'i-lucide-chevron-up' : 'i-lucide-chevron-down'"
                size="lg"
                block
              />
            </template>
            <template #content>
              <div
                class="mt-2 grid grid-cols-1 gap-6 rounded-md p-4 ring ring-accented/40 ring-inset md:grid-cols-3"
              >
                <UFormField :label="t('forms.firstName')" required>
                  <UInput v-model="guest.firstName" />
                </UFormField>

                <UFormField :label="t('forms.lastName')" required>
                  <UInput v-model="guest.lastName" />
                </UFormField>

                <UFormField :label="t('guest.fatherName')">
                  <UInput v-model="guest.fatherName" />
                </UFormField>

                <UFormField :label="t('guest.nationality')">
                  <UInput v-model="guest.nationality" />
                </UFormField>

                <UFormField :label="t('guest.gender')">
                  <UInput v-model="guest.gender" />
                </UFormField>

                <UFormField :label="t('guest.dateOfBirth')">
                  <UInput type="date" v-model="guest.dateOfBirth" />
                </UFormField>

                <UFormField :label="t('guest.placeOfBirth')">
                  <UInput v-model="guest.placeOfBirth" />
                </UFormField>

                <UFormField :label="t('guest.nationalId')">
                  <UInput v-model="guest.nationalId" />
                </UFormField>

                <UFormField :label="t('guest.idPassportNumber')">
                  <UInput v-model="guest.idNumber" />
                </UFormField>

                <UFormField :label="t('guest.occupation')">
                  <UInput v-model="guest.occupation" />
                </UFormField>

                <UFormField :label="t('guest.phone')">
                  <UInput v-model="guest.phone" />
                </UFormField>

                <UFormField :label="t('guest.postalCode')">
                  <UInput v-model="guest.postalCode" />
                </UFormField>

                <UFormField :label="t('forms.address')" class="md:col-span-3">
                  <UInput class="w-full" v-model="guest.address" />
                </UFormField>
              </div>
            </template>
          </UCollapsible>

          <!-- SECTION: Reservation -->
          <UCollapsible>
            <template #default="{ open }">
              <UButton
                :label="t('guest.reservationDetails')"
                color="neutral"
                variant="outline"
                :trailing-icon="open ? 'i-lucide-chevron-up' : 'i-lucide-chevron-down'"
                size="lg"
                block
              />
            </template>
            <template #content>
              <div
                class="mt-2 grid grid-cols-1 gap-6 rounded-md p-4 ring ring-accented/40 ring-inset md:grid-cols-3"
              >
                <UFormField :label="t('guest.roomId')">
                  <UInput v-model.number="reservation?.roomId" type="number" min="1" />
                </UFormField>

                <UFormField :label="t('guest.reservationCode')">
                  <UInput v-model="reservation?.reservationCode" />
                </UFormField>

                <UFormField :label="t('guest.entryDate')" required>
                  <UInput type="date" v-model="reservation?.entryDate" />
                </UFormField>

                <UFormField :label="t('guest.departureDate')">
                  <UInput type="date" v-model="reservation?.departureDate" />
                </UFormField>

                <UFormField :label="t('guest.numberOfPeople')">
                  <UInput type="number" min="1" v-model.number="reservation?.numberOfPeople" />
                </UFormField>

                <UFormField :label="t('guest.durationOfStay')">
                  <UInput type="number" min="1" v-model.number="reservation?.durationOfStay" />
                </UFormField>

                <UFormField :label="t('guest.origin')">
                  <UInput v-model="reservation?.origin" />
                </UFormField>

                <UFormField :label="t('guest.destination')">
                  <UInput v-model="reservation?.destination" />
                </UFormField>

                <UFormField :label="t('guest.purposeOfTravel')">
                  <UInput v-model="reservation?.purposeOfTravel" />
                </UFormField>

                <UFormField :label="t('guest.roomPrice')">
                  <UInput type="number" min="0" v-model.number="reservation?.roomPrice" />
                </UFormField>

                <UFormField :label="t('guest.breakfast')">
                  <UCheckbox v-model="reservation?.breakfast" />
                </UFormField>

                <UFormField :label="t('guest.guide')">
                  <UCheckbox v-model="reservation?.guide" />
                </UFormField>

                <UFormField :label="t('guest.checkInUser')">
                  <UInput v-model="reservation?.userCheckIn" />
                </UFormField>

                <UFormField :label="t('guest.checkOutUser')">
                  <UInput v-model="reservation?.userCheckOut" />
                </UFormField>

                <UFormField :label="t('guest.notes')" class="md:col-span-3">
                  <UTextarea v-model="reservation?.notes" :rows="3" class="w-full" />
                </UFormField>
              </div>
            </template>
          </UCollapsible>

          <!-- SECTION: Payment -->
          <UCollapsible>
            <template #default="{ open }">
              <UButton
                :label="t('guest.payment')"
                color="neutral"
                variant="outline"
                :trailing-icon="open ? 'i-lucide-chevron-up' : 'i-lucide-chevron-down'"
                size="lg"
                block
              />
            </template>
            <template #content>
              <div
                class="mt-2 grid grid-cols-1 gap-6 rounded-md p-4 ring ring-accented/40 ring-inset md:grid-cols-2"
              >
                <UFormField :label="t('guest.cash')">
                  <UCheckbox v-model="payment?.isCash" />
                </UFormField>

                <UFormField :label="t('guest.agency')">
                  <UCheckbox v-model="payment?.agency" />
                </UFormField>

                <UFormField :label="t('guest.referrer')">
                  <UInput v-model="payment?.referrer" />
                </UFormField>

                <UFormField :label="t('guest.contractType')">
                  <UInput v-model="payment?.contractType" />
                </UFormField>
              </div>
            </template>
          </UCollapsible>

          <div class="flex justify-end gap-3 pt-4">
            <UButton variant="outline" to="/guests">{{ t("actions.cancel") }}</UButton>
            <UButton type="submit" color="primary" :loading="loading">{{
              t("guests.createGuest")
            }}</UButton>
          </div>
        </form>
      </UCard>
    </div>

    <!-- Right: Sticky Summary -->
    <div class="w-full lg:w-80">
      <UCard class="sticky top-4">
        <h2 class="mb-4 text-xl font-semibold">{{ t("guest.reservationSummary") }}</h2>

        <div class="space-y-2 text-sm">
          <div>
            <strong>{{ t("guest.summaryGuest") }}:</strong> {{ guest.firstName }}
            {{ guest.lastName }}
          </div>
          <div>
            <strong>{{ t("guest.summaryDates") }}:</strong> {{ reservation?.entryDate }} →
            {{ reservation?.departureDate }}
          </div>
          <div>
            <strong>{{ t("guest.roomId") }}:</strong> {{ reservation?.roomId }}
          </div>
          <div>
            <strong>{{ t("guest.summaryPeople") }}:</strong> {{ reservation?.numberOfPeople }}
          </div>
          <div>
            <strong>{{ t("guest.summaryPrice") }}:</strong> {{ reservation?.roomPrice }}
            {{ t("guest.perNight") }}
          </div>
          <div>
            <strong>{{ t("guest.payment") }}:</strong>
            <span v-if="payment?.isCash">{{ t("guest.cash") }}</span>
            <span v-else-if="payment?.agency">{{ t("guest.agency") }}</span>
            <span v-else>{{ t("guest.unspecified") }}</span>
          </div>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n();
import { z } from "zod/v4";
const loading = ref(false);

const guestSchema = z.object({
  firstName: z.string(),
  lastName: z.string(),
  fatherName: z.string(),
  nationalId: z.string(),
  idNumber: z.string(),
  nationality: z.string(),
  gender: z.string(),
  dateOfBirth: z.string(),
  placeOfBirth: z.string(),
  phone: z.string(),
  address: z.string(),
  postalCode: z.string(),
  occupation: z.string(),
});

type Guest = z.output<typeof guestSchema>;

const guest = reactive<Guest>({
  firstName: "",
  lastName: "",
  fatherName: "",
  nationalId: "",
  idNumber: "",
  nationality: "",
  gender: "",
  dateOfBirth: "",
  placeOfBirth: "",
  phone: "",
  address: "",
  postalCode: "",
  occupation: "",
});

const reservationSchema = z
  .object({
    roomId: z.number().optional(),
    reservationCode: z.string(),
    entryDate: z.string(),
    departureDate: z.string(),
    durationOfStay: z.number(),
    numberOfPeople: z.number(),
    origin: z.string(),
    destination: z.string(),
    purposeOfTravel: z.string(),
    breakfast: z.boolean(),
    guide: z.boolean,
    roomPrice: z.number(),
    userCheckIn: z.string(),
    userCheckOut: z.string(),
    notes: z.string(),
  })
  .optional();

type Reservation = z.output<typeof reservationSchema>;

const reservation = ref<Reservation>({
  roomId: undefined,
  reservationCode: "",
  entryDate: "",
  departureDate: "",
  durationOfStay: 1,
  numberOfPeople: 1,
  origin: "",
  destination: "",
  purposeOfTravel: "",
  breakfast: false,
  guide: false,
  roomPrice: 0,
  userCheckIn: "",
  userCheckOut: "",
  notes: "",
});

const paymentSchema = z
  .object({
    isCash: z.boolean(),
    agency: z.boolean(),
    referrer: z.string(),
    contractType: z.string(),
  })
  .optional();

type Payment = z.output<typeof paymentSchema>;

const payment = ref<Payment>({
  isCash: false,
  agency: false,
  referrer: "",
  contractType: "",
});

const handleSubmit = async () => {
  loading.value = true;
  try {
    const body = {
      ...guest,
      dateOfBirth: new Date(guest.dateOfBirth).toISOString(),
      reservation,
      payment,
    };

    await $fetch("/api/guests", {
      method: "POST",
      body,
    });

    navigateTo("/guests");
  } catch (err) {
    console.error(err);
  } finally {
    loading.value = false;
  }
};
</script>
<style scoped>
div[data-slot="root"]:has(input) {
  width: 100%;
}
</style>
