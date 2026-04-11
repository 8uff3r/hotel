hotel/app/pages/guests/[id].vue ``` ```vue
<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/guests" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        {{ t("actions.backToGuests") }}
      </UButton>
      <div class="flex items-center justify-between">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          {{ t("guests.guestNumber", { id: guest?.id }) }}
        </h1>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <ULoader size="lg" />
    </div>

    <div v-else-if="guest">
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- Main Form -->
        <div class="lg:col-span-2">
          <UCard>
            <template #header>
              <h3 class="text-lg font-semibold">{{ t("guests.information") }}</h3>
            </template>
            <form @submit.prevent="handleSubmit">
              <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
                <!-- First Name -->
                <UFormField :label="t('forms.firstName')" name="firstName" required>
                  <UInput v-model="form.firstName" :disabled="loading || !editing" />
                </UFormField>

                <!-- Last Name -->
                <UFormField :label="t('forms.lastName')" name="lastName" required>
                  <UInput v-model="form.lastName" :disabled="loading || !editing" />
                </UFormField>

                <!-- Email -->
                <UFormField :label="t('forms.email')" name="email">
                  <UInput v-model="form.email" type="email" :disabled="loading || !editing" />
                </UFormField>

                <!-- Phone -->
                <UFormField :label="t('forms.phone')" name="phone">
                  <UInput v-model="form.phone" :disabled="loading || !editing" />
                </UFormField>

                <!-- ID Type -->
                <UFormField :label="t('guests.idType')" name="idType">
                  <USelect
                    v-model="form.idType"
                    :items="idTypeOptions"
                    :disabled="loading || !editing"
                  />
                </UFormField>

                <!-- ID Number -->
                <UFormField :label="t('guests.idNumber')" name="idNumber">
                  <UInput v-model="form.idNumber" :disabled="loading || !editing" />
                </UFormField>

                <!-- Address -->
                <UFormField :label="t('forms.address')" name="address" class="md:col-span-2">
                  <UInput v-model="form.address" :disabled="loading || !editing" />
                </UFormField>

                <!-- City -->
                <UFormField :label="t('forms.city')" name="city">
                  <UInput v-model="form.city" :disabled="loading || !editing" />
                </UFormField>

                <!-- Country -->
                <UFormField :label="t('forms.country')" name="country">
                  <UInput v-model="form.country" :disabled="loading || !editing" />
                </UFormField>

                <!-- Notes -->
                <UFormField :label="t('forms.notes')" name="notes" class="md:col-span-2">
                  <UTextarea v-model="form.notes" :rows="3" :disabled="loading || !editing" />
                </UFormField>
              </div>

              <div class="mt-6 flex justify-end gap-3">
                <UButton v-if="!editing" variant="outline" @click="editing = true">
                  {{ t("actions.editGuest") }}
                </UButton>
                <template v-else>
                  <UButton variant="outline" :disabled="loading" @click="cancelEdit">
                    {{ t("actions.cancel") }}
                  </UButton>
                  <UButton type="submit" color="primary" :loading="loading">
                    {{ t("actions.saveChanges") }}
                  </UButton>
                </template>
              </div>
            </form>
          </UCard>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Quick Info -->
          <UCard>
            <template #header>
              <h3 class="text-lg font-semibold">{{ t("guests.quickInfo") }}</h3>
            </template>
            <div class="space-y-3">
              <div>
                <p class="text-sm text-gray-500">{{ t("guests.fullName") }}</p>
                <p class="font-medium">{{ guest.firstName }} {{ guest.lastName }}</p>
              </div>
              <div v-if="guest.email">
                <p class="text-sm text-gray-500">{{ t("forms.email") }}</p>
                <p class="font-medium">{{ guest.email }}</p>
              </div>
              <div v-if="guest.phone">
                <p class="text-sm text-gray-500">{{ t("forms.phone") }}</p>
                <p class="font-medium">{{ guest.phone }}</p>
              </div>
              <div v-if="guest.idType">
                <p class="text-sm text-gray-500">{{ t("guests.id") }}</p>
                <p class="font-medium">
                  {{ formatIdType(guest.idType) }}
                  <span v-if="guest.idNumber">- {{ guest.idNumber }}</span>
                </p>
              </div>
            </div>
          </UCard>

          <!-- Meta Info -->
          <UCard>
            <template #header>
              <h3 class="text-lg font-semibold">{{ t("guests.metaInformation") }}</h3>
            </template>
            <div class="space-y-2 text-sm">
              <p class="text-gray-500">{{ t("common.created", { value: formatDate(guest.createdAt) }) }}</p>
              <p class="text-gray-500">{{ t("common.updated", { value: formatDate(guest.updatedAt) }) }}</p>
            </div>
          </UCard>
        </div>
      </div>
    </div>

    <!-- Not Found -->
    <UCard v-else>
      <div class="py-12 text-center">
        <p class="text-gray-500">{{ t("guests.notFound") }}</p>
        <UButton to="/guests" class="mt-4">{{ t("actions.backToGuests") }}</UButton>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});
const { t } = useI18n();

interface Guest {
  id: number;
  firstName: string;
  lastName: string;
  email: string | null;
  phone: string | null;
  idType: string | null;
  idNumber: string | null;
  address: string | null;
  city: string | null;
  country: string | null;
  notes: string | null;
  createdAt: string;
  updatedAt: string;
}

const route = useRoute();
const guestId = Number(route.params.id);

const loading = ref(false);
const editing = ref(false);
const guest = ref<Guest | null>(null);

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

const idTypeOptions = computed(() => [
  { value: "passport", label: t("idTypes.passport") },
  { value: "national_id", label: t("idTypes.nationalId") },
  { value: "driver_license", label: t("idTypes.driverLicense") },
  { value: "other", label: t("idTypes.other") },
]);

const fetchGuest = async () => {
  loading.value = true;
  try {
    const response = await $fetch<{ data: Guest }>(`/api/guests/${guestId}`);
    guest.value = response.data;
    resetForm();
  } catch (error) {
    console.error("Failed to fetch guest:", error);
    guest.value = null;
  } finally {
    loading.value = false;
  }
};

const resetForm = () => {
  if (guest.value) {
    form.firstName = guest.value.firstName || "";
    form.lastName = guest.value.lastName || "";
    form.email = guest.value.email || "";
    form.phone = guest.value.phone || "";
    form.idType = guest.value.idType || undefined;
    form.idNumber = guest.value.idNumber || "";
    form.address = guest.value.address || "";
    form.city = guest.value.city || "";
    form.country = guest.value.country || "";
    form.notes = guest.value.notes || "";
  }
};

const cancelEdit = () => {
  editing.value = false;
  resetForm();
};

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const formatIdType = (idType: string | null): string => {
  if (!idType) return "";
  const types: Record<string, string> = {
    passport: t("idTypes.passport"),
    national_id: t("idTypes.nationalId"),
    driver_license: t("idTypes.driverLicense"),
    other: t("idTypes.other"),
  };
  return types[idType] || idType;
};

const handleSubmit = async () => {
  loading.value = true;
  try {
    const body = {
      firstName: form.firstName,
      lastName: form.lastName,
      email: form.email || undefined,
      phone: form.phone || undefined,
      idType: form.idType,
      idNumber: form.idNumber || undefined,
      address: form.address || undefined,
      city: form.city || undefined,
      country: form.country || undefined,
      notes: form.notes || undefined,
    };

    await $fetch(`/api/guests/${guestId}`, {
      method: "put",
      body,
    });

    editing.value = false;
    await fetchGuest();
  } catch (error) {
    console.error("Failed to update guest:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchGuest);
</script>
