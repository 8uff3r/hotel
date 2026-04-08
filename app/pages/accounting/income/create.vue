<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/accounting/income" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        Back to Income
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Record Income</h1>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <!-- Income Date -->
          <UFormField label="Income Date" name="incomeDate" required>
            <UInput v-model="form.incomeDate" type="date" :disabled="loading" />
          </UFormField>

          <!-- Amount -->
          <UFormField label="Amount ($)" name="amount" required>
            <UInput
              v-model.number="form.amount"
              type="number"
              min="0"
              step="0.01"
              :disabled="loading"
            />
          </UFormField>

          <!-- Category -->
          <UFormField label="Category" name="category" required>
            <USelect v-model="form.category" :items="categoryOptions" :disabled="loading" />
          </UFormField>

          <!-- Payment Status -->
          <UFormField label="Payment Status" name="paymentStatus">
            <USelect
              v-model="form.paymentStatus"
              :items="paymentStatusOptions"
              :disabled="loading"
            />
          </UFormField>

          <!-- Description -->
          <UFormField label="Description" name="description" required class="md:col-span-2">
            <UInput
              v-model="form.description"
              placeholder="Income description"
              :disabled="loading"
            />
          </UFormField>

          <!-- Source -->
          <UFormField label="Source" name="source">
            <UInput
              v-model="form.source"
              placeholder="e.g., Guest Name, Company"
              :disabled="loading"
            />
          </UFormField>

          <!-- Reference -->
          <UFormField label="Reference" name="reference">
            <UInput
              v-model="form.reference"
              placeholder="Invoice/Receipt number"
              :disabled="loading"
            />
          </UFormField>

          <!-- Payment Method -->
          <UFormField label="Payment Method" name="paymentMethod">
            <USelect
              v-model="form.paymentMethod"
              :items="paymentMethodOptions"
              :disabled="loading"
            />
          </UFormField>

          <!-- Reservation Link -->
          <UFormField label="Reservation ID (Optional)" name="reservationId">
            <UInput
              v-model.number="form.reservationId"
              type="number"
              placeholder="Link to reservation"
              :disabled="loading"
            />
          </UFormField>

          <!-- Notes -->
          <UFormField label="Notes" name="notes" class="md:col-span-2">
            <UTextarea
              v-model="form.notes"
              placeholder="Additional notes..."
              :rows="3"
              :disabled="loading"
            />
          </UFormField>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/accounting/income" :disabled="loading"> Cancel </UButton>
          <UButton type="submit" color="success" :loading="loading"> Record Income </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager"],
});

const loading = ref(false);

const form = reactive({
  incomeDate: new Date().toISOString().split("T")[0],
  description: "",
  amount: 0,
  category: "room_revenue" as string,
  source: "",
  reference: "",
  paymentMethod: "cash" as string,
  paymentStatus: "received" as string,
  reservationId: undefined as number | undefined,
  notes: "",
});

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
    await $fetch("/api/income", {
      method: "POST",
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
