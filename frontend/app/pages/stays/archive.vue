<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">{{ t("stays.archiveTitle") }}</h1>
      <div class="flex gap-2">
        <USelect v-model="statusFilter" :items="statusOptions" :placeholder="t('stays.status')" />
        <USelect v-model="settlementFilter" :items="settlementOptions" :placeholder="t('stays.settlementStatus')" />
        <UButton to="/stays">{{ t("actions.back") }}</UButton>
      </div>
    </div>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("stays.archiveList") }}</span>
          <span class="text-sm text-gray-500">{{ (stays ?? []).length }} {{ t("stays.count") }}</span>
        </div>
      </template>
      <UTable :rows="stays" :columns="columns" :loading="pending">
        <template #actions-data="{ row }">
          <UButton variant="ghost" icon="i-lucide-eye" size="xs" :to="`/stays/${row.id}`" />
        </template>
      </UTable>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import type { Stay } from "~/utils/client";

definePageMeta({
  requiresPermission: "guests:read",
});

const { t } = useI18n();

const columns = [
  { accessorKey: "acceptanceId", header: t("stays.acceptanceId") },
  { accessorKey: "guest.firstName", header: t("stays.guestFirstName") },
  { accessorKey: "guest.lastName", header: t("stays.guestLastName") },
  { accessorKey: "room.roomNumber", header: t("stays.roomNumber") },
  { accessorKey: "entryDate", header: t("stays.entryDate") },
  { accessorKey: "departureDate", header: t("stays.departureDate") },
  { accessorKey: "status.label", header: t("stays.status") },
  { id: "actions", header: t("actions.actions") },
];

const statusFilter = ref("");
const settlementFilter = ref("");

const statusOptions = [
  { label: t("common.all"), value: "" },
  { label: t("stays.checkedOut"), value: "checked_out" },
  { label: t("stays.cancelled"), value: "cancelled" },
  { label: t("stays.absence"), value: "absence" },
];

const settlementOptions = [
  { label: t("common.all"), value: "" },
  { label: t("stays.cleared"), value: "cleared" },
  { label: t("stays.unsettled"), value: "unsettled" },
];

const queryParams = computed(() => {
  const params = new URLSearchParams();
  if (statusFilter.value) params.append("status", statusFilter.value);
  if (settlementFilter.value) params.append("settlement", settlementFilter.value);
  return params.toString();
});

const { data: stays, pending } = useFetch(() => `/api/stays?${queryParams.value}`, {
  key: "stays-archive",
  transform: (res) => (res as { data?: Stay[] })?.data ?? [],
});
</script>
