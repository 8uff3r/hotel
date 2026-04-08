<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/guests" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        Back to Guests
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Add New Guest</h1>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <!-- First Name -->
          <UFormField label="First Name" name="firstName" required>
            <UInput v-model="form.firstName" placeholder="John" :disabled="loading" />
          </UFormField>

          <!-- Last Name -->
          <UFormField label="Last Name" name="lastName" required>
            <UInput v-model="form.lastName" placeholder="Doe" :disabled="loading" />
          </UFormField>

          <!-- Email -->
          <UFormField label="Email" name="email">
            <UInput
              v-model="form.email"
              type="email"
              placeholder="john@example.com"
              :disabled="loading"
            />
          </UFormField>

          <!-- Phone -->
          <UFormField label="Phone" name="phone">
            <UInput v-model="form.phone" placeholder="+1 234 567 8900" :disabled="loading" />
          </UFormField>

          <!-- ID Type -->
          <UFormField label="ID Type" name="idType">
            <USelect
              v-model="form.idType"
              :items="idTypeOptions"
              placeholder="Select ID type"
              :disabled="loading"
            />
          </UFormField>

          <!-- ID Number -->
          <UFormField label="ID Number" name="idNumber">
            <UInput v-model="form.idNumber" placeholder="ID number" :disabled="loading" />
          </UFormField>

          <!-- Address -->
          <UFormField label="Address" name="address" class="md:col-span-2">
            <UInput v-model="form.address" placeholder="Street address" :disabled="loading" />
          </UFormField>

          <!-- City -->
          <UFormField label="City" name="city">
            <UInput v-model="form.city" placeholder="City" :disabled="loading" />
          </UFormField>

          <!-- Country -->
          <UFormField label="Country" name="country">
            <UInput v-model="form.country" placeholder="Country" :disabled="loading" />
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
          <UButton variant="outline" to="/guests" :disabled="loading"> Cancel </UButton>
          <UButton type="submit" color="primary" :loading="loading"> Create Guest </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

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

const idTypeOptions = [
  { value: "passport", label: "Passport" },
  { value: "national_id", label: "National ID" },
  { value: "driver_license", label: "Driver License" },
  { value: "other", label: "Other" },
];

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
