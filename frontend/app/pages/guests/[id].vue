<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/guests" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />{{ t("actions.backToGuests") }}</UButton
      >
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ t("guests.guestNumber", { id: guest?.id }) }}
      </h1>
    </div>

    <div v-if="pending" class="flex justify-center py-12">
      <UIcon name="i-lucide-loader-2" class="h-8 w-8 animate-spin" />
    </div>

    <div v-else-if="guest">
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- Main Form -->
        <div class="lg:col-span-2">
          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <span class="text-lg font-semibold">{{ t("guest.personalInfo") }}</span>
              </div>
            </template>

            <form @submit.prevent="handleSubmit">
              <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
                <!-- First Name -->
                <UFormField :label="t('forms.firstName')" name="firstName" required>
                  <UInput v-model="form.firstName" :disabled="saving" />
                </UFormField>

                <!-- Last Name -->
                <UFormField :label="t('forms.lastName')" name="lastName" required>
                  <UInput v-model="form.lastName" :disabled="saving" />
                </UFormField>

                <!-- Father Name -->
                <UFormField :label="t('guest.fatherName')" name="fatherName">
                  <UInput v-model="form.fatherName" :disabled="saving" />
                </UFormField>

                <!-- Nationality -->
                <UFormField :label="t('guest.nationality')" name="nationalityID" required>
                  <HSelectMenu
                    v-model="form.nationalityID"
                    :items="countries ?? []"
                    :disabled="saving"
                    :placeholder="t('guest.defaultNationality')"
                  />
                </UFormField>

                <!-- Gender -->
                <UFormField :label="t('guest.gender')" name="gender">
                  <USelect
                    v-model="form.gender"
                    class="w-full"
                    :items="[
                      { value: 'male', label: t('guest.male') },
                      { value: 'female', label: t('guest.female') },
                    ]"
                    :disabled="saving"
                  />
                </UFormField>

                <!-- Date of Birth -->
                <UFormField :label="t('guest.dateOfBirth')" name="dateOfBirth">
                  <HDate v-model="form.dateOfBirth" />
                </UFormField>

                <!-- Place of Birth -->
                <UFormField :label="t('guest.placeOfBirth')" name="placeOfBirth">
                  <UInput v-model="form.placeOfBirth" :disabled="saving" />
                </UFormField>

                <!-- National ID -->
                <UFormField :label="t('guest.nationalId')" name="nationalId">
                  <UInput v-model="form.nationalId" :disabled="saving" />
                </UFormField>

                <!-- ID/Passport Number -->
                <UFormField :label="t('guest.idPassportNumber')" name="idNumber">
                  <UInput v-model="form.idNumber" :disabled="saving" />
                </UFormField>

                <!-- Occupation -->
                <UFormField :label="t('guest.occupation')" name="occupation">
                  <UInput
                    v-model="form.occupation"
                    :disabled="saving"
                    :placeholder="t('guest.defaultOccupation')"
                  />
                </UFormField>

                <!-- Phone -->
                <UFormField :label="t('guest.phone')" name="phone">
                  <UInput v-model="form.phone" :disabled="saving" />
                </UFormField>

                <!-- Email -->
                <UFormField :label="t('forms.email')" name="email">
                  <UInput v-model="form.email" :disabled="saving" />
                </UFormField>

                <!-- Landline -->
                <UFormField :label="t('forms.landline')" name="landline">
                  <UInput v-model="form.landline" :disabled="saving" />
                </UFormField>

                <!-- Postal Code -->
                <UFormField :label="t('guest.postalCode')" name="postalCode">
                  <UInput v-model="form.postalCode" :disabled="saving" />
                </UFormField>

                <!-- Address -->
                <UFormField :label="t('forms.address')" name="address" class="md:col-span-2">
                  <UInput v-model="form.address" :disabled="saving" />
                </UFormField>
              </div>

              <div class="mt-6 flex justify-end gap-3">
                <UButton variant="outline" to="/guests" :disabled="saving">
                  {{ t("actions.cancel") }}
                </UButton>
                <UButton type="submit" color="primary" :loading="saving">
                  {{ t("actions.saveChanges") }}
                </UButton>
              </div>
            </form>
          </UCard>
        </div>

        <!-- Sidebar Info -->
        <div class="space-y-6">
          <!-- Quick Info -->
          <UCard>
            <template #header>
              <span class="font-semibold">{{ t("guests.quickInfo") }}</span>
            </template>
            <div class="space-y-4">
              <div class="flex items-center gap-3">
                <UIcon name="i-lucide-user" class="h-6 w-6 text-primary" />
                <div>
                  <p class="text-sm text-gray-500">{{ t("guests.fullName") }}</p>
                  <p class="font-medium">{{ guest.firstName }} {{ guest.lastName }}</p>
                </div>
              </div>
              <div v-if="guest.fatherName" class="flex items-center gap-3">
                <UIcon name="i-lucide-users" class="h-6 w-6 text-primary" />
                <div>
                  <p class="text-sm text-gray-500">{{ t('guest.fatherName') }}</p>
                  <p class="font-medium">{{ guest.fatherName }}</p>
                </div>
              </div>
              <div v-if="guest.nationality" class="flex items-center gap-3">
                <UIcon name="i-lucide-globe" class="h-6 w-6 text-primary" />
                <div>
                  <p class="text-sm text-gray-500">{{ t('guest.nationality') }}</p>
                  <p class="font-medium">{{ guest.nationality.label }}</p>
                </div>
              </div>
              <div v-if="guest.phone" class="flex items-center gap-3">
                <UIcon name="i-lucide-phone" class="h-6 w-6 text-primary" />
                <div>
                  <p class="text-sm text-gray-500">{{ t('guest.phone') }}</p>
                  <p class="font-medium">{{ guest.phone }}</p>
                </div>
              </div>
              <div v-if="guest.gender" class="flex items-center gap-3">
                <UIcon name="i-lucide-venus-and-mars" class="h-6 w-6 text-primary" />
                <div>
                  <p class="text-sm text-gray-500">{{ t('guest.gender') }}</p>
                  <p class="font-medium capitalize">{{ guest.gender }}</p>
                </div>
              </div>
            </div>
          </UCard>

          <!-- Companions -->
          <UCard v-if="guest.companions?.length">
            <template #header>
              <span class="font-semibold">{{ t("guests.companions") }}</span>
            </template>
            <div class="space-y-3">
              <div
                v-for="companion in guest.companions"
                :key="companion.id"
                class="flex items-center gap-3 rounded-lg border p-3"
              >
                <UIcon name="i-lucide-user" class="h-5 w-5 text-primary" />
                <div>
                  <p class="font-medium">
                    {{ companion.firstName }} {{ companion.lastName }}
                  </p>
                  <p v-if="companion.relation" class="text-sm text-gray-500">
                    {{ companion.relation.label }}
                  </p>
                </div>
              </div>
            </div>
          </UCard>
        </div>
      </div>
    </div>

    <!-- Not Found -->
    <div v-else class="py-12 text-center">
      <UIcon name="i-lucide-alert-circle" class="mx-auto h-12 w-12 text-gray-400" />
      <p class="mt-4 text-lg text-gray-500">{{ t("guests.notFound") }}</p>
      <UButton to="/guests" class="mt-4"> {{ t("actions.backToGuests") }} </UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Guest } from "~/utils/client";
import { useCountriesQuery } from "./queries";

definePageMeta({
  requiresPermission: PERMISSIONS.guests.guests.read,
});

const { t } = useI18n();
const route = useRoute();
const guestId = route.params.id as string;
const saving = ref(false);

const form = ref<Guest>({} as Guest);

const { data: countries } = useCountriesQuery();

const { data: guest, pending } = useAsyncData(async () => {
  const response = await getApiGuestsId({
    path: { id: guestId },
  });

  form.value = {
    firstName: response.data?.firstName ?? "",
    lastName: response.data?.lastName ?? "",
    fatherName: response.data?.fatherName ?? "",
    nationalityID: response.data?.nationalityID,
    gender: response.data?.gender ?? "",
    dateOfBirth: response.data?.dateOfBirth ?? "",
    placeOfBirth: response.data?.placeOfBirth ?? "",
    nationalId: response.data?.nationalId ?? "",
    idNumber: response.data?.idNumber ?? "",
    occupation: response.data?.occupation ?? "",
    phone: response.data?.phone ?? "",
    email: response.data?.email ?? "",
    landline: response.data?.landline ?? "",
    postalCode: response.data?.postalCode ?? "",
    address: response.data?.address ?? "",
  } as Guest;

  return response.data;
});

const handleSubmit = async () => {
  saving.value = true;
  try {
    await putApiGuestsId({
      path: { id: guestId },
      body: form.value,
    });

    await navigateTo("/guests");
  } catch (error) {
    console.error("Failed to update guest:", error);
  } finally {
    saving.value = false;
  }
};
</script>
