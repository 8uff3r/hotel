<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/guests" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />{{ t("actions.backToGuests") }}</UButton
      >
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("guests.guestNumber", { id: guest?.id }) }}
      </h1>
    </div>

    <div v-if="pending" class="flex justify-center py-12">
      <UIcon name="i-lucide-loader-2" class="h-8 w-8 animate-spin" />
    </div>

    <div v-else-if="guest" class="flex flex-col gap-6 lg:flex-row">
      <!-- Left: Main Content -->
      <div class="flex-1">
        <!-- Guest Info -->
        <UCard class="mb-6">
          <template #header>
            <span class="text-lg font-semibold">{{ t("guests.information") }}</span>
          </template>

          <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
            <div class="mb-2 flex items-center gap-2 md:col-span-3">
              <p class="text-sm text-gray-500">{{ t("common.status") }}:</p>
              <UBadge :style="{ backgroundColor: `#${guest.status?.colorHex}` }" variant="soft">
                {{ guest.status?.label }}
              </UBadge>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("forms.firstName") }}</p>
              <p class="font-medium">{{ guest.firstName }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("forms.lastName") }}</p>
              <p class="font-medium">{{ guest.lastName }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.fatherName") }}</p>
              <p class="font-medium">{{ guest.fatherName }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.nationality") }}</p>
              <p class="font-medium">{{ guest.nationality?.label }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.gender") }}</p>
              <p class="font-medium capitalize">{{ guest.gender }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.dateOfBirth") }}</p>
              <p class="font-medium">{{ formatDate(guest.dateOfBirth) }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.placeOfBirth") }}</p>
              <p class="font-medium">{{ guest.placeOfBirth }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.nationalId") }}</p>
              <p class="font-medium">{{ guest.nationalId }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.idPassportNumber") }}</p>
              <p class="font-medium">{{ guest.idNumber }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.passport") }}</p>
              <p class="font-medium">{{ guest.passport }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.occupation") }}</p>
              <p class="font-medium">{{ guest.occupation }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.phone") }}</p>
              <p class="font-medium">{{ guest.phone }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("forms.email") }}</p>
              <p class="font-medium">{{ guest.email }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("forms.landline") }}</p>
              <p class="font-medium">{{ guest.landline }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.postalCode") }}</p>
              <p class="font-medium">{{ guest.postalCode }}</p>
            </div>
            <div class="md:col-span-3">
              <p class="text-sm text-gray-500">{{ t("forms.address") }}</p>
              <p class="font-medium">{{ guest.address }}</p>
            </div>
          </div>
        </UCard>

        <!-- Stays Table -->
        <UCard class="mb-6">
          <template #header>
            <span class="text-lg font-semibold">{{ t("guests.stays") }}</span>
          </template>
          <div v-if="stays.length === 0" class="py-4 text-center text-gray-500">
            {{ t("guests.noStays") }}
          </div>
          <UTable v-else :data="stays" :columns="stayColumns" striped>
            <template #room-cell="{ row }">
              <p class="font-medium">{{ row.original.room?.roomNumber }}</p>
            </template>
            <template #dates-cell="{ row }">
              <div>
                <p class="text-sm">
                  <UIcon name="i-lucide-log-in" class="mr-1 inline h-3 w-3" />
                  {{ formatDate(row.original.entryDate) }}
                </p>
                <p class="text-sm">
                  <UIcon name="i-lucide-log-out" class="mr-1 inline h-3 w-3" />
                  {{ formatDate(row.original.departureDate) }}
                </p>
              </div>
            </template>
            <template #status-cell="{ row }">
              <UBadge
                :style="{ backgroundColor: `#${row.original.status?.colorHex}` }"
                variant="soft"
              >
                {{ row.original.status?.label }}
              </UBadge>
            </template>
            <template #actions-cell="{ row }">
              <div class="flex items-center gap-2">
                <UButton variant="ghost" size="sm" @click="openStayModal(row.original)">
                  <UIcon name="i-lucide-eye" class="h-4 w-4" />
                </UButton>
                <UButton variant="ghost" size="sm" :to="`/stays/${row.original.id}`">
                  <UIcon name="i-lucide-pencil" class="h-4 w-4" />
                </UButton>
              </div>
            </template>
          </UTable>
        </UCard>

        <!-- Reservations Table -->
        <UCard>
          <template #header>
            <span class="text-lg font-semibold">{{ t("guests.reservations") }}</span>
          </template>
          <div v-if="reservations.length === 0" class="py-4 text-center text-gray-500">
            {{ t("guests.noReservations") }}
          </div>
          <UTable v-else :data="reservations" :columns="reservationColumns" striped>
            <template #room-cell="{ row }">
              <p class="font-medium">
                {{ row.original.rooms?.map((r) => r.roomNumber).join(", ") }}
              </p>
            </template>
            <template #dates-cell="{ row }">
              <div>
                <p class="text-sm">
                  <UIcon name="i-lucide-log-in" class="mr-1 inline h-3 w-3" />
                  {{ formatDate(row.original.entryDate) }}
                </p>
                <p class="text-sm">
                  <UIcon name="i-lucide-log-out" class="mr-1 inline h-3 w-3" />
                  {{ formatDate(row.original.departureDate) }}
                </p>
              </div>
            </template>
            <template #status-cell="{ row }">
              <UBadge
                :style="{ backgroundColor: `#${row.original.status?.colorHex}` }"
                variant="soft"
              >
                {{ row.original.status?.label }}
              </UBadge>
            </template>
            <template #paymentStatus-cell="{ row }">
              <UBadge
                :style="{ backgroundColor: `#${row.original.payment?.status?.colorHex}` }"
                variant="soft"
              >
                {{ row.original.payment?.status?.label }}
              </UBadge>
            </template>
            <template #actions-cell="{ row }">
              <div class="flex items-center gap-2">
                <UButton variant="ghost" size="sm" @click="openReservationModal(row.original)">
                  <UIcon name="i-lucide-eye" class="h-4 w-4" />
                </UButton>
                <UButton variant="ghost" size="sm" :to="`/reservations/${row.original.id}`">
                  <UIcon name="i-lucide-pencil" class="h-4 w-4" />
                </UButton>
              </div>
            </template>
          </UTable>
        </UCard>
      </div>

      <!-- Right: Sidebar Info -->
      <div class="w-80 max-lg:hidden">
        <UCard>
          <template #header>
            <span class="font-semibold">{{ t("guests.quickInfo") }}</span>
          </template>
          <div class="space-y-4">
            <div class="flex items-center gap-3">
              <UIcon name="i-lucide-user" class="h-6 w-6 text-primary" />
              <div>
                <p class="text-sm text-gray-500">{{ t("guests.fullName") }}</p>
                <p class="font-medium">{{ guest.firstName }} {{ guest.lastName }}</p>
              </div>
            </div>
            <div v-if="guest.fatherName" class="flex items-center gap-3">
              <UIcon name="i-lucide-users" class="h-6 w-6 text-primary" />
              <div>
                <p class="text-sm text-gray-500">{{ t("guest.fatherName") }}</p>
                <p class="font-medium">{{ guest.fatherName }}</p>
              </div>
            </div>
            <div v-if="guest.nationality" class="flex items-center gap-3">
              <UIcon name="i-lucide-globe" class="h-6 w-6 text-primary" />
              <div>
                <p class="text-sm text-gray-500">{{ t("guest.nationality") }}</p>
                <p class="font-medium">{{ guest.nationality.label }}</p>
              </div>
            </div>
            <div v-if="guest.phone" class="flex items-center gap-3">
              <UIcon name="i-lucide-phone" class="h-6 w-6 text-primary" />
              <div>
                <p class="text-sm text-gray-500">{{ t("guest.phone") }}</p>
                <p class="font-medium">{{ guest.phone }}</p>
              </div>
            </div>
            <div v-if="guest.gender" class="flex items-center gap-3">
              <UIcon name="i-lucide-venus-and-mars" class="h-6 w-6 text-primary" />
              <div>
                <p class="text-sm text-gray-500">{{ t("guest.gender") }}</p>
                <p class="font-medium capitalize">{{ guest.gender }}</p>
              </div>
            </div>
          </div>
        </UCard>

        <UCard v-if="guest.companions?.length" class="mt-6">
          <template #header>
            <span class="font-semibold">{{ t("guests.companions") }}</span>
          </template>
          <div class="space-y-3">
            <div
              v-for="companion in guest.companions"
              :key="companion.id"
              class="flex items-center gap-3 rounded-lg border p-3"
            >
              <UIcon name="i-lucide-user" class="h-5 w-5 text-primary" />
              <div>
                <p class="font-medium">{{ companion.firstName }} {{ companion.lastName }}</p>
                <p v-if="companion.relation" class="text-sm text-gray-500">
                  {{ companion.relation.label }}
                </p>
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </div>

    <!-- Not Found -->
    <div v-else class="py-12 text-center">
      <UIcon name="i-lucide-alert-circle" class="mx-auto h-12 w-12 text-gray-400" />
      <p class="mt-4 text-lg text-gray-500">{{ t("guests.notFound") }}</p>
      <UButton to="/guests" class="mt-4"> {{ t("actions.backToGuests") }} </UButton>
    </div>

    <!-- Stay Modal -->
    <UModal v-model:open="showStayModal">
      <template #body>
        <div v-if="selectedStay" class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">
              {{ t("stays.detailTitle") }} {{ selectedStay.acceptanceId }}
            </h3>
            <UBadge
              :style="{ backgroundColor: `#${selectedStay.status?.colorHex}` }"
              variant="soft"
            >
              {{ selectedStay.status?.label }}
            </UBadge>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <p class="text-sm text-gray-500">{{ t("guests.columns.room") }}</p>
              <p class="font-medium">{{ selectedStay.room?.roomNumber }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guests.columns.dates") }}</p>
              <p class="font-medium">
                {{ formatDate(selectedStay.entryDate) }} -
                {{ formatDate(selectedStay.departureDate) }}
              </p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.numberOfPeople") }}</p>
              <p class="font-medium">{{ selectedStay.numberOfPeople }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.roomPrice") }}</p>
              <p class="font-medium">${{ selectedStay.roomPrice?.toFixed(2) }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.origin") }}</p>
              <p class="font-medium">{{ selectedStay.origin }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.destination") }}</p>
              <p class="font-medium">{{ selectedStay.destination }}</p>
            </div>
          </div>
        </div>
      </template>
      <template #footer>
        <UButton variant="outline" @click="showStayModal = false">{{ t("actions.close") }}</UButton>
        <UButton v-if="selectedStay" color="primary" :to="`/stays/${selectedStay.id}`">{{
          t("guests.edit")
        }}</UButton>
      </template>
    </UModal>

    <!-- Reservation Modal -->
    <UModal v-model:open="showReservationModal">
      <template #body>
        <div v-if="selectedReservation" class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">
              {{ t("reservations.reservation_details") }} #{{ selectedReservation.id }}
            </h3>
            <UBadge
              :style="{ backgroundColor: `#${selectedReservation.status?.colorHex}` }"
              variant="soft"
            >
              {{ selectedReservation.status?.label }}
            </UBadge>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <p class="text-sm text-gray-500">{{ t("guests.columns.room") }}</p>
              <p class="font-medium">
                {{ selectedReservation.rooms?.map((r) => r.roomNumber).join(", ") }}
              </p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guests.columns.dates") }}</p>
              <p class="font-medium">
                {{ formatDate(selectedReservation.entryDate) }} -
                {{ formatDate(selectedReservation.departureDate) }}
              </p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.numberOfPeople") }}</p>
              <p class="font-medium">{{ selectedReservation.numberOfPeople }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.roomPrice") }}</p>
              <p class="font-medium">${{ selectedReservation.roomPrice?.toFixed(2) }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.origin") }}</p>
              <p class="font-medium">{{ selectedReservation.origin }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">{{ t("guest.destination") }}</p>
              <p class="font-medium">{{ selectedReservation.destination }}</p>
            </div>
          </div>
        </div>
      </template>
      <template #footer>
        <UButton variant="outline" @click="showReservationModal = false">{{
          t("actions.close")
        }}</UButton>
        <UButton
          v-if="selectedReservation"
          color="primary"
          :to="`/reservations/${selectedReservation.id}`"
          >{{ t("guests.edit") }}</UButton
        >
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import type { Stay, Reservation } from "~/utils/client";
import { getApiGuestsIdWithStay } from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.guests.guests.read,
});

const { t } = useI18n();
const route = useRoute();
const guestId = route.params.id as string;

const { data: guestData, pending } = useAsyncData(async () => {
  const response = await getApiGuestsIdWithStay({
    path: { id: guestId },
  });
  return response.data;
});

const guest = computed(() => guestData.value?.guest);
const stays = computed(() => guestData.value?.stays ?? []);
const reservations = computed(() => guestData.value?.reservations ?? []);

const stayColumns: TableColumn<Stay>[] = [
  { accessorKey: "id", header: () => t("guests.columns.id") },
  { accessorKey: "room", header: () => t("guests.columns.room") },
  { accessorKey: "dates", header: () => t("guests.columns.dates") },
  { accessorKey: "status", header: () => t("guests.columns.status") },
  { accessorKey: "actions", header: () => t("guests.columns.actions") },
];

const reservationColumns: TableColumn<Reservation>[] = [
  { accessorKey: "id", header: () => t("guests.columns.id") },
  { accessorKey: "room", header: () => t("guests.columns.room") },
  { accessorKey: "dates", header: () => t("guests.columns.dates") },
  { accessorKey: "status", header: () => t("guests.columns.status") },
  { accessorKey: "paymentStatus", header: () => t("common.payment_status") },
  { accessorKey: "actions", header: () => t("guests.columns.actions") },
];

const showStayModal = ref(false);
const selectedStay = ref<Stay | null>(null);

const openStayModal = (stay: Stay) => {
  selectedStay.value = stay;
  showStayModal.value = true;
};

const showReservationModal = ref(false);
const selectedReservation = ref<Reservation | null>(null);

const openReservationModal = (reservation: Reservation) => {
  selectedReservation.value = reservation;
  showReservationModal.value = true;
};

const formatDate = (date: Date | string | undefined) => {
  if (!date) return "";
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};
</script>
