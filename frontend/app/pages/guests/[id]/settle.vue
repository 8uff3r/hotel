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
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">{{ t("accounting.roomCharges") }}</h3>
          </template>
          <div v-if="settlement?.reservations?.length">
            <table class="w-full">
              <thead>
                <tr class="border-b text-left text-sm text-gray-500">
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
                  <td class="py-3">{{ res.reservationCode }}</td>
                  <td class="py-3">{{ formatDate(res.checkInDate) }}</td>
                  <td class="py-3">{{ formatDate(res.checkOutDate) }}</td>
                  <td class="py-3">
                    <UBadge :color="getStatusColor(res.status)" variant="soft" size="sm">
                      {{ res.status }}
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

        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">{{ t("accounting.parkingCharges") }}</h3>
          </template>
          <div v-if="settlement?.parkingTransactions?.length">
            <table class="w-full">
              <thead>
                <tr class="border-b text-left text-sm text-gray-500">
                  <th class="pb-2">License Plate</th>
                  <th class="pb-2">Entry</th>
                  <th class="pb-2">Exit</th>
                  <th class="pb-2 text-right">Hours</th>
                  <th class="pb-2 text-right">Rate</th>
                  <th class="pb-2 text-right">{{ t("accounting.amount") }}</th>
                  <th class="pb-2 text-right">{{ t("accounting.paid") }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="pt in settlement.parkingTransactions" :key="pt.id" class="border-b">
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
                {{ t("common.cancel") }}
              </UButton>
              <UButton type="submit" color="success" :loading="settling">
                {{ t("accounting.recordPayment") }}
              </UButton>
            </div>
          </form>
        </UCard>
      </div>

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

        <UButton
          v-if="settlement?.balance && settlement.balance > 0 && !showSettlementForm"
          color="success"
          block
          @click="showSettlementForm = true"
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

interface Settlement {
  reservations: ReservationSettlement[];
  parkingTransactions: ParkingSettlement[];
  totalRoom: number;
  totalParking: number;
  totalDue: number;
  totalPaid: number;
  balance: number;
}

const settlement = ref<Settlement | null>(null);

const form = reactive({
  amount: 0,
  paymentMethod: "cash",
  reference: "",
  notes: "",
});

const paymentMethodOptions = [
  { value: "cash", label: "Cash" },
  { value: "credit_card", label: "Credit Card" },
  { value: "debit_card", label: "Debit Card" },
  { value: "bank_transfer", label: "Bank Transfer" },
  { value: "other", label: "Other" },
];

const fetchSettlement = async () => {
  loading.value = true;
  try {
    const response = await $fetch(`/api/guests/${guestId}/settle`);
    settlement.value = response as Settlement;
    form.amount = settlement.value?.balance || 0;
  } catch (error) {
    console.error("Failed to fetch settlement:", error);
  } finally {
    loading.value = false;
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

const getStatusColor = (status: string): "success" | "warning" | "info" | "error" | "neutral" => {
  const colors: Record<string, "success" | "warning" | "info" | "error" | "neutral"> = {
    pending: "warning",
    confirmed: "info",
    checked_in: "success",
    checked_out: "neutral",
    cancelled: "error",
  };
  return colors[status] || "neutral";
};

const handleSettle = async () => {
  if (form.amount <= 0) return;

  settling.value = true;
  try {
    await $fetch(`/api/guests/${guestId}/settle`, {
      method: "POST",
      body: {
        amount: form.amount,
        paymentMethod: form.paymentMethod,
        reference: form.reference,
        notes: form.notes,
      },
    });

    showSettlementForm.value = false;
    await fetchSettlement();
  } catch (error) {
    console.error("Failed to settle:", error);
  } finally {
    settling.value = false;
  }
};

onMounted(fetchSettlement);
</script>
