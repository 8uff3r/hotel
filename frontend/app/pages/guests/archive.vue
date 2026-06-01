<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import type { Guest } from "~/utils/client";

const { t } = useI18n();
const columns = computed<TableColumn<Guest>[]>(() => [
  { accessorKey: "id", header: t("guests.columns.id") },
  { accessorKey: "name", header: t("guests.columns.name") },
  { accessorKey: "phone", header: t("guests.columns.phone") },
  { accessorKey: "actions", header: t("guests.columns.actions") },
]);

const deleting = ref(false);
const deleteModalOpen = ref(false);
const selectedGuest = ref<Guest | null>(null);
const page = ref(1);

const filters = reactive({
  search: "",
});

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  totalPages: 0,
});

let searchTimeout: ReturnType<typeof setTimeout> | null = null;

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    pagination.page = 1;
    refresh();
  }, 300);
};

const {
  data: guests,
  isPending: pending,
  refresh,
} = useQuery({
  key: () => ["guests", "archived", "list", page.value],
  query: async () => {
    const response = await getApiGuestsArchived({
      query: {
        page: page.value,
      },
    });
    pagination.total = response.data?.total ?? 0;
    pagination.totalPages = response.data?.totalPages ?? 0;
    return response.data?.data;
  },
});

const clearFilters = () => {
  filters.search = "";
  pagination.page = 1;
  refresh();
};

const confirmDelete = (guest: Guest) => {
  selectedGuest.value = guest;
  deleteModalOpen.value = true;
};

const deleteGuest = async () => {
  if (!selectedGuest.value) return;

  deleting.value = true;
  try {
    await $fetch(`/api/guests/${selectedGuest.value.id}`, { method: "DELETE" });
    deleteModalOpen.value = false;
  } catch (error) {
    console.error("Failed to delete guest:", error);
  } finally {
    deleting.value = false;
  }
};
</script>
<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("guests.title") }}</h1>
      <UButton to="/guests/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />
        {{ t("guests.addGuest") }}
      </UButton>
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          :placeholder="t('guests.searchPlaceholder')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <UButton variant="outline" @click="clearFilters"> {{ t("actions.clear") }} </UButton>
      </div>
    </UCard>

    <!-- Guests Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("guests.list") }}</span>
          <span class="text-sm text-gray-500">{{
            t("guests.count", { count: pagination.total })
          }}</span>
        </div>
      </template>

      <UTable :data="guests ?? []" :columns="columns" :loading="pending" striped>
        <template #id-cell="{ row }">
          <NuxtLink
            :to="`/guests/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            #{{ row.original.id }}
          </NuxtLink>
        </template>

        <template #name-cell="{ row }">
          <div>
            <p class="font-medium">{{ row.original.firstName }} {{ row.original.lastName }}</p>
          </div>
        </template>

        <template #phone-cell="{ row }">
          {{ row.original.phone || "-" }}
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/guests/${row.original.id}`">
              <UIcon name="i-lucide-eye" class="h-4 w-4" />
            </UButton>
            <UButton variant="ghost" size="sm" :to="`/guests/${row.original.id}`">
              <UIcon name="i-lucide-pencil" class="h-4 w-4" />
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

    <!-- Delete Confirmation Modal -->
    <UModal v-model="deleteModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("actions.confirmDelete") }}</h2>
      </template>
      <template #body>
        <p>{{ t("guests.confirmDelete", { name: selectedGuest?.firstName }) }}</p>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="deleteModalOpen = false">{{
            t("actions.cancel")
          }}</UButton>
          <UButton color="error" :loading="deleting" @click="deleteGuest">{{
            t("actions.delete")
          }}</UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>
