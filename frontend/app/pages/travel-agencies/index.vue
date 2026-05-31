<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <UButton to="/" variant="ghost" size="sm" class="mb-2">
          <UIcon name="i-lucide-arrow-left" class="mr-1" />{{ t("actions.back") }}
        </UButton>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          {{ t("travelAgency.title") }}
        </h1>
      </div>
      <UButton to="/travel-agencies/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />{{ t("travelAgency.addAgency") }}
      </UButton>
    </div>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("travelAgency.list") }}</span>
          <span class="text-sm text-gray-500">{{
            t("travelAgency.count", { count: pagination.total })
          }}</span>
        </div>
      </template>
      <UTable :data="agencies" :columns="columns" :loading="loading" striped>
        <template #name-cell="{ row }">
          <NuxtLink
            :to="`/travel-agencies/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            {{ row.original.name }}
          </NuxtLink>
        </template>
        <template #ceo-cell="{ row }">
          {{ row.original.ceoFirstName }} {{ row.original.ceoLastName }}
        </template>
        <template #location-cell="{ row }">
          <template v-if="row.original.city && row.original.province">
            {{ row.original.city }}, {{ row.original.province }}
          </template>
          <template v-else-if="row.original.city">{{ row.original.city }}</template>
          <template v-else-if="row.original.province">{{ row.original.province }}</template>
        </template>
        <template #status-cell="{ row }">
          <UBadge
            :color="row.original.status === 'enabled' ? 'success' : 'neutral'"
            variant="subtle"
          >
            {{
              row.original.status === "enabled"
                ? t("travelAgency.enabled")
                : t("travelAgency.disabled")
            }}
          </UBadge>
        </template>
        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/travel-agencies/${row.original.id}`">
              <UIcon name="i-lucide-eye" class="h-4 w-4" />
            </UButton>
            <UButton variant="ghost" size="sm" color="error" @click="confirmDelete(row.original)">
              <UIcon name="i-lucide-trash-2" class="h-4 w-4" />
            </UButton>
          </div>
        </template>
      </UTable>
      <template #footer>
        <div class="flex items-center justify-between">
          <span class="text-sm text-gray-500">
            {{
              t("pagination.pageOf", { page: pagination.page, totalPages: pagination.totalPages })
            }}
          </span>
          <UPagination v-model="page" :page-count="pagination.limit" :total="pagination.total" />
        </div>
      </template>
    </UCard>

    <UModal v-model="deleteModalOpen">
      <template #header
        ><h2 class="text-lg font-semibold">{{ t("actions.confirmDelete") }}</h2></template
      >
      <template #body>
        <p>{{ t("travelAgency.confirmDelete", { name: selectedAgency?.name }) }}</p>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="deleteModalOpen = false">{{
            t("actions.cancel")
          }}</UButton>
          <UButton color="error" :loading="deleting" @click="deleteAgency">{{
            t("actions.delete")
          }}</UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import { getApiTravelAgencies } from "~/utils/client";

definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

const { t } = useI18n();

const columns: TableColumn<any>[] = [
  { accessorKey: "name", header: t("travelAgency.columns.name") },
  { accessorKey: "ceo", header: t("travelAgency.columns.ceo") },
  { accessorKey: "location", header: t("travelAgency.columns.location") },
  { accessorKey: "status", header: t("travelAgency.columns.status") },
  { accessorKey: "actions", header: t("travelAgency.columns.actions") },
];

const deleting = ref(false);
const deleteModalOpen = ref(false);
const selectedAgency = ref<any>(null);
const page = ref(1);

const pagination = reactive({ page: 1, limit: 10, total: 0, totalPages: 0 });

const {
  data: agencies,
  isLoading: loading,
  refetch,
} = useQuery({
  key: ["agencies", "list", pagination],
  query: async () => {
    const response = await getApiTravelAgencies({
      query: pagination,
    });
    agencies.value = response.data?.data || [];
    pagination.total = response.data?.total ?? 0;
    pagination.totalPages = response.data?.totalPages ?? 0;
    return response.data?.data;
  },
});

const confirmDelete = (agency: any) => {
  selectedAgency.value = agency;
  deleteModalOpen.value = true;
};

const deleteAgency = async () => {
  if (!selectedAgency.value) return;
  deleting.value = true;
  try {
    await $fetch(`/api/travel-agencies/${selectedAgency.value.id}`, { method: "DELETE" });
    deleteModalOpen.value = false;
    refetch();
  } catch (error) {
    console.error(error);
  } finally {
    deleting.value = false;
  }
};
</script>
