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
          @click="checkOut"
          :loading="actionLoading"
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
            <strong>{{ t("stays.entryDate") }}:</strong> {{ formatDate(stay.entryDate) }}
          </p>
          <p>
            <strong>{{ t("stays.departureDate") }}:</strong> {{ formatDate(stay.departureDate) }}
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
          <UButton size="xs" @click="showPayment = true">{{ t("stays.pay") }}</UButton>
        </div>
        <div v-else>{{ t("stays.noInvoice") }}</div>
      </UCard>
    </div>

    <!-- Invoice Items -->
    <UCard class="mt-6" v-if="invoice?.items?.length">
      <template #header>
        <h2 class="font-semibold">{{ t("stays.invoiceItems") }}</h2>
      </template>
      <UTable :rows="invoice.items" :columns="itemColumns">
        <template #actions-data="{ row }">
          <UButton
            v-if="((row as unknown as InvoiceItem).remainingAmount ?? 0) > 0"
            size="xs"
            @click="settleItem(row as unknown as InvoiceItem)"
            >{{ t("stays.settle") }}</UButton
          >
          <UBadge v-else color="success">{{ t("stays.settled") }}</UBadge>
        </template>
      </UTable>
    </UCard>

    <!-- Payment Modal -->
    <UModal v-model="showPayment">
      <UCard>
        <template #header>
          <h3 class="font-semibold">{{ t("stays.paymentTitle") }}</h3>
        </template>
        <UForm :state="paymentState" @submit="makePayment">
          <UFormGroup :label="t('stays.amount')" name="amount">
            <UInput v-model.number="paymentState.amount" type="number" />
          </UFormGroup>
          <UFormGroup :label="t('stays.paymentMethod')" name="paymentMethod" class="mt-2">
            <USelect v-model="paymentState.paymentMethod" :items="paymentMethodOptions" />
          </UFormGroup>
          <div class="mt-4 flex gap-2">
            <UButton type="submit" :loading="paymentLoading">{{ t("stays.pay") }}</UButton>
            <UButton variant="outline" @click="showPayment = false">{{
              t("actions.cancel")
            }}</UButton>
          </div>
        </UForm>
      </UCard>
    </UModal>

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

    <!-- Item Settlement Modal -->
    <UModal v-model="showItemSettlement">
      <UCard>
        <template #header>
          <h3 class="font-semibold">{{ t("stays.itemSettlement") }}</h3>
        </template>
        <UForm :state="itemSettlementState" @submit="doItemSettlement">
          <UFormGroup :label="t('stays.amount')" name="amount">
            <UInput v-model.number="itemSettlementState.amount" type="number" />
          </UFormGroup>
          <UFormGroup :label="t('stays.paymentMethod')" name="paymentMethod" class="mt-2">
            <USelect v-model="itemSettlementState.paymentMethod" :items="paymentMethodOptions" />
          </UFormGroup>
          <div class="mt-4 flex gap-2">
            <UButton type="submit" :loading="itemSettlementLoading">{{
              t("stays.settle")
            }}</UButton>
            <UButton variant="outline" @click="showItemSettlement = false">{{
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
  postApiStaysIdCheckOut,
  postApiStaysIdInvoiceItems,
  postApiStaysIdInvoicePay,
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
  key: `stay-rooms-${stayId}`,
  transform: (res) => (res as { data?: Room[] })?.data ?? [],
});

const { data: services } = useFetch("/api/services", {
  key: `stay-services-${stayId}`,
  transform: (res) => (res as { data?: Service[] })?.data ?? [],
});

const roomOptions = computed(() =>
  (rooms.value ?? []).map((r) => ({ label: r.roomNumber ?? "", value: r.id }))
);

const serviceOptions = computed(() =>
  (services.value ?? []).map((s) => ({ label: s.name ?? "", value: s.id }))
);

const itemColumns = [
  { accessorKey: "itemType", header: t("stays.itemType") },
  { accessorKey: "description", header: t("stays.description") },
  { accessorKey: "quantity", header: t("stays.quantity") },
  { accessorKey: "unitPrice", header: t("stays.unitPrice") },
  { accessorKey: "totalPrice", header: t("stays.totalPrice") },
  { accessorKey: "paidAmount", header: t("stays.paidAmount") },
  { accessorKey: "remainingAmount", header: t("stays.remainingAmount") },
  { id: "actions", header: t("actions.actions") },
];

const actionLoading = ref(false);
const showPayment = ref(false);
const paymentLoading = ref(false);
const paymentState = reactive({ amount: 0, paymentMethod: undefined as number | undefined });

const showRoomChange = ref(false);
const roomChangeLoading = ref(false);
const roomChangeState = reactive({ newRoomId: undefined as number | undefined });

const showService = ref(false);
const serviceLoading = ref(false);
const serviceState = reactive({
  serviceId: undefined as number | undefined,
  quantity: 1,
  description: "",
});

const showItemSettlement = ref(false);
const itemSettlementLoading = ref(false);
const itemSettlementState = reactive({ amount: 0, paymentMethod: undefined as number | undefined });
const currentItem = ref<InvoiceItem | null>(null);

const paymentMethodOptions = ref([
  { label: t("payment.cash"), value: 1 },
  { label: t("payment.card"), value: 2 },
  { label: t("payment.transfer"), value: 3 },
  { label: t("payment.contractingParty"), value: 4 },
  { label: t("payment.agency"), value: 5 },
]);

function formatDate(date: string | undefined) {
  if (!date) return "-";
  return new Date(date).toLocaleDateString("fa-IR");
}

async function checkIn() {
  actionLoading.value = true;
  try {
    const response = await postApiStaysIdCheckIn({ path: { id: stayId }, body: {} });
    const res = response.data;
    if (res?.earlyCheckInPrompt) {
      alert(res.promptMessage);
    }
    refreshStay();
    refreshInvoice();
  } catch (e) {
    console.error(e);
  } finally {
    actionLoading.value = false;
  }
}

async function checkOut() {
  actionLoading.value = true;
  try {
    await postApiStaysIdCheckOut({ path: { id: stayId }, body: {} });
    refreshStay();
  } catch (e) {
    console.error(e);
  } finally {
    actionLoading.value = false;
  }
}

async function makePayment() {
  paymentLoading.value = true;
  try {
    await postApiStaysIdInvoicePay({
      path: { id: stayId },
      body: { amount: paymentState.amount, paymentMethod: paymentState.paymentMethod },
    });
    showPayment.value = false;
    refreshInvoice();
  } catch (e) {
    console.error(e);
  } finally {
    paymentLoading.value = false;
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

function settleItem(item: InvoiceItem) {
  currentItem.value = item;
  itemSettlementState.amount = item.remainingAmount ?? 0;
  showItemSettlement.value = true;
}

async function doItemSettlement() {
  itemSettlementLoading.value = true;
  try {
    await postApiStaysIdInvoiceItems({
      path: { id: stayId },
      body: {
        itemType: currentItem.value?.itemType,
        description: (currentItem.value?.description ?? "") + " (settlement)",
        quantity: 1,
        unitPrice: 0,
        totalPrice: 0,
      },
    });
    await postApiStaysIdInvoicePay({
      path: { id: stayId },
      body: {
        amount: itemSettlementState.amount,
        paymentMethod: itemSettlementState.paymentMethod,
      },
    });
    showItemSettlement.value = false;
    refreshInvoice();
  } catch (e) {
    console.error(e);
  } finally {
    itemSettlementLoading.value = false;
  }
}
</script>
