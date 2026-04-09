<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/accounting/accounts" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        Back to Accounts
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Add New Account</h1>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <!-- Account Code -->
          <UFormField label="Account Code" name="accountCode" required>
            <UInput
              v-model="form.accountCode"
              placeholder="e.g., 1000, 2000, 4000"
              :disabled="loading"
            />
          </UFormField>

          <!-- Account Name -->
          <UFormField label="Account Name" name="accountName" required>
            <UInput
              v-model="form.accountName"
              placeholder="e.g., Cash, Accounts Receivable"
              :disabled="loading"
            />
          </UFormField>

          <!-- Account Type -->
          <UFormField label="Account Type" name="accountType" required>
            <USelect v-model="form.accountType" :items="accountTypeOptions" :disabled="loading" />
          </UFormField>

          <!-- Normal Balance -->
          <UFormField label="Normal Balance" name="normalBalance" required>
            <USelect
              v-model="form.normalBalance"
              :items="normalBalanceOptions"
              :disabled="loading"
            />
          </UFormField>

          <!-- Description -->
          <UFormField label="Description" name="description" class="md:col-span-2">
            <UTextarea
              v-model="form.description"
              placeholder="Account description..."
              :rows="3"
              :disabled="loading"
            />
          </UFormField>

          <!-- Is Active -->
          <UFormField label="Status" name="isActive">
            <div class="flex items-center gap-2">
              <UCheckbox v-model="form.isActive" :disabled="loading" />
              <span class="text-sm text-gray-600">Active Account</span>
            </div>
          </UFormField>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/accounting/accounts" :disabled="loading">
            Cancel
          </UButton>
          <UButton type="submit" color="primary" :loading="loading"> Create Account </UButton>
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
  accountCode: "",
  accountName: "",
  accountType: "expense" as string,
  normalBalance: "debit" as string,
  description: "",
  isActive: true,
});

const accountTypeOptions = [
  { value: "asset", label: "Asset" },
  { value: "liability", label: "Liability" },
  { value: "equity", label: "Equity" },
  { value: "revenue", label: "Revenue" },
  { value: "expense", label: "Expense" },
];

const normalBalanceOptions = [
  { value: "debit", label: "Debit" },
  { value: "credit", label: "Credit" },
];

const handleSubmit = async () => {
  if (!form.accountCode || !form.accountName) {
    return;
  }

  loading.value = true;
  try {
    await $fetch("/api/accounts", {
      method: "POST",
      body: {
        accountCode: form.accountCode,
        accountName: form.accountName,
        accountType: form.accountType,
        normalBalance: form.normalBalance,
        description: form.description || undefined,
        isActive: form.isActive,
      },
    });

    await navigateTo("/accounting/accounts");
  } catch (error) {
    console.error("Failed to create account:", error);
  } finally {
    loading.value = false;
  }
};
</script>
