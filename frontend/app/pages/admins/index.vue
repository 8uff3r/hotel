<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">مدیریت مدیران سیستم</h1>
      <UButton to="/admins/create" icon="i-lucide-plus">
        افزودن مدیر
      </UButton>
    </div>

    <UTable :rows="admins" :columns="columns" :loading="pending">
      <template #actions-data="{ row }">
        <UButton variant="ghost" icon="i-lucide-pencil" size="xs" :to="`/admins/${row.id}`" />
      </template>
    </UTable>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresPermission: "users:read",
});

const columns = [
  { key: "firstName", label: "نام" },
  { key: "lastName", label: "نام خانوادگی" },
  { key: "email", label: "ایمیل" },
  { key: "username", label: "نام کاربری" },
  { key: "role", label: "نقش" },
  { key: "isSuperAdmin", label: "سوپر ادمین" },
  { key: "actions", label: "عملیات" },
];

const { data: admins, pending } = useFetch("/api/admins", {
  key: "admins-list",
  transform: (res: any) => res?.data ?? [],
});
</script>
