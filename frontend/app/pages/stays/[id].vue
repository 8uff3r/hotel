<template>
  <div class="p-6" v-if="stay">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl font-bold">{{ t("stays.detailTitle") }} {{ stay.acceptanceId }}</h1>
      <div class="flex gap-2">
        <UButton v-if="stay.status?.slug === 'waiting'" @click="checkIn" :loading="actionLoading">
          {{ t("stays.checkIn") }}
        </UButton>
        <UButton
          v-if="stay.status?.slug === 'resident'"
          :to="`/guests/${stay.guestId}/settle`"
          color="error"
        >
          {{ t("stays.checkOut") }}
        </UButton>
        <UButton
          v-if="stay.status?.slug === 'resident'"
          @click="showRoomChange = true"
          variant="soft"
        >
          {{ t("stays.changeRoom") }}
        </UButton>
        <UButton v-if="stay.status?.slug === 'resident'" @click="showService = true" variant="soft">
          {{ t("stays.addService") }}
        </UButton>
        <UButton v-if="stay.status?.slug === 'resident'" @click="showExtend = true" variant="soft">
          {{ t("stays.changeDuration") }}
        </UButton>
        <UButton to="/stays">{{ t("actions.back") }}</UButton>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-6">
      <UCard>
        <template #header>
          <h2 class="font-semibold">{{ t("stays.guestInfo") }}</h2>
        </template>
        <div class="space-y-2">
          <p>
            <strong>{{ t("guests.name") }}:</strong> {{ stay.guest?.firstName }}
            {{ stay.guest?.lastName }}
          </p>
          <p>
            <strong>{{ t("guests.nationalId") }}:</strong> {{ stay.guest?.nationalId }}
          </p>
          <p>
            <strong>{{ t("guests.phone") }}:</strong> {{ stay.guest?.phone }}
          </p>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <h2 class="font-semibold">{{ t("stays.roomInfo") }}</h2>
        </template>
        <div class="space-y-2">
          <p>
            <strong>{{ t("rooms.roomNumber") }}:</strong> {{ stay.room?.roomNumber }}
          </p>
          <p>
            <strong>{{ t("rooms.capacity") }}:</strong> {{ stay.room?.capacity }}
          </p>
          <p>
            <strong>{{ t("rooms.type") }}:</strong> {{ stay.room?.roomType?.label }}
          </p>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <h2 class="font-semibold">{{ t("stays.stayInfo") }}</h2>
        </template>
        <div class="space-y-2">
          <p>
            <strong>{{ t("stays.entryDate") }}:</strong> {{ formatDate(stay.entryDate ?? '') }}
          </p>
          <p>
            <strong>{{ t("stays.departureDate") }}:</strong> {{ formatDate(stay.departureDate ?? '') }}
          </p>
          <p>
            <strong>{{ t("stays.numberOfPeople") }}:</strong> {{ stay.numberOfPeople }}
          </p>
          <p>
            <strong>{{ t("stays.roomPrice") }}:</strong> {{ stay.roomPrice }}
          </p>
          <p v-if="stay.earlyCheckInFee && stay.earlyCheckInFee > 0">
            <strong>{{ t("stays.earlyCheckInFee") }}:</strong> {{ stay.earlyCheckInFee }}
          </p>
          <p v-if="stay.halfDayFee && stay.halfDayFee > 0">
            <strong>{{ t("stays.halfDayFee") }}:</strong> {{ stay.halfDayFee }}
          </p>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <h2 class="font-semibold">{{ t("stays.invoice") }}</h2>
        </template>
        <div v-if="invoice" class="space-y-2">
          <p>
            <strong>{{ t("stays.totalAmount") }}:</strong> {{ invoice.totalAmount }}
          </p>
          <p>
            <strong>{{ t("stays.paidAmount") }}:</strong> {{ invoice.paidAmount }}
          </p>
          <p>
            <strong>{{ t("stays.remainingAmount") }}:</strong> {{ invoice.remainingAmount }}
          </p>
          <p>
            <strong>{{ t("stays.paymentStatus") }}:</strong> {{ invoice.paymentStatus }}
          </p>
          <UButton size="xs" :to="`/guests/${stay.guestId}/settle`">{{ t("stays.goToSettlement") }}</UButton>
        </div>
        <div v-else>{{ t("stays.noInvoice") }}</div>
      </UCard>
    </div>

    <!-- Invoice Items -->
    <UCard class="mt-6" v-if="invoice?.items?.length">
      <template #header>
        <h2 class="font-semibold">{{ t("stays.invoiceItems") }}</h2>
      </template>
      <UTable :data="invoice.items" :columns="itemColumns">
        <template #actions-cell="{ row }">
          <UBadge v-if="((row.original as unknown as InvoiceItem).remainingAmount ?? 0) > 0" color="warning">
            {{ t("stays.unpaid") }}
          </UBadge>
          <UBadge v-else color="success">{{ t("stays.settled") }}</UBadge>
        </template>
      </UTable>
    </UCard>

    <!-- Room Change Modal -->
    <UModal v-model="showRoomChange">
      <UCard>
        <template #header>
          <h3 class="font-semibold">{{ t("stays.changeRoom") }}</h3>
        </template>
        <UForm :state="roomChangeState" @submit="doRoomChange">
          <UFormGroup :label="t('stays.newRoom')" name="newRoomId">
            <USelect v-model="roomChangeState.newRoomId" :items="roomOptions" />
          </UFormGroup>
          <div class="mt-4 flex gap-2">
            <UButton type="submit" :loading="roomChangeLoading">{{ t("actions.change") }}</UButton>
            <UButton variant="outline" @click="showRoomChange = false">{{
              t("actions.cancel")
            }}</UButton>
          </div>
        </UForm>
      </UCard>
    </UModal>

    <!-- Add Service Modal -->
    <UModal v-model="showService">
      <UCard>
        <template #header>
          <h3 class="font-semibold">{{ t("stays.addService") }}</h3>
        </template>
        <UForm :state="serviceState" @submit="doAddService">
          <UFormGroup :label="t('stays.service')" name="serviceId">
            <USelect v-model="serviceState.serviceId" :items="serviceOptions" />
          </UFormGroup>
          <UFormGroup :label="t('stays.quantity')" name="quantity" class="mt-2">
            <UInput v-model.number="serviceState.quantity" type="number" />
          </UFormGroup>
          <UFormGroup :label="t('stays.description')" name="description" class="mt-2">
            <UTextarea v-model="serviceState.description" />
          </UFormGroup>
          <div class="mt-4 flex gap-2">
            <UButton type="submit" :loading="serviceLoading">{{ t("stays.add") }}</UButton>
            <UButton variant="outline" @click="showService = false">{{
              t("actions.cancel")
            }}</UButton>
          </div>
        </UForm>
      </UCard>
    </UModal>

    <!-- Change Duration Modal -->
    <UModal v-model="showExtend">
      <UCard>
        <template #header>
          <h3 class="font-semibold">{{ t("stays.changeDuration") }}</h3>
        </template>
        <UForm :state="extendState" @submit="doChangeDuration">
          <UFormGroup :label="t('stays.newDuration')" name="durationOfStay" required>
            <UInput v-model.number="extendState.durationOfStay" type="number" min="1" />
          </UFormGroup>
          <div class="mt-4 flex gap-2">
            <UButton type="submit" :loading="extendLoading">{{
              t("stays.updateDuration")
            }}</UButton>
            <UButton variant="outline" @click="showExtend = false">{{
              t("actions.cancel")
            }}</UButton>
          </div>
        </UForm>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import {
  postApiStaysIdChangeRoom,
  postApiStaysIdCheckIn,
  postApiStaysIdServices,
} from "~/utils/client";
import type { Invoice, InvoiceItem, Room, Service, Stay } from "~/utils/client";

const route = useRoute();
const stayId = route.params.id as string;
const { t } = useI18n();

const { data: stay, refresh: refreshStay } = useFetch(`/api/stays/${stayId}`, {
  key: `stay-${stayId}`,
  transform: (res) => (res as { data?: Stay })?.data,
});

const { data: invoice, refresh: refreshInvoice } = useFetch(`/api/stays/${stayId}/invoice`, {
  key: `stay-invoice-${stayId}`,
  transform: (res) => (res as { data?: Invoice })?.data,
});

const { data: rooms } = useFetch("/api/rooms", {
  key: "rooms",
  transform: (res) =>
    (res as { data?: Room[] })?.data?.map((r) => ({ value: r.id, label: r.roomNumber })) ?? [],
});

const { data: services } = useFetch("/api/services", {
  key: "services",
  transform: (res) =>
    (res as { data?: Service[] })?.data?.map((s) => ({ value: s.id, label: s.name })) ?? [],
});

const { data: paymentMethods } = useFetch("/api/accounting/payment-methods", {
  key: "payment-methods",
  transform: (res) =>
    (res as { data?: { id: number; label: string }[] })?.data?.map((m) => ({
      value: m.id,
      label: m.label,
    })) ?? [],
});

const actionLoading = ref(false);
const roomChangeLoading = ref(false);
const serviceLoading = ref(false);
const extendLoading = ref(false);

const showRoomChange = ref(false);
const showService = ref(false);
const showExtend = ref(false);

const roomChangeState = reactive({ newRoomId: undefined as number | undefined });
const serviceState = reactive({
  serviceId: undefined as number | undefined,
  quantity: 1,
  description: "",
});
const extendState = reactive({ durationOfStay: 1 });

const itemColumns = [
  { accessorKey: "itemType", header: t("stays.itemType") },
  { accessorKey: "description", header: t("stays.description") },
  { accessorKey: "quantity", header: t("stays.quantity") },
  { accessorKey: "unitPrice", header: t("stays.unitPrice") },
  { accessorKey: "totalPrice", header: t("stays.totalPrice") },
  { accessorKey: "paidAmount", header: t("stays.paidAmount") },
  { accessorKey: "remainingAmount", header: t("stays.remainingAmount") },
  { accessorKey: "actions", header: t("actions.actions") },
];

const roomOptions = computed(() => rooms.value ?? []);
const serviceOptions = computed(() => services.value ?? []);
const paymentMethodOptions = computed(() => paymentMethods.value ?? []);

async function checkIn() {
  actionLoading.value = true;
  try {
    await postApiStaysIdCheckIn({ path: { id: stayId }, body: {} });
    refreshStay();
  } catch (e) {
    console.error(e);
  } finally {
    actionLoading.value = false;
  }
}

async function doRoomChange() {
  roomChangeLoading.value = true;
  try {
    await postApiStaysIdChangeRoom({
      path: { id: stayId },
      body: { newRoomId: roomChangeState.newRoomId },
    });
    showRoomChange.value = false;
    refreshStay();
  } catch (e) {
    console.error(e);
  } finally {
    roomChangeLoading.value = false;
  }
}

async function doAddService() {
  serviceLoading.value = true;
  try {
    await postApiStaysIdServices({
      path: { id: stayId },
      body: {
        serviceId: serviceState.serviceId,
        quantity: serviceState.quantity,
        description: serviceState.description,
      },
    });
    showService.value = false;
    refreshInvoice();
  } catch (e) {
    console.error(e);
  } finally {
    serviceLoading.value = false;
  }
}

async function doChangeDuration() {
  extendLoading.value = true;
  try {
    await $fetch(`/api/stays/${stayId}/change-duration`, {
      method: "POST",
      body: { durationOfStay: extendState.durationOfStay },
    });
    showExtend.value = false;
    refreshStay();
    refreshInvoice();
  } catch (e) {
    console.error(e);
  } finally {
    extendLoading.value = false;
  }
}

function formatDate(date: string | Date) {
  if (!date) return "-";
  const d = typeof date === "string" ? new Date(date) : date;
  return d.toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}
</script>
