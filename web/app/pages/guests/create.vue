<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/guests" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        {{ t("actions.backToGuests") }}
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("guests.addNew") }}</h1>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <!-- First Name -->
          <UFormField :label="t('forms.firstName')" name="firstName" required>
            <UInput v-model="form.firstName" :placeholder="t('forms.firstNamePlaceholder')" :disabled="loading" />
          </UFormField>

          <!-- Last Name -->
          <UFormField :label="t('forms.lastName')" name="lastName" required>
            <UInput v-model="form.lastName" :placeholder="t('forms.lastNamePlaceholder')" :disabled="loading" />
          </UFormField>

          <!-- Email -->
          <UFormField :label="t('forms.email')" name="email">
            <UInput
              v-model="form.email"
              type="email"
              :placeholder="t('guests.emailPlaceholder')"
              :disabled="loading"
            />
          </UFormField>

          <!-- Phone -->
          <UFormField :label="t('forms.phone')" name="phone">
            <UInput v-model="form.phone" :placeholder="t('guests.phonePlaceholder')" :disabled="loading" />
          </UFormField>

          <!-- ID Type -->
          <UFormField :label="t('guests.idType')" name="idType">
            <USelect
              v-model="form.idType"
              :items="idTypeOptions"
              :placeholder="t('guests.selectIdType')"
              :disabled="loading"
            />
          </UFormField>

          <!-- ID Number -->
          <UFormField :label="t('guests.idNumber')" name="idNumber">
            <UInput v-model="form.idNumber" :placeholder="t('guests.idNumberPlaceholder')" :disabled="loading" />
          </UFormField>

          <!-- Address -->
          <UFormField :label="t('forms.address')" name="address" class="md:col-span-2">
            <UInput v-model="form.address" :placeholder="t('forms.addressPlaceholder')" :disabled="loading" />
          </UFormField>

          <!-- City -->
          <UFormField :label="t('forms.city')" name="city">
            <UInput v-model="form.city" :placeholder="t('forms.cityPlaceholder')" :disabled="loading" />
          </UFormField>

          <!-- Country -->
          <UFormField :label="t('forms.country')" name="country">
            <UInput v-model="form.country" :placeholder="t('forms.countryPlaceholder')" :disabled="loading" />
          </UFormField>

          <!-- Notes -->
          <UFormField :label="t('forms.notes')" name="notes" class="md:col-span-2">
            <UTextarea
              v-model="form.notes"
              :placeholder="t('forms.notesPlaceholder')"
              :rows="3"
              :disabled="loading"
            />
          </UFormField>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/guests" :disabled="loading"> {{ t("actions.cancel") }} </UButton>
          <UButton type="submit" color="primary" :loading="loading"> {{ t("guests.createGuest") }} </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});
const { t } = useI18n();

const loading = ref(false);

const form = reactive({
  firstName: "",
  lastName: "",
  email: "",
  phone: "",
  idType: undefined as string | undefined,
  idNumber: "",
  address: "",
  city: "",
  country: "",
  notes: "",
});

const submitBody = computed(() => ({
  firstName: form.firstName,
  lastName: form.lastName,
  email: form.email || undefined,
  phone: form.phone || undefined,
  idType: form.idType || undefined,
  idNumber: form.idNumber || undefined,
  address: form.address || undefined,
  city: form.city || undefined,
  country: form.country || undefined,
  notes: form.notes || undefined,
}));

const idTypeOptions = computed(() => [
  { value: "passport", label: t("idTypes.passport") },
  { value: "national_id", label: t("idTypes.nationalId") },
  { value: "driver_license", label: t("idTypes.driverLicense") },
  { value: "other", label: t("idTypes.other") },
]);

const handleSubmit = async () => {
  loading.value = true;
  try {
    await $fetch("/api/guests", {
      method: "POST",
      body: submitBody.value,
    });

    await navigateTo("/guests");
  } catch (error) {
    console.error("Failed to create guest:", error);
  } finally {
    loading.value = false;
  }
};
</script>
