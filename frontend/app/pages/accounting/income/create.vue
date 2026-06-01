<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/accounting/income" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        {{ t("accounting.back_to_income") }}
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("accounting.record_income") }}
      </h1>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <!-- Income Date -->
          <UFormField :label="t('accounting.income_date')" name="incomeDate" required>
            <UInput v-model="form.incomeDate" type="date" :disabled="loading" />
          </UFormField>

          <!-- Amount -->
          <UFormField :label="t('accounting.amount')" name="amount" required>
            <UInput
              v-model.number="form.amount"
              type="number"
              min="0"
              step="0.01"
              :disabled="loading"
            />
          </UFormField>

          <!-- Category -->
          <UFormField :label="t('accounting.category')" name="category" required>
            <USelect v-model="form.category" :items="categoryOptions" :disabled="loading" />
          </UFormField>

          <!-- Payment Status -->
          <UFormField :label="t('common.payment_status')" name="paymentStatus">
            <USelect
              v-model="form.paymentStatus"
              :items="paymentStatusOptions"
              :disabled="loading"
            />
          </UFormField>

          <!-- Description -->
          <UFormField
            :label="t('common.description')"
            name="description"
            required
            class="md:col-span-2"
          >
            <UInput
              v-model="form.description"
              :placeholder="t('accounting.income_description')"
              :disabled="loading"
            />
          </UFormField>

          <!-- Source -->
          <UFormField :label="t('accounting.source')" name="source">
            <UInput
              v-model="form.source"
              :placeholder="t('accounting.e_g_guest_name_company')"
              :disabled="loading"
            />
          </UFormField>

          <!-- Reference -->
          <UFormField :label="t('accounting.reference')" name="reference">
            <UInput
              v-model="form.reference"
              :placeholder="t('accounting.invoice_receipt_number')"
              :disabled="loading"
            />
          </UFormField>

          <!-- Payment Method -->
          <UFormField :label="t('accounting.payment_method')" name="paymentMethod">
            <USelect
              v-model="form.paymentMethod"
              :items="paymentMethodOptions"
              :disabled="loading"
            />
          </UFormField>

          <!-- Reservation Link -->
          <UFormField :label="t('accounting.reservation_id_optional')" name="reservationId">
            <UInput
              v-model.number="form.reservationId"
              type="number"
              :placeholder="t('accounting.link_to_reservation')"
              :disabled="loading"
            />
          </UFormField>

          <!-- Notes -->
          <UFormField :label="t('common.notes')" name="notes" class="md:col-span-2">
            <UTextarea
              v-model="form.notes"
              :placeholder="t('accounting.additional_notes')"
              :rows="3"
              :disabled="loading"
            />
          </UFormField>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/accounting/income" :disabled="loading">
            {{ t("common.cancel") }}
          </UButton>
          <UButton type="submit" color="success" :loading="loading">
            {{ t("accounting.record_income") }}
          </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { postApiAccountingIncome } from "~/utils/client";

const categoryOptions = [
  { value: "room_revenue", label: "Room Revenue" },
  { value: "food_beverage", label: "Food & Beverage" },
  { value: "laundry", label: "Laundry" },
  { value: "spa", label: "Spa" },
  { value: "meeting_rooms", label: "Meeting Rooms" },
  { value: "other", label: "Other" },
];

const paymentStatusOptions = [
  { value: "pending", label: "Pending" },
  { value: "received", label: "Received" },
  { value: "refunded", label: "Refunded" },
];

const paymentMethodOptions = [
  { value: "cash", label: "Cash" },
  { value: "credit_card", label: "Credit Card" },
  { value: "debit_card", label: "Debit Card" },
  { value: "bank_transfer", label: "Bank Transfer" },
  { value: "other", label: "Other" },
];

const handleSubmit = async () => {
  if (!form.description || form.amount <= 0) {
    return;
  }

  loading.value = true;
  try {
    await postApiAccountingIncome({
      body: {
        incomeDate: form.incomeDate,
        description: form.description,
        amount: form.amount,
        category: form.category,
        source: form.source || undefined,
        reference: form.reference || undefined,
        paymentMethod: form.paymentMethod,
        paymentStatus: form.paymentStatus,
        reservationId: form.reservationId,
        notes: form.notes || undefined,
      },
    });

    await navigateTo("/accounting/income");
  } catch (error) {
    console.error("Failed to create income:", error);
  } finally {
    loading.value = false;
  }
};
</script>
