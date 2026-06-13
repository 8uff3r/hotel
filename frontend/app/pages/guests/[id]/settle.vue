<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" :to="`/guests/${guestId}`" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        {{ t("actions.backToGuest") }}
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("accounting.settleAccount") }}
      </h1>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <ULoader size="lg" />
    </div>

    <div v-else class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <div class="space-y-6 lg:col-span-2">
        <!-- Stays / Room Charges -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">{{ t("accounting.roomCharges") }}</h3>
          </template>
          <div v-if="settlement?.stays?.length">
            <div
              v-for="stay in settlement.stays"
              :key="stay.id"
              class="mb-4 border-b pb-4 last:border-0"
            >
              <div class="mb-2 flex items-center justify-between">
                <div>
                  <p class="font-medium">{{ stay.acceptanceId }}</p>
                  <p class="text-sm text-gray-500">
                    {{ formatDate(stay.entryDate) }} - {{ formatDate(stay.departureDate) }}
                  </p>
                </div>
                <UBadge variant="soft">
                  {{ stay.statusLabel || stay.status }}
                </UBadge>
              </div>
              <UTable :data="stay.items" :columns="itemColumns" striped>
                <template #description-cell="{ row }">
                  {{ row.original.description || row.original.itemType }}
                </template>
                <template #totalPrice-cell="{ row }">
                  ${{ row.original.totalPrice?.toFixed(2) }}
                </template>
                <template #paidAmount-cell="{ row }">
                  ${{ row.original.paidAmount?.toFixed(2) }}
                </template>
                <template #remainingAmount-cell="{ row }">
                  ${{ row.original.remainingAmount?.toFixed(2) }}
                </template>
              </UTable>
              <div class="mt-2 flex justify-between font-medium">
                <span>{{ t("accounting.stayTotal") }}:</span>
                <span>${{ stay.totalAmount?.toFixed(2) }}</span>
              </div>
              <div class="flex justify-between text-sm text-gray-500">
                <span>{{ t("accounting.paid") }}:</span>
                <span>${{ stay.paidAmount?.toFixed(2) }}</span>
              </div>
              <div class="flex justify-between text-sm text-red-600">
                <span>{{ t("accounting.remainingAmount") }}:</span>
                <span>${{ stay.remainingAmount?.toFixed(2) }}</span>
              </div>
            </div>
          </div>
          <div v-else class="py-4 text-center text-gray-500">
            {{ t("accounting.noRoomCharges") }}
          </div>
        </UCard>

        <!-- Parking Charges -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">{{ t("accounting.parkingCharges") }}</h3>
          </template>
          <div v-if="settlement?.parkingTransactions?.length">
            <UTable
              :data="settlement.parkingTransactions"
              :columns="parkingColumns"
              striped
            >
              <template #entryTime-cell="{ row }">
                {{ formatDateTime(row.original.entryTime) }}
              </template>
              <template #hoursParked-cell="{ row }">
                {{ row.original.hoursParked?.toFixed(1) }}
              </template>
              <template #amountDue-cell="{ row }">
                ${{ row.original.amountDue?.toFixed(2) }}
              </template>
              <template #amountPaid-cell="{ row }">
                ${{ row.original.amountPaid?.toFixed(2) }}
              </template>
            </UTable>
          </div>
          <div v-else class="py-4 text-center text-gray-500">
            {{ t("accounting.noParkingCharges") }}
          </div>
        </UCard>

        <!-- Restaurant Charges -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">{{ t("accounting.restaurantCharges") }}</h3>
          </template>
          <div v-if="settlement?.restaurantBills?.length">
            <UTable
              :data="settlement.restaurantBills"
              :columns="restaurantColumns"
              striped
            >
              <template #billDate-cell="{ row }">
                {{ formatDate(row.original.billDate) }}
              </template>
              <template #totalAmount-cell="{ row }">
                ${{ row.original.totalAmount?.toFixed(2) }}
              </template>
              <template #isExternal-cell="{ row }">
                <UBadge
                  v-if="row.original.isExternal"
                  color="warning"
                  variant="soft"
                  size="sm"
                >
                  {{ t("accounting.external") }}
                </UBadge>
                <span v-else class="text-sm text-gray-500">---</span>
              </template>
              <template #notes-cell="{ row }">
                <span class="max-w-xs truncate">{{ row.original.notes || "---" }}</span>
              </template>
            </UTable>
          </div>
          <div v-else class="py-4 text-center text-gray-500">
            {{ t("accounting.noRestaurantCharges") }}
          </div>
        </UCard>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">{{ t("accounting.summary") }}</h3>
          </template>
          <div class="space-y-3">
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t("accounting.roomTotal") }}</span>
              <span class="font-medium">${{ settlement?.totalRoom?.toFixed(2) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t("accounting.parkingTotal") }}</span>
              <span class="font-medium">${{ settlement?.totalParking?.toFixed(2) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t("accounting.restaurantTotal") }}</span>
              <span class="font-medium">${{ settlement?.totalRestaurant?.toFixed(2) }}</span>
            </div>
            <div class="border-t pt-3">
              <div class="flex justify-between">
                <span class="font-medium">{{ t("accounting.totalDue") }}</span>
                <span class="text-xl font-bold">${{ settlement?.totalDue?.toFixed(2) }}</span>
              </div>
            </div>
            <div class="flex justify-between text-green-600">
              <span>{{ t("accounting.totalPaid") }}</span>
              <span>${{ settlement?.totalPaid?.toFixed(2) }}</span>
            </div>
            <div class="border-t pt-3">
              <div class="flex justify-between">
                <span class="font-medium text-red-600">{{ t("accounting.balance") }}</span>
                <span class="text-xl font-bold text-red-600"
                  >${{ settlement?.balance?.toFixed(2) }}</span
                >
              </div>
            </div>
          </div>
        </UCard>

        <!-- Pay Now Button (only if balance > 0) -->
        <UButton
          v-if="settlement && settlement.balance > 0"
          color="success"
          block
          @click="showPaymentForm = true"
        >
          <UIcon name="i-lucide-credit-card" class="mr-2" />
          {{ t("accounting.payNow") }}
        </UButton>

        <!-- Checkout Button (only if balance == 0) -->
        <UButton
          v-if="settlement && settlement.canCheckout"
          color="error"
          block
          @click="handleCheckout"
          :loading="checkingOut"
        >
          <UIcon name="i-lucide-log-out" class="mr-2" />
          {{ t("accounting.checkout") }}
        </UButton>

        <!-- Payment Form Modal -->
        <UModal v-model="showPaymentForm">
          <template #body>
            <div class="space-y-4">
              <h3 class="text-lg font-semibold">{{ t("accounting.recordPayment") }}</h3>
              <div class="grid grid-cols-1 gap-4">
                <UFormField :label="t('accounting.amount')" name="amount">
                  <UInput v-model.number="paymentForm.amount" type="number" min="0" step="0.01" />
                </UFormField>
                <UFormField :label="t('accounting.paymentMethod')" name="paymentMethod">
                  <USelect v-model="paymentForm.paymentMethod" :items="paymentMethodOptions" />
                </UFormField>
                <UFormField :label="t('accounting.reference')" name="reference">
                  <UInput v-model="paymentForm.reference" />
                </UFormField>
                <UFormField :label="t('common.notes')" name="notes">
                  <UTextarea v-model="paymentForm.notes" :rows="2" />
                </UFormField>
              </div>
            </div>
          </template>
          <template #footer>
            <UButton variant="outline" @click="showPaymentForm = false">{{
              t("actions.cancel")
            }}</UButton>
            <UButton color="success" @click="handlePayment" :loading="paying">
              {{ t("accounting.recordPayment") }}
            </UButton>
          </template>
        </UModal>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import {
  getApiAccountingPaymentMethods,
  getApiGuestsIdSettle,
  postApiGuestsIdSettle,
  postApiGuestsIdCheckout,
} from "~/utils/client";

const { t } = useI18n();
const route = useRoute();
const guestId = route.params.id as string;
const qc = useQueryCache();

const { data: settlement, isPending: loading } = useQuery({
  key: () => ["guests", guestId, "settle"],
  query: async () => {
    const res = await getApiGuestsIdSettle({ path: { id: guestId } });
    return res.data as Settlement;
  },
});

const { data: paymentMethods } = useQuery({
  key: ["payment-methods", "list"],
  query: async () => {
    const res = await getApiAccountingPaymentMethods();
    return res.data?.data ?? [];
  },
});

const paying = ref(false);
const checkingOut = ref(false);
const showPaymentForm = ref(false);

interface StayItem {
  id: number;
  description: string;
  itemType: string;
  totalPrice: number;
  paidAmount: number;
  remainingAmount: number;
}

interface StayInvoiceSettlement {
  id: number;
  acceptanceId: string;
  entryDate: string;
  departureDate: string;
  status: string;
  statusLabel: string;
  totalAmount: number;
  paidAmount: number;
  remainingAmount: number;
  items: StayItem[];
}

interface ParkingSettlement {
  id: number;
  licensePlate: string;
  entryTime: string;
  exitTime: string;
  hoursParked: number;
  rateApplied: number;
  amountDue: number;
  amountPaid: number;
  status: string;
}

interface RestaurantSettlement {
  id: number;
  billDate: string;
  totalAmount: number;
  notes: string;
  isExternal: boolean;
}

interface Settlement {
  stays: StayInvoiceSettlement[];
  parkingTransactions: ParkingSettlement[];
  restaurantBills: RestaurantSettlement[];
  totalRoom: number;
  totalParking: number;
  totalRestaurant: number;
  totalDue: number;
  totalPaid: number;
  balance: number;
  canCheckout: boolean;
}

interface PaymentMethod {
  id?: number;
  slug?: string;
  label?: string;
}

const paymentForm = reactive({
  amount: 0,
  paymentMethod: undefined as number | undefined,
  reference: "",
  notes: "",
});

const paymentMethodOptions = computed(() => {
  return (
    (paymentMethods.value as PaymentMethod[] | undefined)?.map((pm) => ({
      value: pm.id,
      label: pm.label,
    })) ?? []
  );
});

const itemColumns: TableColumn<StayItem>[] = [
  { accessorKey: "description", header: t("accounting.item") },
  { accessorKey: "itemType", header: t("stays.itemType") },
  { accessorKey: "totalPrice", header: t("accounting.amount") },
  { accessorKey: "paidAmount", header: t("accounting.paid") },
  { accessorKey: "remainingAmount", header: t("accounting.remainingAmount") },
];

const parkingColumns: TableColumn<ParkingSettlement>[] = [
  { accessorKey: "licensePlate", header: t("parking.licensePlate") },
  { accessorKey: "entryTime", header: t("parking.entryTime") },
  { accessorKey: "hoursParked", header: t("parking.hours") },
  { accessorKey: "amountDue", header: t("accounting.amount") },
  { accessorKey: "amountPaid", header: t("accounting.paid") },
];

const restaurantColumns: TableColumn<RestaurantSettlement>[] = [
  { accessorKey: "billDate", header: t("accounting.billDate") },
  { accessorKey: "totalAmount", header: t("accounting.amount") },
  { accessorKey: "isExternal", header: t("accounting.external") },
  { accessorKey: "notes", header: t("common.notes") },
];

const handlePayment = async () => {
  if (paymentForm.amount <= 0) return;
  if (!paymentForm.paymentMethod) return;

  paying.value = true;
  try {
    const invoiceIds =
      settlement.value?.stays?.filter((s) => s.remainingAmount > 0)?.map((s) => s.id) ?? [];

    const parkingTxnIds =
      settlement.value?.parkingTransactions
        ?.filter((p) => p.amountDue > p.amountPaid)
        ?.map((p) => p.id) ?? [];

    const restaurantBillIds = settlement.value?.restaurantBills?.map((b) => b.id) ?? [];

    await postApiGuestsIdSettle({
      path: { id: guestId },
      body: {
        invoiceIds,
        parkingTxnIds,
        restaurantBillIds,
        amount: paymentForm.amount,
        paymentMethod: paymentForm.paymentMethod,
        reference: paymentForm.reference,
        notes: paymentForm.notes,
      },
    });

    showPaymentForm.value = false;
    paymentForm.amount = 0;
    paymentForm.reference = "";
    paymentForm.notes = "";
    qc.invalidateQueries({ key: ["guests", guestId, "settle"] });
  } catch (error) {
    console.error("Failed to settle:", error);
  } finally {
    paying.value = false;
  }
};

const handleCheckout = async () => {
  checkingOut.value = true;
  try {
    await postApiGuestsIdCheckout({
      path: { id: guestId },
      body: { paymentMethod: paymentForm.paymentMethod ?? 0 },
    });
    await navigateTo(`/guests/${guestId}`);
  } catch (error) {
    console.error("Failed to checkout:", error);
  } finally {
    checkingOut.value = false;
  }
};

const formatDate = (date: string) => {
  if (!date) return "-";
  return new Date(date).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
};

const formatDateTime = (date: string) => {
  if (!date) return "-";
  return new Date(date).toLocaleString("en-US", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};
</script>
