<template>
  <UModal v-model:open="open">
    <template #title>
      <div class="flex items-center gap-2">
        <span>{{ t("rooms.roomRack.roomDetails", { room: room?.roomNumber }) }}</span>
        <UBadge
          :style="{ backgroundColor: `#${room?.status?.colorHex || '94a3b8'}` }"
          variant="soft"
          size="sm"
        >
          {{ room?.status?.label }}
        </UBadge>
      </div>
    </template>

    <template #content>
      <div class="space-y-6 p-4">
        <!-- Guest Info -->
        <div v-if="room.currentReservation?.guest">
          <h4 class="mb-3 font-semibold text-gray-900 dark:text-white">
            <UIcon name="i-lucide-user" class="mr-1 inline h-4 w-4" />
            {{ t("rooms.roomRack.guestInfo") }}
          </h4>
          <div class="grid grid-cols-2 gap-4 rounded-lg bg-gray-50 p-3 dark:bg-gray-900">
            <div>
              <p class="text-xs text-gray-500">{{ t("guests.firstName") }}</p>
              <p class="text-sm font-medium text-gray-900 dark:text-white">
                {{ room.currentReservation.guest.firstName }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("guests.lastName") }}</p>
              <p class="text-sm font-medium text-gray-900 dark:text-white">
                {{ room.currentReservation.guest.lastName }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("guests.fatherName") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">
                {{ room.currentReservation.guest.fatherName || "—" }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("guests.nationality") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">{{ guestNationality }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("guests.phone") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">
                {{ room.currentReservation.guest.phone || "—" }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("guests.nationalId") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">
                {{ room.currentReservation.guest.nationalId || "—" }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("guests.idNumber") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">
                {{ room.currentReservation.guest.idNumber || "—" }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("guests.occupation") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">
                {{ room.currentReservation.guest.occupation || "—" }}
              </p>
            </div>
          </div>
        </div>

        <!-- Reservation Info -->
        <div v-if="room.currentReservation">
          <h4 class="mb-3 font-semibold text-gray-900 dark:text-white">
            <UIcon name="i-lucide-calendar" class="mr-1 inline h-4 w-4" />
            {{ t("rooms.roomRack.reservationInfo") }}
          </h4>
          <div class="grid grid-cols-2 gap-4 rounded-lg bg-gray-50 p-3 dark:bg-gray-900">
            <div>
              <p class="text-xs text-gray-500">{{ t("reservations.guest_count") }}</p>
              <p class="text-sm font-medium text-gray-900 dark:text-white">
                {{ room.currentReservation.numberOfPeople }}
              </p>
            </div>
            <div v-if="room.currentReservation.origin">
              <p class="text-xs text-gray-500">{{ t("reservations.origin") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">
                {{ room.currentReservation.origin }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("reservations.entry_date") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">
                {{ formatDate(room.currentReservation.entryDate) }}
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("reservations.departure_date") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">
                {{ formatDate(room.currentReservation.departureDate) }}
              </p>
            </div>
            <div v-if="room.currentReservation.durationOfStay">
              <p class="text-xs text-gray-500">{{ t("reservations.duration") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">
                {{ room.currentReservation.durationOfStay }} {{ t("common.nights") }}
              </p>
            </div>
          </div>
        </div>

        <!-- Room Info -->
        <div>
          <h4 class="mb-3 font-semibold text-gray-900 dark:text-white">
            <UIcon name="i-lucide-door-open" class="mr-1 inline h-4 w-4" />
            {{ t("rooms.roomRack.roomInfo") }}
          </h4>
          <div class="grid grid-cols-2 gap-4 rounded-lg bg-gray-50 p-3 dark:bg-gray-900">
            <div>
              <p class="text-xs text-gray-500">{{ t("rooms.columns.roomNumber") }}</p>
              <p class="text-sm font-medium text-gray-900 dark:text-white">{{ room.roomNumber }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("rooms.columns.type") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">{{ room.roomType?.label }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("rooms.columns.floor") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">{{ room.floor }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("rooms.capacity") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">{{ room.capacity }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-500">{{ t("rooms.basePrice") }}</p>
              <p class="text-sm text-gray-700 dark:text-gray-300">{{ formatPrice }}</p>
            </div>
          </div>
        </div>

        <!-- Amenities -->
        <div v-if="room.amenities?.length">
          <h4 class="mb-3 font-semibold text-gray-900 dark:text-white">
            <UIcon name="i-lucide-star" class="mr-1 inline h-4 w-4" />
            {{ t("rooms.amenities") }}
          </h4>
          <div class="flex flex-wrap gap-2">
            <UBadge
              v-for="amenity in room.amenities"
              :key="amenity.id"
              variant="soft"
              color="neutral"
              size="sm"
            >
              {{ amenity.label }}
            </UBadge>
          </div>
        </div>

        <!-- Description -->
        <div v-if="room.description">
          <h4 class="mb-2 font-semibold text-gray-900 dark:text-white">
            {{ t("rooms.description") }}
          </h4>
          <p class="text-sm text-gray-600 dark:text-gray-400">{{ room.description }}</p>
        </div>

        <!-- Action: Change Status -->
        <div class="border-t border-gray-200 pt-4 dark:border-gray-700">
          <UButton color="primary" size="xl" class="w-full" @click="$emit('change-status', room)">
            <UIcon name="i-lucide-rotate-ccw" class="mr-2 h-5 w-5" />
            {{ t("rooms.roomRack.changeStatus") }}
          </UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
const props = defineProps<{
  room: any;
}>();

const emit = defineEmits<{
  "change-status": [room: any];
}>();

const open = defineModel<boolean>("open", { default: false });

const { t } = useI18n();

const guestNationality = computed(
  () => props.room.currentReservation?.guest?.nationality?.label || "—"
);

const formatPrice = computed(() => {
  const price = props.room.basePrice;
  if (price == null) return "—";
  return `$${Number(price).toFixed(2)}`;
});

function formatDate(date: string | undefined) {
  if (!date) return "—";
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}
</script>
