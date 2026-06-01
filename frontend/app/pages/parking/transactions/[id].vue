<template>
  <div>
    <div class="mb-6">
      <UButton to="/parking/transactions" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />{{
          t("parking.back_to_transactions")
        }}</UButton
      >
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("parking.transaction_details") }}
      </h1>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <UIcon name="i-lucide-loader-2" class="h-8 w-8 animate-spin" />
    </div>

    <div v-else-if="transaction">
      <div class="mb-6 grid grid-cols-1 gap-4 md:grid-cols-3">
        <UCard>
          <template #header>
            <span class="text-sm text-gray-500">{{ t("parking.license_plate") }}</span>
          </template>
          <div class="text-2xl font-bold">{{ transaction.licensePlate }}</div>
        </UCard>

        <UCard>
          <template #header>
            <span class="text-sm text-gray-500">{{ t("common.status") }}</span>
          </template>
          <UBadge :color="getStatusColor(transaction.status) as any" size="lg">
            {{ transaction.status }}
          </UBadge>
        </UCard>

        <UCard>
          <template #header>
            <span class="text-sm text-gray-500">{{ t("parking.payment") }}</span>
          </template>
          <UBadge :color="getPaymentColor(transaction.paymentStatus) as any" size="lg">
            {{ transaction.paymentStatus }}
          </UBadge>
        </UCard>
      </div>

      <UCard>
        <template #header>
          <span class="font-semibold">{{ t("parking.transaction_information") }}</span>
        </template>

        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <div>
            <div class="mb-1 text-sm text-gray-500">{{ t("parking.entry_time") }}</div>
            <div class="font-medium">{{ formatDate(transaction.entryTime) }}</div>
          </div>

          <div>
            <div class="mb-1 text-sm text-gray-500">{{ t("parking.exit_time") }}</div>
            <div class="font-medium">
              {{ transaction.exitTime ? formatDate(transaction.exitTime) : "-" }}
            </div>
          </div>

          <div>
            <div class="mb-1 text-sm text-gray-500">{{ t("parking.duration") }}</div>
            <div class="font-medium">
              {{
                transaction.hoursParked
                  ? t("parking.hours_long", { hours: transaction.hoursParked.toFixed(1) })
                  : "-"
              }}
            </div>
          </div>

          <div>
            <div class="mb-1 text-sm text-gray-500">{{ t("parking.rate_applied") }}</div>
            <div class="font-medium">{{ transaction.rateApplied || t("parking.hourly") }}</div>
          </div>

          <div>
            <div class="mb-1 text-sm text-gray-500">{{ t("parking.amount_due") }}</div>
            <div class="text-xl font-medium">
              ${{ transaction.amountDue?.toFixed(2) || "0.00" }}
            </div>
          </div>

          <div>
            <div class="mb-1 text-sm text-gray-500">{{ t("parking.amount_paid") }}</div>
            <div class="text-xl font-medium">
              ${{ transaction.amountPaid?.toFixed(2) || "0.00" }}
            </div>
          </div>

          <div v-if="transaction.paymentMethod">
            <div class="mb-1 text-sm text-gray-500">{{ t("parking.payment_method") }}</div>
            <div class="font-medium">{{ transaction.paymentMethod }}</div>
          </div>

          <div v-if="transaction.notes">
            <div class="mb-1 text-sm text-gray-500">{{ t("common.notes") }}</div>
            <div class="font-medium">{{ transaction.notes }}</div>
          </div>
        </div>

        <div v-if="transaction.status === 'active'" class="mt-6 flex justify-end">
          <UButton color="warning" @click="openCheckout">
            <UIcon name="i-lucide-log-out" class="mr-2" />
            {{ t("parking.check_out_now") }}
          </UButton>
        </div>
      </UCard>
    </div>

    <UModal v-model="checkoutModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("parking.check_out_vehicle") }}</h2>
      </template>
      <template #body>
        <div class="space-y-4">
          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.payment_method") }}</label>
            <USelect v-model="checkoutForm.paymentMethod" :items="paymentMethods" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("parking.amount_paid") }}</label>
            <UInput v-model="checkoutForm.amountPaid" type="number" step="0.01" />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("common.notes") }}</label>
            <UTextarea v-model="checkoutForm.notes" :rows="2" />
          </div>
        </div>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="checkoutModalOpen = false">{{
            t("common.cancel")
          }}</UButton>
          <UButton color="warning" :loading="checkingOut" @click="confirmCheckout">{{
            t("parking.check_out")
          }}</UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n();

const route = useRoute();
const transactionId = Number(route.params.id);

const transaction = ref<any>(null);
const loading = ref(true);
const checkoutModalOpen = ref(false);
const checkingOut = ref(false);

const checkoutForm = reactive({
  paymentMethod: "cash",
  amountPaid: "",
  notes: "",
});

const paymentMethods = [
  { value: "cash", label: t("parking.payment_cash") },
  { value: "card", label: t("parking.payment_card") },
  { value: "guest_account", label: t("parking.payment_guest_account") },
];

const fetchTransaction = async () => {
  try {
    transaction.value = await $fetch(`/api/parking/transactions/${transactionId}`);
    checkoutForm.amountPaid = transaction.value.amountDue?.toString() || "0";
  } catch (error) {
    console.error("Failed to fetch transaction:", error);
  } finally {
    loading.value = false;
  }
};

const openCheckout = () => {
  checkoutForm.amountPaid = transaction.value.amountDue?.toString() || "0";
  checkoutForm.paymentMethod = "cash";
  checkoutForm.notes = "";
  checkoutModalOpen.value = true;
};

const confirmCheckout = async () => {
  checkingOut.value = true;
  try {
    await $fetch(`/api/parking/transactions/${transactionId}/check-out`, {
      method: "POST",
      body: {
        rateType: "hourly",
        amountPaid: parseFloat(checkoutForm.amountPaid) || 0,
        paymentMethod: checkoutForm.paymentMethod,
        notes: checkoutForm.notes || null,
      },
    });
    checkoutModalOpen.value = false;
    await fetchTransaction();
  } catch (error) {
    console.error("Failed to check out:", error);
  } finally {
    checkingOut.value = false;
  }
};

const formatDate = (date: string) => {
  return new Date(date).toLocaleString();
};

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    active: "info",
    completed: "success",
    cancelled: "error",
  };
  return colors[status] || "neutral";
};

const getPaymentColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: "warning",
    paid: "success",
    waived: "info",
  };
  return colors[status] || "neutral";
};

onMounted(fetchTransaction);
</script>
