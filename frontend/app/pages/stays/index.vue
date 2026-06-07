<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">پذیرش / اقامت‌ها</h1>
      <UButton to="/stays/create" icon="i-lucide-plus">
        ثبت پذیرش جدید
      </UButton>
    </div>

    <UTable :rows="stays" :columns="columns" :loading="pending">
      <template #actions-data="{ row }">
        <UButton variant="ghost" icon="i-lucide-eye" size="xs" :to="`/stays/${row.id}`" />
      </template>
    </UTable>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresPermission: "guests:read",
});

const columns = [
  { key: "acceptanceId", label: "شناسه پذیرش" },
  { key: "guest.firstName", label: "نام مهمان" },
  { key: "guest.lastName", label: "نام خانوادگی" },
  { key: "room.roomNumber", label: "شماره اتاق" },
  { key: "entryDate", label: "تاریخ ورود" },
  { key: "departureDate", label: "تاریخ خروج" },
  { key: "status.label", label: "وضعیت" },
  { key: "actions", label: "عملیات" },
];

const { data: stays, pending } = useFetch("/api/stays", {
  key: "stays-list",
  transform: (res: any) => res?.data ?? [],
});
</script>
