<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/accounting/expenses" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        Back to Expenses
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Record Expense</h1>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <!-- Expense Date -->
          <UFormField label="Expense Date" name="expenseDate" required>
            <UInput v-model="form.expenseDate" type="date" :disabled="loading" />
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
              placeholder="Expense description"
              :disabled="loading"
            />
          </UFormField>

          <!-- Vendor -->
          <UFormField label="Vendor" name="vendor">
            <UInput v-model="form.vendor" placeholder="Vendor/Supplier name" :disabled="loading" />
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

          <!-- Receipt Number -->
          <UFormField label="Receipt Number" name="receiptNumber">
            <UInput v-model="form.receiptNumber" placeholder="Receipt number" :disabled="loading" />
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
          <UButton variant="outline" to="/accounting/expenses" :disabled="loading">
            Cancel
          </UButton>
          <UButton type="submit" color="error" :loading="loading"> Record Expense </UButton>
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
  expenseDate: new Date().toISOString().split("T")[0],
  description: "",
  amount: 0,
  category: "other" as string,
  vendor: "",
  reference: "",
  paymentMethod: "cash" as string,
  paymentStatus: "pending" as string,
  receiptNumber: "",
  notes: "",
});

const categoryOptions = [
  { value: "food_beverage", label: "Food & Beverage" },
  { value: "housekeeping", label: "Housekeeping" },
  { value: "maintenance", label: "Maintenance" },
  { value: "utilities", label: "Utilities" },
  { value: "salaries", label: "Salaries" },
  { value: "marketing", label: "Marketing" },
  { value: "supplies", label: "Supplies" },
  { value: "insurance", label: "Insurance" },
  { value: "taxes", label: "Taxes" },
  { value: "rent", label: "Rent" },
  { value: "other", label: "Other" },
];

const paymentStatusOptions = [
  { value: "pending", label: "Pending" },
  { value: "paid", label: "Paid" },
  { value: "cancelled", label: "Cancelled" },
];

const paymentMethodOptions = [
  { value: "cash", label: "Cash" },
  { value: "credit_card", label: "Credit Card" },
  { value: "debit_card", label: "Debit Card" },
  { value: "bank_transfer", label: "Bank Transfer" },
  { value: "check", label: "Check" },
  { value: "other", label: "Other" },
];

const handleSubmit = async () => {
  if (!form.description || form.amount <= 0) {
    return;
  }

  loading.value = true;
  try {
    await $fetch("/api/expenses", {
      method: "POST",
      body: {
        expenseDate: form.expenseDate,
        description: form.description,
        amount: form.amount,
        category: form.category,
        vendor: form.vendor || undefined,
        reference: form.reference || undefined,
        paymentMethod: form.paymentMethod,
        paymentStatus: form.paymentStatus,
        receiptNumber: form.receiptNumber || undefined,
        notes: form.notes || undefined,
      },
    });

    await navigateTo("/accounting/expenses");
  } catch (error) {
    console.error("Failed to create expense:", error);
  } finally {
    loading.value = false;
  }
};
</script>
