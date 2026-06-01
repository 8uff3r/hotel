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
        <!-- Room Charges -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold">{{ t("accounting.roomCharges") }}</h3>
              <UCheckbox
                v-if="settlement?.reservations?.length"
                v-model="selectAllRoom"
                @update:model-value="toggleSelectAll"
                :label="t('accounting.selectAll')"
              />
            </div>
          </template>
          <div v-if="settlement?.reservations?.length">
            <table class="w-full">
              <thead>
                <tr class="border-b text-left text-sm text-gray-500">
                  <th class="pb-2 w-10"></th>
                  <th class="pb-2">{{ t("reservations.code") }}</th>
                  <th class="pb-2">{{ t("reservations.checkIn") }}</th>
                  <th class="pb-2">{{ t("reservations.checkOut") }}</th>
                  <th class="pb-2">{{ t("reservations.status") }}</th>
                  <th class="pb-2 text-right">{{ t("accounting.amount") }}</th>
                  <th class="pb-2 text-right">{{ t("accounting.paid") }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="res in settlement.reservations" :key="res.id" class="border-b">
                  <td class="py-3">
                    <UCheckbox v-model="selectedReservationIds" :value="res.id" />
                  </td>
                  <td class="py-3">{{ res.reservationCode }}</td>
                  <td class="py-3">{{ formatDate(res.checkInDate) }}</td>
                  <td class="py-3">{{ formatDate(res.checkOutDate) }}</td>
                  <td class="py-3">
                    <UBadge :color="getStatusColor(res.status)" variant="soft" size="sm">
                      {{ res.statusLabel || res.status }}
                    </UBadge>
                  </td>
                  <td class="py-3 text-right">${{ res.roomPrice?.toFixed(2) }}</td>
                  <td class="py-3 text-right text-green-600">${{ res.paidAmount?.toFixed(2) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-else class="py-4 text-center text-gray-500">
            {{ t("accounting.noRoomCharges") }}
          </div>
        </UCard>

        <!-- Parking Charges -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold">{{ t("accounting.parkingCharges") }}</h3>
              <UCheckbox
                v-if="settlement?.parkingTransactions?.length"
                v-model="selectAllParking"
                @update:model-value="toggleSelectAllParking"
                :label="t('accounting.selectAll')"
              />
            </div>
          </template>
          <div v-if="settlement?.parkingTransactions?.length">
            <table class="w-full">
              <thead>
                <tr class="border-b text-left text-sm text-gray-500">
                  <th class="pb-2 w-10"></th>
                  <th class="pb-2">{{ t("parking.licensePlate") }}</th>
                  <th class="pb-2">{{ t("parking.entryTime") }}</th>
                  <th class="pb-2">{{ t("parking.exitTime") }}</th>
                  <th class="pb-2 text-right">{{ t("parking.hours") }}</th>
                  <th class="pb-2 text-right">{{ t("parking.rate") }}</th>
                  <th class="pb-2 text-right">{{ t("accounting.amount") }}</th>
                  <th class="pb-2 text-right">{{ t("accounting.paid") }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="pt in settlement.parkingTransactions" :key="pt.id" class="border-b">
                  <td class="py-3">
                    <UCheckbox v-model="selectedParkingTxnIds" :value="pt.id" />
                  </td>
                  <td class="py-3">{{ pt.licensePlate }}</td>
                  <td class="py-3">{{ formatDateTime(pt.entryTime) }}</td>
                  <td class="py-3">{{ formatDateTime(pt.exitTime) }}</td>
                  <td class="py-3 text-right">{{ pt.hoursParked?.toFixed(1) }}</td>
                  <td class="py-3 text-right">${{ pt.rateApplied?.toFixed(2) }}</td>
                  <td class="py-3 text-right">${{ pt.amountDue?.toFixed(2) }}</td>
                  <td class="py-3 text-right text-green-600">${{ pt.amountPaid?.toFixed(2) }}</td>
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
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold">{{ t("accounting.restaurantCharges") }}</h3>
              <UCheckbox
                v-if="settlement?.restaurantBills?.length"
                v-model="selectAllRestaurant"
                @update:model-value="toggleSelectAllRestaurant"
                :label="t('accounting.selectAll')"
              />
            </div>
          </template>
          <div v-if="settlement?.restaurantBills?.length">
            <table class="w-full">
              <thead>
                <tr class="border-b text-left text-sm text-gray-500">
                  <th class="pb-2 w-10"></th>
                  <th class="pb-2">{{ t("accounting.billDate") }}</th>
                  <th class="pb-2">{{ t("accounting.amount") }}</th>
                  <th class="pb-2">{{ t("accounting.external") }}</th>
                  <th class="pb-2">{{ t("common.notes") }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="bill in settlement.restaurantBills" :key="bill.id" class="border-b">
                  <td class="py-3">
                    <UCheckbox v-model="selectedRestaurantBillIds" :value="bill.id" />
                  </td>
                  <td class="py-3">{{ formatDate(bill.billDate) }}</td>
                  <td class="py-3">${{ bill.totalAmount?.toFixed(2) }}</td>
                  <td class="py-3">
                    <UBadge v-if="bill.isExternal" color="warning" variant="soft" size="sm">
                      {{ t("accounting.external") }}
                    </UBadge>
                    <span v-else class="text-sm text-gray-500">---</span>
                  </td>
                  <td class="py-3 max-w-xs truncate">{{ bill.notes || "---" }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-else class="py-4 text-center text-gray-500">
            {{ t("accounting.noRestaurantCharges") }}
          </div>
        </UCard>

        <!-- Payment Form -->
        <UCard v-if="showSettlementForm">
          <template #header>
            <h3 class="text-lg font-semibold">{{ t("accounting.recordPayment") }}</h3>
          </template>
          <form @submit.prevent="handleSettle">
            <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
              <UFormField :label="t('accounting.amount')" name="amount" required>
                <UInput v-model.number="form.amount" type="number" min="0" step="0.01" />
              </UFormField>

              <UFormField :label="t('accounting.paymentMethod')" name="paymentMethod">
                <USelect v-model="form.paymentMethod" :items="paymentMethodOptions" />
              </UFormField>

              <UFormField :label="t('accounting.reference')" name="reference" class="md:col-span-2">
                <UInput v-model="form.reference" :placeholder="t('accounting.receiptNumber')" />
              </UFormField>

              <UFormField :label="t('common.notes')" name="notes" class="md:col-span-2">
                <UTextarea v-model="form.notes" :rows="2" />
              </UFormField>
            </div>

            <div class="mt-6 flex justify-end gap-3">
              <UButton variant="outline" @click="showSettlementForm = false">
                {{ t("actions.cancel") }}
              </UButton>
              <UButton type="submit" color="success" :loading="settling">
                {{ t("accounting.recordPayment") }}
              </UButton>
            </div>
          </form>
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
            <div class="border-t pt-3">
              <div class="flex justify-between">
                <span class="text-gray-500">{{ t("accounting.selected") }}</span>
                <span class="font-semibold">${{ selectedTotal.toFixed(2) }}</span>
              </div>
            </div>
          </div>
        </UCard>

        <UButton
          v-if="selectedTotal > 0 && !showSettlementForm"
          color="success"
          block
          @click="openSettlementForm"
        >
          <UIcon name="i-lucide-credit-card" class="mr-2" />
          {{ t("accounting.processPayment") }}
        </UButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

const { t } = useI18n();
const route = useRoute();
const guestId = route.params.id as string;

const loading = ref(true);
const settling = ref(false);
const showSettlementForm = ref(false);

interface ReservationSettlement {
  id: number;
  reservationCode: string;
  checkInDate: string;
  checkOutDate: string;
  status: string;
  statusLabel: string;
  roomPrice: number;
  paidAmount: number;
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
  reservations: ReservationSettlement[];
  parkingTransactions: ParkingSettlement[];
  restaurantBills: RestaurantSettlement[];
  totalRoom: number;
  totalParking: number;
  totalRestaurant: number;
  totalDue: number;
  totalPaid: number;
  balance: number;
}

interface PaymentMethod {
  id: number;
  slug: string;
  label: string;
}

const settlement = ref<Settlement | null>(null);
const paymentMethods = ref<PaymentMethod[]>([]);

const selectedReservationIds = ref<number[]>([]);
const selectedParkingTxnIds = ref<number[]>([]);
const selectedRestaurantBillIds = ref<number[]>([]);
const selectAllRoom = ref(false);
const selectAllParking = ref(false);
const selectAllRestaurant = ref(false);

const toggleSelectAll = (val: boolean) => {
  if (val && settlement.value) {
    selectedReservationIds.value = settlement.value.reservations.map(r => r.id);
  } else {
    selectedReservationIds.value = [];
  }
};

const toggleSelectAllParking = (val: boolean) => {
  if (val && settlement.value) {
    selectedParkingTxnIds.value = settlement.value.parkingTransactions.map(p => p.id);
  } else {
    selectedParkingTxnIds.value = [];
  }
};

const toggleSelectAllRestaurant = (val: boolean) => {
  if (val && settlement.value) {
    selectedRestaurantBillIds.value = settlement.value.restaurantBills.map(b => b.id);
  } else {
    selectedRestaurantBillIds.value = [];
  }
};

const selectedTotal = computed(() => {
  let total = 0;
  if (settlement.value) {
    for (const r of settlement.value.reservations) {
      if (selectedReservationIds.value.includes(r.id)) {
        total += r.roomPrice;
      }
    }
    for (const p of settlement.value.parkingTransactions) {
      if (selectedParkingTxnIds.value.includes(p.id)) {
        total += p.amountDue;
      }
    }
    for (const b of settlement.value.restaurantBills) {
      if (selectedRestaurantBillIds.value.includes(b.id)) {
        total += b.totalAmount;
      }
    }
  }
  return total;
});

const form = reactive({
  amount: 0,
  paymentMethod: undefined as number | undefined,
  reference: "",
  notes: "",
});

const paymentMethodOptions = computed(() => {
  return paymentMethods.value.map(pm => ({
    value: pm.id,
    label: pm.label,
  }));
});

const getStatusColor = (status: string): "success" | "warning" | "info" | "error" | "neutral" => {
  const colors: Record<string, "success" | "warning" | "info" | "error" | "neutral"> = {
    pending: "warning",
    confirmed: "info",
    checked_in: "success",
    checked_out: "neutral",
    cancelled: "error",
    no_show: "error",
  };
  return colors[status] || "neutral";
};

const fetchPaymentMethods = async () => {
  try {
    const res = await $fetch("/api/accounting/payment-methods");
    paymentMethods.value = (res as any).data || [];
    if (paymentMethods.value.length > 0) {
      form.paymentMethod = paymentMethods.value[0].id;
    }
  } catch (e) {
    console.error("Failed to fetch payment methods:", e);
  }
};

const fetchSettlement = async () => {
  loading.value = true;
  try {
    const response: any = await $fetch(`/api/guests/${guestId}/settle`);
    settlement.value = response as Settlement;
  } catch (error) {
    console.error("Failed to fetch settlement:", error);
  } finally {
    loading.value = false;
  }
};

const openSettlementForm = () => {
  form.amount = selectedTotal.value;
  showSettlementForm.value = true;
};

const handleSettle = async () => {
  if (form.amount <= 0) return;
  if (!form.paymentMethod) return;

  settling.value = true;
  try {
    await $fetch(`/api/guests/${guestId}/settle`, {
      method: "POST",
      body: {
        reservationIds: selectedReservationIds.value,
        parkingTxnIds: selectedParkingTxnIds.value,
        restaurantBillIds: selectedRestaurantBillIds.value,
        amount: form.amount,
        paymentMethod: form.paymentMethod,
        reference: form.reference,
        notes: form.notes,
      },
    });

    showSettlementForm.value = false;
    selectedReservationIds.value = [];
    selectedParkingTxnIds.value = [];
    selectedRestaurantBillIds.value = [];
    selectAllRoom.value = false;
    selectAllParking.value = false;
    selectAllRestaurant.value = false;
    await fetchSettlement();
  } catch (error) {
    console.error("Failed to settle:", error);
  } finally {
    settling.value = false;
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
