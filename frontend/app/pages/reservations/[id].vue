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
          Reservation #{{ reservation?.id }}
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
      <UIcon name="i-lucide-loader" size="lg" />
    </div>

    <div v-else-if="reservation" class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <!-- Main Details -->
      <div class="space-y-6 lg:col-span-2">
        <!-- Guest & Room Info -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">{{ t("reservations.reservation_details") }}</h3>
          </template>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <!-- Guest Info -->
            <div>
              <p class="mb-1 text-sm text-gray-500">Guest</p>
              <p class="text-lg font-medium">
                {{ reservation.guest?.firstName }} {{ reservation.guest?.lastName }}
              </p>
              <p class="text-gray-600">{{ reservation.guest?.email }}</p>
              <p class="text-gray-600">{{ reservation.guest?.phone }}</p>
              <p v-if="reservation.guest?.address" class="text-gray-600">
                {{ reservation.guest?.address }},
                <!-- FIXME: -->
                <!-- {{ reservation.guest?.city }}, -->
                {{ reservation.guest?.nationality?.label }}
              </p>
            </div>

            <!-- Room Info -->
            <div v-for="room in reservation.rooms ?? []">
              <p class="mb-1 text-sm text-gray-500">Room</p>
              <p class="text-lg font-medium">{{ room?.roomNumber }}</p>
              <p class="text-gray-600 capitalize">{{ room?.roomType }}</p>
              <p class="text-gray-600">Floor {{ room?.floor }}</p>
              <p class="text-gray-600">Capacity: {{ room?.capacity }} guests</p>
            </div>
          </div>
        </UCard>

        <!-- Dates & Status -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">Stay Information</h3>
          </template>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <p class="text-sm text-gray-500">Check-in Date</p>
              <p class="font-medium">{{ formatDate(reservation.entryDate) }}</p>
              <p v-if="reservation.entryDate" class="text-sm text-green-600">
                <UIcon name="i-lucide-check-circle" class="mr-1 inline" />
                Checked in: {{ formatDateTime(reservation.entryDate) }}
              </p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Check-out Date</p>
              <p class="font-medium">{{ formatDate(reservation.entryDate) }}</p>
              <p v-if="reservation.departureDate" class="text-sm text-green-600">
                <UIcon name="i-lucide-check-circle" class="mr-1 inline" />
                Checked out: {{ formatDateTime(reservation.departureDate) }}
              </p>
            </div>
          </div>
        </UCard>
      </div>

      <!-- Sidebar Actions -->
      <div class="space-y-6">
        <!-- Actions -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">Actions</h3>
          </template>
          <div class="space-y-3">
            <UButton
              v-if="reservation.status?.slug === 'confirmed'"
              color="success"
              block
              :loading="processing"
              @click="handleCheckIn()"
            >
              <UIcon name="i-lucide-log-in" class="mr-2" />
              Check In Guest
            </UButton>
            <UButton
              v-if="reservation.status?.slug === 'checked_in'"
              color="warning"
              block
              :loading="processing"
              @click="handleCheckOut()"
            >
              <UIcon name="i-lucide-log-out" class="mr-2" />
              Check Out Guest
            </UButton>
            <UButton
              v-if="reservation.status?.slug === 'confirmed'"
              color="error"
              variant="outline"
              block
              :loading="processing"
              @click="handleCancel()"
            >
              <UIcon name="i-lucide-x" class="mr-2" />
              Cancel Reservation
            </UButton>
            <UButton
              v-if="reservation.status?.slug === 'no_show'"
              color="info"
              block
              :loading="processing"
              @click="handleCheckIn()"
            >
              <UIcon name="i-lucide-user-check" class="mr-2" />
              Mark as Arrived
            </UButton>
            <UButton
              v-if="reservation.guestId"
              color="success"
              variant="soft"
              block
              :to="`/guests/${reservation.guestId}/settle`"
            >
              <UIcon name="i-lucide-credit-card" class="mr-2" />
              Settle Account
            </UButton>
          </div>
        </UCard>

        <!-- Room Preview -->
        <UCard v-for="room in reservation.rooms">
          <template #header>
            <h3 class="text-lg font-semibold">Room Details</h3>
          </template>
          <div class="space-y-3">
            <div>
              <p class="text-sm text-gray-500">Room Number</p>
              <p class="font-medium">{{ room.roomNumber }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Type</p>
              <p class="font-medium capitalize">{{ room.roomType }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Base Price</p>
              <p class="font-medium">${{ room.basePrice }}/night</p>
            </div>
            <div v-if="room.amenities">
              <p class="text-sm text-gray-500">Amenities</p>
              <div class="mt-1 flex flex-wrap gap-1">
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
    <UCard v-else>
      <div class="py-12 text-center">
        <p class="text-gray-500">Reservation not found</p>
        <UButton to="/reservations" class="mt-4">Back to Reservations</UButton>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});
const { t } = useI18n();

const route = useRoute();
const reservationId = computed(() => "7");

const {
  data: reservation,
  isLoading: loading,
  refetch,
} = useQuery({
  key: () => ["reservations", "get", reservationId.value],
  query: async () => {
    const response = await getApiReservationIdDetailed({
      path: {
        id: reservationId.value,
      },
    });
    return response.data;
  },
});

const formatDate = (date: string | undefined) => {
  if (!date) return "";
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
};

const formatDateTime = (date: string) => {
  return new Date(date).toLocaleString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const { mutate: handleCheckIn, isLoading: processingCheckIn } = useMutation({
  mutation: async () => {
    await postApiReservationIdCheckOut({
      path: {
        id: reservationId.value,
      },
    });
  },
  onSettled: () => {
    refetch();
  },
});

const { mutate: handleCheckOut, isLoading: processingCheckOut } = useMutation({
  mutation: async () => {
    await postApiReservationIdCheckOut({
      path: {
        id: reservationId.value,
      },
    });
  },
  onSettled: () => {
    refetch();
  },
});

const processing = computed(
  () => processingCheckIn.value || processingCheckOut.value || processingCancel.value
);

const { mutate: handleCancel, isLoading: processingCancel } = useMutation({
  mutation: async () => {
    //FIXME:
    await putApiReservationId({} as any);
  },
  onSettled: () => {
    refetch();
  },
});
</script>
