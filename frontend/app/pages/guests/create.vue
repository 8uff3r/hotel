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
        <UForm @submit="handleSubmit" :state="form" :schema class="flex flex-col gap-2">
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
                <UFormField :label="t('forms.firstName')" name="guest.firstName" required>
                  <UInput v-model="form.guest.firstName" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('forms.lastName')" name="guest.lastName" required>
                  <UInput v-model="form.guest.lastName" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.fatherName')" name="guest.fatherName">
                  <UInput v-model="form.guest.fatherName" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.nationality')" name="guest.nationality">
                  <UInput v-model="form.guest.nationality" :disabled="loading" :placeholder="t('guest.defaultNationality')" />
                </UFormField>

                <UFormField :label="t('guest.gender')" name="guest.gender">
                  <USelect
                    v-model="form.guest.gender"
                    class="w-full"
                    :items="[
                      { value: 'male', label: t('guest.male') },
                      { value: 'female', label: t('guest.female') },
                    ]"
                    :disabled="loading"
                  />
                </UFormField>

                <UFormField :label="t('guest.dateOfBirth')" name="guest.dateOfBirth">
                  <HDate v-model="form.guest.dateOfBirth" />
                </UFormField>

                <UFormField :label="t('guest.placeOfBirth')" name="guest.placeOfBirth">
                  <UInput v-model="form.guest.placeOfBirth" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.nationalId')" name="guest.nationalId">
                  <UInput v-model="form.guest.nationalId" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.idPassportNumber')" name="guest.idNumber">
                  <UInput v-model="form.guest.idNumber" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.occupation')" name="guest.occupation">
                  <UInput v-model="form.guest.occupation" :disabled="loading" :placeholder="t('guest.defaultOccupation')" />
                </UFormField>

                <UFormField :label="t('guest.phone')" name="guest.phone">
                  <UInput v-model="form.guest.phone" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.postalCode')" name="guest.postalCode">
                  <UInput v-model="form.guest.postalCode" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('forms.address')" name="guest.address" class="md:col-span-3">
                  <UInput class="w-full" v-model="form.guest.address" :disabled="loading" />
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
                <UFormField :label="t('guest.reservationCode')" name="reservation.reservationCode" class="md:col-span-3">
                  <UInput v-model="form.reservation.reservationCode" :disabled="true" />
                </UFormField>

                <UFormField :label="t('guest.entryDate')" name="reservation.entryDate" required>
                  <HDate v-model="form.reservation.entryDate" />
                </UFormField>

                <UFormField :label="t('guest.departureDate')" name="reservation.departureDate">
                  <HDate v-model="form.reservation.departureDate" />
                </UFormField>

                <UFormField :label="t('guest.numberOfPeople')" name="reservation.numberOfPeople">
                  <UInput
                    type="number"
                    min="1"
                    v-model.number="form.reservation.numberOfPeople"
                    :disabled="loading"
                  />
                </UFormField>

                <UFormField :label="t('guest.durationOfStay')" name="reservation.durationOfStay">
                  <UInput
                    type="number"
                    min="1"
                    v-model.number="form.reservation.durationOfStay"
                    :disabled="loading"
                  />
                </UFormField>

                <UFormField :label="t('guest.origin')" name="reservation.origin">
                  <UInput v-model="form.reservation.origin" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.destination')" name="reservation.destination">
                  <UInput v-model="form.reservation.destination" :disabled="loading" :placeholder="t('guest.defaultDestination')" />
                </UFormField>

                <UFormField :label="t('guest.purposeOfTravel')" name="reservation.purposeOfTravel">
                  <UInput v-model="form.reservation.purposeOfTravel" :disabled="loading" :placeholder="t('guest.defaultPurposeOfTravel')" />
                </UFormField>

                <UFormField :label="t('guest.roomPrice')" name="reservation.roomPrice">
                  <UInput
                    type="number"
                    min="0"
                    v-model.number="form.reservation.roomPrice"
                    :disabled="loading"
                  />
                </UFormField>

                <UFormField
                  :label="t('guest.selectRoom')"
                  name="reservation.rooms"
                  class="md:col-span-3"
                >
                  <USelect
                    v-model="form.roomIds"
                    :items="rooms ?? []"
                    multiple
                    :disabled="loading"
                  />
                </UFormField>

                <UFormField :label="t('guest.breakfast')" name="reservation.breakfast">
                  <UCheckbox v-model="form.reservation.breakfast" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.guide')" name="reservation.guide">
                  <UCheckbox v-model="form.reservation.guide" :disabled="loading" />
                </UFormField>

                <UFormField
                  :label="t('guest.notes')"
                  name="reservation.notes"
                  class="md:col-span-3"
                >
                  <UTextarea
                    v-model="form.reservation.notes"
                    :rows="3"
                    class="w-full"
                    :disabled="loading"
                  />
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
                <UFormField :label="t('guest.cash')" name="payment.isCash">
                  <UCheckbox v-model="form.payment.isCash" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.agency')" name="payment.agency">
                  <UCheckbox v-model="form.payment.agency" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.referrer')" name="payment.referrer">
                  <UInput v-model="form.payment.referrer" :disabled="loading" />
                </UFormField>

                <UFormField :label="t('guest.contractType')" name="payment.contractType">
                  <UInput v-model="form.payment.contractType" :disabled="loading" />
                </UFormField>
              </div>
            </template>
          </UCollapsible>

          <!-- SECTION: Companions -->
          <UCollapsible>
            <template #default="{ open }">
              <UButton
                :label="t('guests.companions')"
                color="neutral"
                variant="outline"
                :trailing-icon="open ? 'i-lucide-chevron-up' : 'i-lucide-chevron-down'"
                size="lg"
                block
              />
            </template>
            <template #content>
              <div class="mt-2 rounded-md p-4 ring ring-accented/40 ring-inset">
                <div class="mb-4 flex items-center justify-between">
                  <span class="text-sm font-medium">{{ t("guests.companionsList") }}</span>
                  <UButton size="sm" color="primary" @click="showCompanionModal = true">
                    <UIcon name="i-lucide-plus" class="mr-1" />
                    {{ t("guests.addCompanion") }}
                  </UButton>
                </div>

                <div v-if="companions.length === 0" class="py-4 text-center text-gray-500">
                  {{ t("guests.noCompanions") }}
                </div>

                <UTable v-else :data="companions" :columns="companionColumns" striped>
                  <template #actions-cell="{ row }">
                    <UButton
                      variant="ghost"
                      color="error"
                      size="xs"
                      @click="removeCompanion(row.index)"
                    >
                      <UIcon name="i-lucide-trash-2" />
                    </UButton>
                  </template>
                </UTable>
              </div>
            </template>
          </UCollapsible>

          <div class="flex justify-end gap-3 pt-4">
            <UButton variant="outline" to="/guests" :disabled="loading">{{
              t("actions.cancel")
            }}</UButton>
            <UButton type="submit" color="primary" :loading="loading">{{
              t("guests.createGuest")
            }}</UButton>
          </div>
        </UForm>
      </UCard>
    </div>

    <!-- Right: Sticky Summary -->
    <div class="w-full lg:w-80">
      <UCard class="sticky top-4">
        <h2 class="mb-4 text-xl font-semibold">{{ t("guest.reservationSummary") }}</h2>

        <div class="space-y-2 text-sm">
          <div>
            <strong>{{ t("guest.summaryGuest") }}:</strong> {{ form.guest.firstName }}
            {{ form.guest.lastName }}
          </div>
          <div>
            <strong>{{ t("guest.summaryDates") }}:</strong>
            {{ formatDate(form.reservation?.entryDate) }} →
            {{ formatDate(form.reservation?.departureDate) }}
          </div>
          <div>
            <strong>{{ t("guest.summaryPeople") }}:</strong> {{ form.reservation?.numberOfPeople }}
          </div>
          <div>
            <strong>{{ t("guest.summaryPrice") }}:</strong> {{ form.reservation?.roomPrice }}
            {{ t("guest.perNight") }}
          </div>
          <div>
            <strong>{{ t("guest.payment") }}:</strong>
            <span v-if="form.payment?.isCash">{{ t("guest.cash") }}</span>
            <span v-else-if="form.payment?.agency">{{ t("guest.agency") }}</span>
            <span v-else>{{ t("guest.unspecified") }}</span>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Add Companion Modal -->
    <UModal v-model:open="showCompanionModal" :title="t('guests.addCompanion')">
      <template #content>
        <form @submit.prevent="addCompanion" class="p-4">
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <UFormField :label="t('forms.firstName')" name="newCompanion.firstName" required>
                <UInput v-model="newCompanion.firstName" :disabled="addingCompanion" />
              </UFormField>
              <UFormField :label="t('forms.lastName')" name="newCompanion.lastName" required>
                <UInput v-model="newCompanion.lastName" :disabled="addingCompanion" />
              </UFormField>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <UFormField :label="t('guest.nationalId')" name="newCompanion.nationalId">
                <UInput v-model="newCompanion.nationalId" :disabled="addingCompanion" />
              </UFormField>
              <UFormField :label="t('guest.idPassportNumber')" name="newCompanion.idNumber">
                <UInput v-model="newCompanion.idNumber" :disabled="addingCompanion" />
              </UFormField>
            </div>

            <UFormField :label="t('guests.companionRelation')" name="newCompanion.relation">
              <USelect
                v-model="newCompanion.relation"
                :items="relationOptions"
                :disabled="addingCompanion"
              />
            </UFormField>
          </div>
          <div class="mt-6 flex justify-end gap-3">
            <UButton variant="outline" type="button" @click="showCompanionModal = false">
              {{ t("actions.cancel") }}
            </UButton>
            <UButton type="submit" color="primary" :loading="addingCompanion">
              {{ t("guests.addCompanion") }}
            </UButton>
          </div>
        </form>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { FormSubmitEvent } from "@nuxt/ui";
import type z from "zod";
import { useAuthStore } from "~/stores/auth";
import { zGuestWithReservationRequest } from "~/utils/client/zod.gen";

definePageMeta({
  requiresRole: ["admin", "manager"],
});

const { t } = useI18n();

const loading = ref(false);

const generateRegisterNumber = () => {
  const now = new Date();
  const y = now.getFullYear().toString().slice(-2);
  const m = String(now.getMonth() + 1).padStart(2, "0");
  const d = String(now.getDate()).padStart(2, "0");
  const h = String(now.getHours()).padStart(2, "0");
  const min = String(now.getMinutes()).padStart(2, "0");
  const s = String(now.getSeconds()).padStart(2, "0");
  return `${y}${m}${d}-${h}${min}${s}`;
};

const schema = zGuestWithReservationRequest;
type Schema = z.output<typeof schema>;

const form = ref<Required<Schema> & { roomIds: number[] }>({
  roomIds: [],
  reservation: {
    rooms: [],
    reservationCode: generateRegisterNumber(),
    destination: t("guest.defaultDestination"),
    purposeOfTravel: t("guest.defaultPurposeOfTravel"),
  },
  guest: {
    gender: "male",
    nationality: t("guest.defaultNationality"),
    occupation: t("guest.defaultOccupation"),
  } as any,
  payment: {},
});

const companions = ref<
  Array<{
    firstName: string;
    lastName: string;
    nationalId: string;
    idNumber: string;
    relation: string;
  }>
>([]);

const showCompanionModal = ref(false);
const addingCompanion = ref(false);

const newCompanion = ref({
  firstName: "",
  lastName: "",
  nationalId: "",
  idNumber: "",
  relation: "",
});

const relationOptions = [
  { value: "spouse", label: "همسر" },
  { value: "child", label: "فرزند" },
  { value: "parent", label: "والد" },
  { value: "sibling", label: "خواهر/برادر" },
  { value: "relative", label: "خویشاوند" },
  { value: "friend", label: "دوست" },
  { value: "colleague", label: "همکار" },
  { value: "other", label: "سایر" },
];

const companionColumns = [
  { accessorKey: "firstName", header: t("forms.firstName") },
  { accessorKey: "lastName", header: t("forms.lastName") },
  { accessorKey: "nationalId", header: t("guest.nationalId") },
  { accessorKey: "idNumber", header: t("guest.idPassportNumber") },
  { accessorKey: "relation", header: t("guests.companionRelation") },
  { accessorKey: "actions", header: t("guests.columns.actions") },
];

const addCompanion = () => {
  if (newCompanion.value.firstName && newCompanion.value.lastName) {
    companions.value.push({ ...newCompanion.value });
    newCompanion.value = {
      firstName: "",
      lastName: "",
      nationalId: "",
      idNumber: "",
      relation: "",
    };
    showCompanionModal.value = false;
  }
};

const removeCompanion = (index: number) => {
  companions.value.splice(index, 1);
};

const { data: rooms } = useAsyncData(async () => {
  const res = await getApiRooms({});
  return res.data?.map((v) => ({
    value: v.id,
    label: `${v.name ?? v.roomNumber}`,
  }));
});

const selectedRooms = ref<number>();
const authStore = useAuthStore();

const handleSubmit = async (event: FormSubmitEvent<Schema>) => {
  loading.value = true;
  try {
    const body: Schema = {
      guest: {
        ...event.data.guest,
        hotelId: authStore.currentHotelId ?? "default",
        dateOfBirth: event.data.guest?.dateOfBirth
          ? new Date(event.data.guest.dateOfBirth).toISOString()
          : undefined,
      } as Schema["guest"],
      reservation: event.data.reservation,
      payment: event.data.payment,
      companions: companions.value,
    } as unknown as Schema;

    await postApiGuestsWithReservation({
      body,
    });

    navigateTo("/guests");
  } catch (error) {
    console.error("Failed to create guest:", error);
  } finally {
    loading.value = false;
  }
};

const formatDate = (date: string | undefined) => {
  if (!date) return "-";
  return new Date(date).toLocaleDateString("fa-IR");
};
</script>
<style scoped>
div[data-slot="root"]:has(input) {
  width: 100%;
}
</style>
