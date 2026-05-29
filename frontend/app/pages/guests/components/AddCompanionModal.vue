<script setup lang="ts" generic="T">
import type { FormSubmitEvent } from "@nuxt/ui";
import { companionSchema, type Companion } from "../utils";
import { getApiGuestsRelations } from "~/utils/client";

const props = defineProps<{
  loading: boolean;
}>();
const state = defineModel<Companion>("state", { default: {} });
const open = defineModel<boolean>("open", { default: false });

const { t } = useI18n();
const { data: relations } = useAsyncData("family-relations", async () => {
  const res = await getApiGuestsRelations({ query: { limit: -1 } });
  return res.data;
});

const emit = defineEmits<{
  submit: [FormSubmitEvent<Companion>];
}>();

const { data: countries } = useAsyncData("countries", async () => {
  const res = await getApiCommonCountries({ query: { limit: -1 } });
  return res.data;
});
</script>
<template>
  <UModal :title="t('guests.addCompanion')" v-model:open="open">
    <template #content>
      <UForm @submit="(e) => emit('submit', e)" :schema="companionSchema" :state class="p-4">
        <div class="space-y-4">
          <UFormField :label="t('guests.companionRelation')" name="relation" required>
            <HSelect v-model="state.relation" :items="relations ?? []" :disabled="loading" />
          </UFormField>
          <div class="grid grid-cols-2 gap-4">
            <UFormField :label="t('forms.firstName')" name="firstName" required>
              <UInput v-model="state.firstName" :disabled="loading" />
            </UFormField>
            <UFormField :label="t('forms.lastName')" name="lastName" required>
              <UInput v-model="state.lastName" :disabled="loading" />
            </UFormField>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <UFormField :label="t('guest.gender')" name="gender">
              <USelect
                v-model="state.gender"
                class="w-full"
                :items="[
                  { value: 'male', label: t('guest.male') },
                  { value: 'female', label: t('guest.female') },
                ]"
                :disabled="loading"
              />
            </UFormField>
            <UFormField :label="t('guest.nationality')" name="nationalityID">
              <HSelectMenu
                v-model="state.nationalityID"
                :items="countries ?? []"
                :disabled="loading"
                :placeholder="t('guest.defaultNationality')"
              />
            </UFormField>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <UFormField :label="t('guest.nationalId')" name="nationalId">
              <UInput v-model="state.nationalId" :disabled="loading" />
            </UFormField>
            <UFormField :label="t('guest.idPassportNumber')" name="idNumber">
              <UInput v-model="state.idNumber" :disabled="loading" />
            </UFormField>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <UFormField :label="t('guest.dateOfBirth')" name="dateOfBirth">
              <HDate v-model="state.dateOfBirth" />
            </UFormField>
            <UFormField :label="t('guest.phone')" name="phone">
              <UInput v-model="state.phone" :disabled="loading" />
            </UFormField>
          </div>
        </div>
        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" type="button" @click="open = false">
            {{ t("actions.cancel") }}
          </UButton>
          <UButton type="submit" color="primary" :loading="loading">
            {{ t("guests.addCompanion") }}
          </UButton>
        </div>
      </UForm>
    </template>
  </UModal>
</template>
