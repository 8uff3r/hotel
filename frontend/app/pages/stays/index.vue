<template>
  <div class="p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl font-bold">{{ t("stays.title") }}</h1>
      <UButton to="/stays/create" icon="i-lucide-plus"> {{ t("stays.create") }} </UButton>
    </div>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("stays.list") }}</span>
          <span class="text-sm text-gray-500"
            >{{ (stays ?? []).length }} {{ t("stays.count") }}</span
          >
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

const { data: stays, pending } = useFetch("/api/stays", {
  key: "stays-list",
  transform: (res) => (res as { data?: Stay[] })?.data ?? [],
});
</script>
