<script setup lang="ts">
import { getApiCommonCountries } from "~/utils/client";
import type { CreateRequest } from "../utils";

defineProps<{
  loading: boolean;
}>();
const form = defineModel<Required<CreateRequest>>({ default: {} });
const { t } = useI18n();

const { data: countries } = useAsyncData("countries", async () => {
  const res = await getApiCommonCountries({ query: { limit: -1 } });
  return res.data;
});
</script>
<template>
  <UFormField :label="t('forms.firstName')" name="guest.firstName" required>
    <UInput v-model="form.guest.firstName" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('forms.lastName')" name="guest.lastName" required>
    <UInput v-model="form.guest.lastName" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.fatherName')" name="guest.fatherName">
    <UInput v-model="form.guest.fatherName" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.nationality')" name="guest.nationalityID">
    <HSelectMenu
      v-model="form.guest.nationalityID"
      :items="countries ?? []"
      :disabled="loading"
      :placeholder="t('guest.defaultNationality')"
    />
  </UFormField>

  <UFormField :label="t('guest.gender')" name="guest.gender">
    <USelect
      v-model="form.guest.gender"
      class="w-full"
      :items="[
        { value: 'male', label: t('guest.male') },
        { value: 'female', label: t('guest.female') },
      ]"
      :disabled="loading"
    />
  </UFormField>

  <UFormField :label="t('guest.dateOfBirth')" name="guest.dateOfBirth">
    <HDate v-model="form.guest.dateOfBirth" />
  </UFormField>

  <UFormField :label="t('guest.placeOfBirth')" name="guest.placeOfBirth">
    <UInput v-model="form.guest.placeOfBirth" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.nationalId')" name="guest.nationalId">
    <UInput v-model="form.guest.nationalId" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.idPassportNumber')" name="guest.idNumber">
    <UInput v-model="form.guest.idNumber" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.occupation')" name="guest.occupation">
    <UInput
      v-model="form.guest.occupation"
      :disabled="loading"
      :placeholder="t('guest.defaultOccupation')"
    />
  </UFormField>

  <UFormField :label="t('guest.phone')" name="guest.phone">
    <UInput v-model="form.guest.phone" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('guest.postalCode')" name="guest.postalCode">
    <UInput v-model="form.guest.postalCode" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('forms.address')" name="guest.address" class="md:col-span-3">
    <UInput class="w-full" v-model="form.guest.address" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('forms.email')" name="guest.email" class="md:col-span-3">
    <UInput class="w-full" v-model="form.guest.email" :disabled="loading" />
  </UFormField>

  <UFormField :label="t('forms.landline')" name="guest.landline" class="md:col-span-3">
    <UInput class="w-full" v-model="form.guest.landline" :disabled="loading" />
  </UFormField>
</template>
