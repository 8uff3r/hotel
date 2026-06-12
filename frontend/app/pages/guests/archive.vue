<template>
  <div class="p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl font-bold">{{ t("guests.archiveTitle") }}</h1>
      <div class="flex gap-2">
        <USelect v-model="statusFilter" :items="statusOptions" :placeholder="t('guests.status')" />
        <UButton to="/guests">{{ t("actions.back") }}</UButton>
      </div>
    </div>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("guests.archiveList") }}</span>
          <span class="text-sm text-gray-500"
            >{{ (guests ?? []).length }} {{ t("guests.count") }}</span
          >
        </div>
      </template>
      <UTable :rows="guests" :columns="columns" :loading="pending">
        <template #actions-data="{ row }">
          <UButton variant="ghost" icon="i-lucide-eye" size="xs" :to="`/guests/${row.id}`" />
        </template>
      </UTable>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { Guest } from "~/utils/client";

definePageMeta({
  requiresPermission: "guests:read",
});

const { t } = useI18n();

const columns = [
  { accessorKey: "firstName", header: t("guests.firstName") },
  { accessorKey: "lastName", header: t("guests.lastName") },
  { accessorKey: "nationalId", header: t("guests.nationalId") },
  { accessorKey: "passport", header: t("guests.passport") },
  { accessorKey: "status", header: t("guests.status") },
  { id: "actions", header: t("actions.actions") },
];

const statusFilter = ref("");

const statusOptions = [
  { label: t("common.all"), value: "" },
  { label: t("guests.checkedOut"), value: "checked_out" },
  { label: t("guests.cancelled"), value: "cancelled" },
  { label: t("guests.absence"), value: "absence" },
];

const { data: guests, pending } = useFetch(
  () => `/api/guests/archived?status=${statusFilter.value}`,
  {
    key: "guests-archive",
    transform: (res) => (res as { data?: Guest[] })?.data ?? [],
  }
);
</script>
