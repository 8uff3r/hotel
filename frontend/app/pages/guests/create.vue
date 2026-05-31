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
        <UForm
          @submit="handleSubmit"
          :state="form"
          :schema="createSchema"
          class="flex flex-col gap-2"
        >
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
                <GuestFormFields v-model="form" :loading />
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
                <ReservationFormFields v-model="form" :loading />
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
            <UButton
              type="submit"
              color="primary"
              :loading="loading"
              @click="
                () => {
                  console.log(createSchema.parse(form));
                }
              "
              >{{ t("guests.createGuest") }}</UButton
            >
          </div>
        </UForm>
      </UCard>
    </div>

    <!-- Right: Sticky Summary -->
    <div class="w-80 max-lg:hidden">
      <StickySummary :value="form" />
    </div>

    <!-- Add Companion Modal -->

    <AddCompanionModal
      v-model:state="newCompanion"
      v-model:open="showCompanionModal"
      :loading="addingCompanion"
      @submit="addCompanion"
    />
  </div>
</template>

<script setup lang="ts">
import type { FormSubmitEvent } from "@nuxt/ui";
import { useAuthStore } from "~/stores/auth";
import AddCompanionModal from "./components/AddCompanionModal.vue";
import { type Companion, type CreateRequest, createSchema } from "./utils";
import StickySummary from "./components/StickySummary.vue";
import GuestFormFields from "./components/GuestFormFields.vue";
import ReservationFormFields from "./components/ReservationFormFields.vue";
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

const form = ref<Required<CreateRequest> & { roomIds: number[] }>({
  roomIds: [],
  reservation: {
    rooms: [],
    reservationCode: generateRegisterNumber(),
    destination: t("guest.defaultDestination"),
    purposeOfTravel: t("guest.defaultPurposeOfTravel"),
  },
  guest: {
    gender: "male",
  } as any,
  payment: {},
  companions: [],
});

const companions = ref<Companion[]>([]);

const showCompanionModal = ref(false);
const addingCompanion = ref(false);

const newCompanion = ref<Companion>({});

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
    newCompanion.value = {};
    showCompanionModal.value = false;
  }
};

const removeCompanion = (index: number) => {
  companions.value.splice(index, 1);
};

const authStore = useAuthStore();

const handleSubmit = async (event: FormSubmitEvent<CreateRequest>) => {
  loading.value = true;
  try {
    const body: CreateRequest = {
      guest: {
        ...event.data.guest,
        hotelId: authStore.currentHotelId ?? "default",
        dateOfBirth: event.data.guest?.dateOfBirth
          ? new Date(event.data.guest.dateOfBirth).toISOString()
          : undefined,
      } as CreateRequest["guest"],
      reservation: event.data.reservation,
      payment: event.data.payment,
      companions: companions.value,
    } as unknown as CreateRequest;

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
</script>
<style scoped>
div[data-slot="root"]:has(input) {
  width: 100%;
}
</style>
