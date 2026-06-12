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
            <div v-for="stay in settlement.stays" :key="stay.id" class="mb-4 border-b pb-4 last:border-0">
              <div class="flex items-center justify-between mb-2">
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
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b text-left text-gray-500">
                    <th class="pb-1">{{ t("accounting.item") }}</th>
                    <th class="pb-1 text-right">{{ t("accounting.amount") }}</th>
                    <th class="pb-1 text-right">{{ t("accounting.paid") }}</th>
                    <th class="pb-1 text-right">{{ t("accounting.remaining") }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in stay.items" :key="item.id" class="border-b">
                    <td class="py-1">{{ item.description || item.itemType }}</td>
                    <td class="py-1 text-right">${{ item.totalPrice?.toFixed(2) }}</td>
                    <td class="py-1 text-right">${{ item.paidAmount?.toFixed(2) }}</td>
                    <td class="py-1 text-right">${{ item.remainingAmount?.toFixed(2) }}</td>
                  </tr>
                </tbody>
              </table>
              <div class="mt-2 flex justify-between font-medium">
                <span>{{ t("accounting.stayTotal") }}:</span>
                <span>${{ stay.totalAmount?.toFixed(2) }}</span>
              </div>
              <div class="flex justify-between text-sm text-gray-500">
                <span>{{ t("accounting.paid") }}:</span>
                <span>${{ stay.paidAmount?.toFixed(2) }}</span>
              </div>
              <div class="flex justify-between text-sm text-red-600">
                <span>{{ t("accounting.remaining") }}:</span>
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
            <table class="w-full">
              <thead>
                <tr class="border-b text-left text-sm text-gray-500">
                  <th class="pb-2">{{ t("parking.licensePlate") }}</th>
                  <th class="pb-2">{{ t("parking.entryTime") }}</th>
                  <th class="pb-2 text-right">{{ t("parking.hours") }}</th>
                  <th class="pb-2 text-right">{{ t("accounting.amount") }}</th>
                  <th class="pb-2 text-right">{{ t("accounting.paid") }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="pt in settlement.parkingTransactions" :key="pt.id" class="border-b">
                  <td class="py-3">{{ pt.licensePlate }}</td>
                  <td class="py-3">{{ formatDateTime(pt.entryTime) }}</td>
                  <td class="py-3 text-right">{{ pt.hoursParked?.toFixed(1) }}</td>
                  <td class="py-3 text-right">${{ pt.amountDue?.toFixed(2) }}</td>
                  <td class="py-3 text-right">${{ pt.amountPaid?.toFixed(2) }}</td>
                </tr>
              </tbody>
            </table>
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
            <table class="w-full">
              <thead>
                <tr class="border-b text-left text-sm text-gray-500">
                  <th class="pb-2">{{ t("accounting.billDate") }}</th>
                  <th class="pb-2 text-right">{{ t("accounting.amount") }}</th>
                  <th class="pb-2">{{ t("accounting.external") }}</th>
                  <th class="pb-2">{{ t("common.notes") }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="bill in settlement.restaurantBills" :key="bill.id" class="border-b">
                  <td class="py-3">{{ formatDate(bill.billDate) }}</td>
                  <td class="py-3 text-right">${{ bill.totalAmount?.toFixed(2) }}</td>
                  <td class="py-3">
                    <UBadge v-if="bill.isExternal" color="warning" variant="soft" size="sm">
                      {{ t("accounting.external") }}
                    </UBadge>
                    <span v-else class="text-sm text-gray-500">---</span>
                  </td>
                  <td class="max-w-xs truncate py-3">{{ bill.notes || "---" }}</td>
                </tr>
              </tbody>
            </table>
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
            <UButton variant="outline" @click="showPaymentForm = false">{{ t("actions.cancel") }}</UButton>
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
import {
  getApiAccountingPaymentMethods,
  getApiGuestsIdSettle,
  postApiGuestsIdSettle,
  postApiGuestsIdCheckout,
} from "~/utils/client";

const { t } = useI18n();
const route = useRoute();
const guestId = route.params.id as string;

const loading = ref(true);
const paying = ref(false);
const checkingOut = ref(false);
const showPaymentForm = ref(false);

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
  items: Array<{
    id: number;
    description: string;
    itemType: string;
    totalPrice: number;
    paidAmount: number;
    remainingAmount: number;
  }>;
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

const settlement = ref<Settlement | null>(null);
const paymentMethods = ref<PaymentMethod[]>([]);

const paymentForm = reactive({
  amount: 0,
  paymentMethod: undefined as number | undefined,
  reference: "",
  notes: "",
});

const paymentMethodOptions = computed(() => {
  return paymentMethods.value.map((pm) => ({
    value: pm.id,
    label: pm.label,
  }));
});

const fetchPaymentMethods = async () => {
  try {
    const res = await getApiAccountingPaymentMethods();
    paymentMethods.value = res.data?.data ?? [];
  } catch (e) {
    console.error("Failed to fetch payment methods:", e);
  }
};

const fetchSettlement = async () => {
  loading.value = true;
  try {
    const response = await getApiGuestsIdSettle({ path: { id: guestId } });
    settlement.value = response.data as Settlement;
  } catch (error) {
    console.error("Failed to fetch settlement:", error);
  } finally {
    loading.value = false;
  }
};

const handlePayment = async () => {
  if (paymentForm.amount <= 0) return;
  if (!paymentForm.paymentMethod) return;

  paying.value = true;
  try {
    // Gather invoice IDs for active stays
    const invoiceIds = settlement.value?.stays
      ?.filter((s) => s.remainingAmount > 0)
      ?.map((s) => s.id) ?? [];

    const parkingTxnIds = settlement.value?.parkingTransactions
      ?.filter((p) => p.amountDue > p.amountPaid)
      ?.map((p) => p.id) ?? [];

    const restaurantBillIds = settlement.value?.restaurantBills
      ?.map((b) => b.id) ?? [];

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
    await fetchSettlement();
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

onMounted(() => {
  fetchPaymentMethods();
  fetchSettlement();
});
</script>
