<template>
  <div>
    <div class="mb-6">
      <UButton to="/travel-agencies" variant="ghost" size="sm" class="mb-2">
        <UIcon name="i-lucide-arrow-left" class="mr-1" />{{ t("actions.back") }}
      </UButton>
      <h1 class="text-3xl font-bold">{{ t("travelAgency.createAgency") }}</h1>
    </div>
    <UCard>
      <form @submit.prevent="createAgency">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("travelAgency.name") }} *</label>
            <UInput v-model="form.name" :placeholder="t('travelAgency.namePlaceholder')" required />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("travelAgency.status") }}</label>
            <USelect v-model="form.status" :items="statusOptions" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium">{{
              t("travelAgency.ceoFirstName")
            }}</label>
            <UInput
              v-model="form.ceoFirstName"
              :placeholder="t('travelAgency.ceoFirstNamePlaceholder')"
            />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium">{{
              t("travelAgency.ceoLastName")
            }}</label>
            <UInput
              v-model="form.ceoLastName"
              :placeholder="t('travelAgency.ceoLastNamePlaceholder')"
            />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("travelAgency.province") }}</label>
            <UInput v-model="form.province" :placeholder="t('travelAgency.provincePlaceholder')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium">{{ t("travelAgency.city") }}</label>
            <UInput v-model="form.city" :placeholder="t('travelAgency.cityPlaceholder')" />
          </div>
        </div>
        <div class="mt-6 flex justify-end gap-3">
          <UButton type="button" variant="outline" to="/travel-agencies">{{
            t("common.cancel")
          }}</UButton>
          <UButton type="submit" color="primary" :loading="loading">{{
            t("travelAgency.createAgency")
          }}</UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { postApiTravelAgencies } from "~/utils/client";

const { t } = useI18n();
const router = useRouter();
const loading = ref(false);

const statusOptions = [
  { label: t("travelAgency.enabled"), value: "enabled" },
  { label: t("travelAgency.disabled"), value: "disabled" },
];

const form = reactive({
  name: "",
  ceoFirstName: "",
  ceoLastName: "",
  province: "",
  city: "",
  status: "enabled",
});

const createAgency = async () => {
  loading.value = true;
  try {
    await postApiTravelAgencies({
      body: { ...form },
    });
    router.push("/travel-agencies");
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false;
  }
};
</script>
