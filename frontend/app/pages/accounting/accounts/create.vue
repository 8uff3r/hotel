<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/accounting/accounts" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        {{ t("accounting.back_to_accounts") }}
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("accounting.add_new_account") }}
      </h1>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <!-- Account Code -->
          <UFormField :label="t('accounting.account_code')" name="accountCode" required>
            <UInput
              v-model="form.accountCode"
              :placeholder="t('accounting.e_g_1000_2000_4000')"
              :disabled="loading"
            />
          </UFormField>

          <!-- Account Name -->
          <UFormField :label="t('accounting.account_name')" name="accountName" required>
            <UInput
              v-model="form.accountName"
              :placeholder="t('accounting.e_g_cash_accounts_receivable')"
              :disabled="loading"
            />
          </UFormField>

          <!-- Account Type -->
          <UFormField :label="t('accounting.account_type')" name="accountType" required>
            <USelect v-model="form.accountType" :items="accountTypeOptions" :disabled="loading" />
          </UFormField>

          <!-- Normal Balance -->
          <UFormField :label="t('accounting.normal_balance')" name="normalBalance" required>
            <USelect
              v-model="form.normalBalance"
              :items="normalBalanceOptions"
              :disabled="loading"
            />
          </UFormField>

          <!-- Description -->
          <UFormField :label="t('common.description')" name="description" class="md:col-span-2">
            <UTextarea
              v-model="form.description"
              :placeholder="t('accounting.account_description')"
              :rows="3"
              :disabled="loading"
            />
          </UFormField>

          <!-- Is Active -->
          <UFormField :label="t('common.status')" name="isActive">
            <div class="flex items-center gap-2">
              <UCheckbox v-model="form.isActive" :disabled="loading" />
              <span class="text-sm text-gray-600">{{ t("accounting.active_account") }}</span>
            </div>
          </UFormField>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/accounting/accounts" :disabled="loading">{{
            t("common.cancel")
          }}</UButton>
          <UButton type="submit" color="primary" :loading="loading">
            {{ t("accounting.create_account") }}
          </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager"],
});

const { t } = useI18n();
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
