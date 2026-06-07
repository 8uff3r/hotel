<template>
  <div class="p-6" v-if="stay">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">جزئیات پذیرش {{ stay.acceptanceId }}</h1>
      <div class="flex gap-2">
        <UButton v-if="stay.status?.slug === 'waiting'" @click="checkIn" :loading="actionLoading">
          ورود
        </UButton>
        <UButton v-if="stay.status?.slug === 'resident'" @click="checkOut" :loading="actionLoading" color="red">
          خروج
        </UButton>
        <UButton to="/stays">بازگشت</UButton>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-6">
      <UCard>
        <template #header>
          <h2 class="font-semibold">اطلاعات مهمان</h2>
        </template>
        <div class="space-y-2">
          <p><strong>نام:</strong> {{ stay.guest?.firstName }} {{ stay.guest?.lastName }}</p>
          <p><strong>کد ملی:</strong> {{ stay.guest?.nationalId }}</p>
          <p><strong>تلفن:</strong> {{ stay.guest?.phone }}</p>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <h2 class="font-semibold">اطلاعات اتاق</h2>
        </template>
        <div class="space-y-2">
          <p><strong>شماره اتاق:</strong> {{ stay.room?.roomNumber }}</p>
          <p><strong>ظرفیت:</strong> {{ stay.room?.capacity }}</p>
          <p><strong>نوع:</strong> {{ stay.room?.type?.label }}</p>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <h2 class="font-semibold">اطلاعات اقامت</h2>
        </template>
        <div class="space-y-2">
          <p><strong>تاریخ ورود:</strong> {{ formatDate(stay.entryDate) }}</p>
          <p><strong>تاریخ خروج:</strong> {{ formatDate(stay.departureDate) }}</p>
          <p><strong>تعداد نفرات:</strong> {{ stay.numberOfPeople }}</p>
          <p><strong>قیمت اتاق:</strong> {{ stay.roomPrice }}</p>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <h2 class="font-semibold">صورتحساب</h2>
        </template>
        <div v-if="invoice" class="space-y-2">
          <p><strong>مبلغ کل:</strong> {{ invoice.totalAmount }}</p>
          <p><strong>پرداخت شده:</strong> {{ invoice.paidAmount }}</p>
          <p><strong>مانده:</strong> {{ invoice.remainingAmount }}</p>
          <p><strong>وضعیت:</strong> {{ invoice.paymentStatus }}</p>
          <UButton size="xs" @click="showPayment = true">پرداخت</UButton>
        </div>
        <div v-else>صورتحساب یافت نشد</div>
      </UCard>
    </div>

    <!-- Payment Modal -->
    <UModal v-model="showPayment">
      <UCard>
        <template #header>
          <h3 class="font-semibold">پرداخت صورتحساب</h3>
        </template>
        <UForm :state="paymentState" @submit="makePayment">
          <UFormGroup label="مبلغ" name="amount">
            <UInput v-model.number="paymentState.amount" type="number" />
          </UFormGroup>
          <UFormGroup label="روش پرداخت" name="paymentMethod" class="mt-2">
            <USelect v-model="paymentState.paymentMethod" :options="paymentMethodOptions" />
          </UFormGroup>
          <div class="mt-4 flex gap-2">
            <UButton type="submit" :loading="paymentLoading">پرداخت</UButton>
            <UButton variant="outline" @click="showPayment = false">انصراف</UButton>
          </div>
        </UForm>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
const route = useRoute();
const stayId = route.params.id as string;

const { data: stay, refresh: refreshStay } = useFetch(`/api/stays/${stayId}`, {
  key: `stay-${stayId}`,
});

const { data: invoice, refresh: refreshInvoice } = useFetch(`/api/stays/${stayId}/invoice`, {
  key: `stay-invoice-${stayId}`,
});

const actionLoading = ref(false);
const showPayment = ref(false);
const paymentLoading = ref(false);
const paymentState = reactive({ amount: 0, paymentMethod: undefined as number | undefined });

const paymentMethodOptions = ref([
  { label: "نقدی", value: 1 },
  { label: "کارت", value: 2 },
  { label: "انتقال بانکی", value: 3 },
]);

function formatDate(date: string) {
  if (!date) return "-";
  return new Date(date).toLocaleDateString("fa-IR");
}

async function checkIn() {
  actionLoading.value = true;
  try {
    await $fetch(`/api/stays/${stayId}/check-in`, { method: "POST" });
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
    await $fetch(`/api/stays/${stayId}/check-out`, { method: "POST" });
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
    await $fetch(`/api/stays/${stayId}/invoice/pay`, {
      method: "POST",
      body: paymentState,
    });
    showPayment.value = false;
    refreshInvoice();
  } catch (e) {
    console.error(e);
  } finally {
    paymentLoading.value = false;
  }
}
</script>
