<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/reservations" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />{{
          t("reservations.back_to_reservations")
        }}</UButton
      >
      <div class="flex items-center justify-between">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          {{ t("reservations.edit") }} #{{ reservation?.id }}
        </h1>
        <div class="flex items-center gap-2">
          <UBadge
            v-if="reservation"
            :style="{ backgroundColor: `#${reservation.status?.colorHex}` }"
            variant="solid"
            size="lg"
          >
            {{ reservation.status?.label }}
          </UBadge>
        </div>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <UIcon name="i-lucide-loader-2" class="h-8 w-8 animate-spin" />
    </div>

    <div v-else-if="reservation" class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <!-- Main Form -->
      <div class="lg:col-span-2">
        <UCard>
          <template #header>
            <span class="text-lg font-semibold">{{ t("reservations.stay_information") }}</span>
          </template>

          <form @submit.prevent="handleSubmit">
            <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
              <!-- Entry Date -->
              <UFormField :label="t('guest.entryDate')" name="entryDate" required>
                <HDate v-model="form.entryDate" :disabled="saving" />
              </UFormField>

              <!-- Departure Date -->
              <UFormField :label="t('guest.departureDate')" name="departureDate">
                <HDate v-model="form.departureDate" :disabled="saving" />
              </UFormField>

              <!-- Number of People -->
              <UFormField :label="t('guest.numberOfPeople')" name="numberOfPeople">
                <UInput
                  v-model.number="form.numberOfPeople"
                  type="number"
                  min="1"
                  :disabled="saving"
                />
              </UFormField>

              <!-- Duration of Stay -->
              <UFormField :label="t('guest.durationOfStay')" name="durationOfStay">
                <UInput
                  v-model.number="form.durationOfStay"
                  type="number"
                  min="0"
                  :disabled="saving"
                />
              </UFormField>

              <!-- Payment Deadline -->
              <UFormField label="مهلت پرداخت" name="paymentDeadline">
                <HDate v-model="form.paymentDeadline" :disabled="saving" />
              </UFormField>

              <!-- Origin -->
              <UFormField :label="t('reservations.origin')" name="origin">
                <UInput v-model="form.origin" :disabled="saving" />
              </UFormField>

              <!-- Destination -->
              <UFormField :label="t('reservations.destination')" name="destination">
                <UInput v-model="form.destination" :disabled="saving" />
              </UFormField>

              <!-- Purpose of Travel -->
              <UFormField :label="t('reservations.purpose_of_travel')" name="purposeOfTravel">
                <UInput v-model="form.purposeOfTravel" :disabled="saving" />
              </UFormField>

              <!-- Room Price -->
              <UFormField :label="t('guest.roomPrice')" name="roomPrice">
                <UInput
                  v-model.number="form.roomPrice"
                  type="number"
                  min="0"
                  step="0.01"
                  :disabled="saving"
                />
              </UFormField>

              <!-- Options -->
              <div class="flex flex-wrap gap-6 md:col-span-2">
                <UFormField :label="t('guest.breakfast')" name="breakfast">
                  <UCheckbox v-model="form.breakfast" :disabled="saving" />
                </UFormField>
                <UFormField :label="t('guest.parking')" name="parking">
                  <UCheckbox v-model="form.parking" :disabled="saving" />
                </UFormField>
                <UFormField :label="t('guest.fullBoard')" name="fullBoard">
                  <UCheckbox v-model="form.fullBoard" :disabled="saving" />
                </UFormField>
              </div>

              <!-- Notes -->
              <UFormField :label="t('common.notes')" name="notes" class="md:col-span-2">
                <UTextarea v-model="form.notes" :disabled="saving" :rows="3" />
              </UFormField>
            </div>

            <div class="mt-6 flex justify-end gap-3">
              <UButton variant="outline" to="/reservations" :disabled="saving">
                {{ t("actions.cancel") }}
              </UButton>
              <UButton type="submit" color="primary" :loading="saving">
                {{ t("actions.saveChanges") }}
              </UButton>
            </div>
          </form>
        </UCard>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <!-- Actions -->
        <UCard>
          <template #header>
            <span class="text-lg font-semibold">{{ t("reservations.actions") }}</span>
          </template>
          <div class="space-y-3">
            <UButton
              v-if="['awaiting_payment', 'verified'].includes(reservation.status?.slug ?? '')"
              color="primary"
              block
              :loading="processing"
              @click="handleAccept()"
            >
              <UIcon name="i-lucide-check" class="mr-2" />
              {{ t("reservations.accept") }}
            </UButton>
            <UButton
              v-if="reservation.status?.slug === 'accepted'"
              color="success"
              block
              :loading="processing"
              @click="handleCheckIn()"
            >
              <UIcon name="i-lucide-log-in" class="mr-2" />
              {{ t("reservations.check_in_guest") }}
            </UButton>
            <UButton
              v-if="reservation.status?.slug === 'checked_in'"
              color="warning"
              block
              :loading="processing"
              @click="handleCheckOut()"
            >
              <UIcon name="i-lucide-log-out" class="mr-2" />
              {{ t("reservations.check_out_guest") }}
            </UButton>
            <UButton
              v-if="['awaiting_payment', 'verified', 'accepted'].includes(reservation.status?.slug ?? '')"
              color="error"
              variant="outline"
              block
              :loading="processing"
              @click="handleCancel()"
            >
              <UIcon name="i-lucide-x" class="mr-2" />
              {{ t("reservations.cancel_reservation") }}
            </UButton>
            <UButton
              v-if="reservation.guest?.id"
              color="success"
              variant="soft"
              block
              :to="`/guests/${reservation.guest.id}/settle`"
            >
              <UIcon name="i-lucide-credit-card" class="mr-2" />
              {{ t("reservations.settle_account") }}
            </UButton>
          </div>
        </UCard>

        <!-- Guest Info -->
        <UCard>
          <template #header>
            <span class="text-lg font-semibold">{{ t("common.guest") }}</span>
          </template>
          <div v-if="reservation.guest" class="space-y-3">
            <div>
              <p class="text-sm text-gray-500">{{ t("guests.fullName") }}</p>
              <p class="font-medium">
                {{ reservation.guest.firstName }} {{ reservation.guest.lastName }}
              </p>
            </div>
            <div v-if="reservation.guest.email">
              <p class="text-sm text-gray-500">{{ t("forms.email") }}</p>
              <p class="font-medium">{{ reservation.guest.email }}</p>
            </div>
            <div v-if="reservation.guest.phone">
              <p class="text-sm text-gray-500">{{ t("guest.phone") }}</p>
              <p class="font-medium">{{ reservation.guest.phone }}</p>
            </div>
            <div v-if="reservation.guest.nationality">
              <p class="text-sm text-gray-500">{{ t("guest.nationality") }}</p>
              <p class="font-medium">{{ reservation.guest.nationality.label }}</p>
            </div>
            <UButton variant="outline" size="sm" block :to="`/guests/${reservation.guestId}`">
              <UIcon name="i-lucide-external-link" class="mr-2" />
              {{ t("guests.guestNumber", { id: reservation.guestId }) }}
            </UButton>
          </div>
        </UCard>

        <!-- Room Details -->
        <UCard v-for="room in reservation.rooms ?? []" :key="room.id">
          <template #header>
            <span class="text-lg font-semibold">{{ t("reservations.room_details") }}</span>
          </template>
          <div class="space-y-3">
            <div>
              <p class="text-sm text-gray-500">{{ t("rooms.roomNumber") }}</p>
              <p class="font-medium">{{ room.roomNumber }}</p>
            </div>
            <div v-if="room.roomType">
              <p class="text-sm text-gray-500">{{ t("rooms.roomType") }}</p>
              <p class="font-medium">{{ room.roomType.label }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("rooms.basePrice") }}</p>
              <p class="font-medium">${{ room.basePrice }}/night</p>
            </div>
            <div v-if="room.amenities?.length">
              <p class="mb-1 text-sm text-gray-500">{{ t("rooms.amenities") }}</p>
              <div class="flex flex-wrap gap-1">
                <UBadge
                  v-for="amenity in room.amenities"
                  :key="amenity.id"
                  variant="soft"
                  size="sm"
                >
                  {{ amenity.label }}
                </UBadge>
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </div>

    <!-- Not Found -->
    <div v-else class="py-12 text-center">
      <UIcon name="i-lucide-alert-circle" class="mx-auto h-12 w-12 text-gray-400" />
      <p class="mt-4 text-lg text-gray-500">{{ t("reservations.not_found") }}</p>
      <UButton to="/reservations" class="mt-4">
        {{ t("reservations.back_to_reservations") }}
      </UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Reservation } from "~/utils/client";
import { getApiReservationIdDetailed, postApiReservationIdAccept, postApiReservationIdCheckIn, postApiReservationIdCheckOut, putApiReservationId } from "~/utils/client";

const { t } = useI18n();
const route = useRoute();
const reservationId = route.params.id as string;
const saving = ref(false);

const form = ref<Reservation>({} as Reservation);

const {
  data: reservation,
  isLoading: loading,
  refetch,
} = useQuery({
  key: () => ["reservations", "get", reservationId],
  query: async () => {
    const response = await getApiReservationIdDetailed({
      path: { id: reservationId },
    });

    form.value = {
      entryDate: response.data?.entryDate ?? "",
      departureDate: response.data?.departureDate ?? "",
      numberOfPeople: response.data?.numberOfPeople ?? 1,
      durationOfStay: response.data?.durationOfStay ?? 0,
      origin: response.data?.origin ?? "",
      destination: response.data?.destination ?? "",
      purposeOfTravel: response.data?.purposeOfTravel ?? "",
      roomPrice: response.data?.roomPrice ?? 0,
      breakfast: response.data?.breakfast ?? false,
      parking: response.data?.parking ?? false,
      fullBoard: response.data?.fullBoard ?? false,
      notes: response.data?.notes ?? "",
    } as Reservation;

    return response.data;
  },
});

const handleSubmit = async () => {
  saving.value = true;
  try {
    await putApiReservationId({
      path: { id: reservationId },
      body: form.value,
    });

    await navigateTo("/reservations");
  } catch (error) {
    console.error("Failed to update reservation:", error);
  } finally {
    saving.value = false;
  }
};

const { mutate: handleAccept, isLoading: processingAccept } = useMutation({
  mutation: async () => {
    await postApiReservationIdAccept({
      path: { id: reservationId },
    });
  },
  onSettled: () => {
    refetch();
  },
});

const { mutate: handleCheckIn, isLoading: processingCheckIn } = useMutation({
  mutation: async () => {
    await postApiReservationIdCheckIn({
      path: { id: reservationId },
    });
  },
  onSettled: () => {
    refetch();
  },
});

const { mutate: handleCheckOut, isLoading: processingCheckOut } = useMutation({
  mutation: async () => {
    await postApiReservationIdCheckOut({
      path: { id: reservationId },
    });
  },
  onSettled: () => {
    refetch();
  },
});

const { mutate: handleCancel, isLoading: processingCancel } = useMutation({
  mutation: async () => {
    await putApiReservationId({ path: { id: reservationId }, body: { status: { slug: "cancelled" } } });
  },
  onSettled: () => {
    refetch();
  },
});

const processing = computed(
  () => processingAccept.value || processingCheckIn.value || processingCheckOut.value || processingCancel.value
);
</script>
