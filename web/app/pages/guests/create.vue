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
                  <USelect
                    v-model="guest.gender"
                    class="w-full"
                    :items="[
                      { value: 'male', label: 'مرد' },
                      { value: 'female', label: 'زن' },
                    ]"
                  />
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
                  <HSelect v-model="room.id" :items="rooms ?? []" />
                </UFormField>

                <UFormField :label="t('guest.reservationCode')">
                  <UInput v-model="reservation.reservationCode" />
                </UFormField>

                <UFormField :label="t('guest.entryDate')" required>
                  <HDate v-model="reservation.entryDate" />
                </UFormField>

                <UFormField :label="t('guest.departureDate')">
                  <UInput type="date" v-model="reservation.departureDate" />
                </UFormField>

                <UFormField :label="t('guest.numberOfPeople')">
                  <UInput type="number" min="1" v-model.number="reservation.numberOfPeople" />
                </UFormField>

                <UFormField :label="t('guest.durationOfStay')">
                  <UInput type="number" min="1" v-model.number="reservation.durationOfStay" />
                </UFormField>

                <UFormField :label="t('guest.origin')">
                  <UInput v-model="reservation.origin" />
                </UFormField>

                <UFormField :label="t('guest.destination')">
                  <UInput v-model="reservation.destination" />
                </UFormField>

                <UFormField :label="t('guest.purposeOfTravel')">
                  <UInput v-model="reservation.purposeOfTravel" />
                </UFormField>

                <UFormField :label="t('guest.roomPrice')">
                  <UInput type="number" min="0" v-model.number="reservation.roomPrice" />
                </UFormField>

                <UFormField :label="t('guest.breakfast')">
                  <UCheckbox v-model="reservation.breakfast" />
                </UFormField>

                <UFormField :label="t('guest.guide')">
                  <UCheckbox v-model="reservation.guide" />
                </UFormField>

                <UFormField :label="t('guest.checkInUser')">
                  <UInput v-model="reservation.userCheckIn" />
                </UFormField>

                <UFormField :label="t('guest.checkOutUser')">
                  <UInput v-model="reservation.userCheckOut" />
                </UFormField>

                <UFormField :label="t('guest.notes')" class="md:col-span-3">
                  <UTextarea v-model="reservation.notes" :rows="3" class="w-full" />
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
                  <UCheckbox v-model="payment.isCash" />
                </UFormField>

                <UFormField :label="t('guest.agency')">
                  <UCheckbox v-model="payment.agency" />
                </UFormField>

                <UFormField :label="t('guest.referrer')">
                  <UInput v-model="payment.referrer" />
                </UFormField>

                <UFormField :label="t('guest.contractType')">
                  <UInput v-model="payment.contractType" />
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
            <strong>{{ t("guest.roomId") }}:</strong> {{ room?.id }}
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
import { zGuest, zReservation, zRoom } from "~/utils/client/zod.gen";
const loading = ref(false);

type Guest = z.output<typeof zGuest>;

const guest = reactive<Guest>({} as any);

type Reservation = z.output<typeof zReservation>;

const reservation = ref<Reservation>({} as any);

type Room = z.output<typeof zRoom>;

const room = ref<Room>({} as any);

const paymentSchema = z
  .object({
    isCash: z.boolean(),
    agency: z.boolean(),
    referrer: z.string(),
    contractType: z.string(),
  })
  .or(z.object({}))
  .default({});

type Payment = z.output<typeof paymentSchema>;

const payment = ref<Payment>({});

const { data: rooms } = useAsyncData(async () => {
  const res = await getApiRooms({});
  return res.data?.map((v) => ({
    id: v.id,
    name: `${v.name ?? v.roomNumber}`,
  }));
});

const handleSubmit = async () => {
  loading.value = true;
  try {
    const body = {
      ...guest,
      dateOfBirth: new Date(guest.dateOfBirth ?? "").toISOString(),
      // reservation: [reservation.value],
      // payment: payment.value,
    };

    await postApiGuests({
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
