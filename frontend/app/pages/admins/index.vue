<template>
  <div class="p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl font-bold">{{ t("admins.title") }}</h1>
      <UButton to="/admins/create" icon="i-lucide-plus">
        {{ t("admins.create") }}
      </UButton>
    </div>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("admins.list") }}</span>
          <span class="text-sm text-gray-500"
            >{{ (admins ?? []).length }} {{ t("admins.count") }}</span
          >
        </div>
      </template>
      <UTable :rows="admins" :columns="columns" :loading="pending">
        <template #actions-data="{ row }">
          <UButton variant="ghost" icon="i-lucide-pencil" size="xs" :to="`/admins/${row.id}`" />
        </template>
      </UTable>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { SanitizedAdmin } from "~/utils/client";

definePageMeta({
  requiresPermission: "users:read",
});

const { t } = useI18n();

const columns = [
  { accessorKey: "firstName", header: t("admins.firstName") },
  { accessorKey: "lastName", header: t("admins.lastName") },
  { accessorKey: "email", header: t("admins.email") },
  { accessorKey: "username", header: t("admins.username") },
  { accessorKey: "role", header: t("admins.role") },
  { accessorKey: "isSuperAdmin", header: t("admins.isSuperAdmin") },
  { id: "actions", header: t("actions.actions") },
];

const { data: admins, pending } = useFetch("/api/admins", {
  key: "admins-list",
  transform: (res) => (res as { data?: SanitizedAdmin[] })?.data ?? [],
});
</script>
