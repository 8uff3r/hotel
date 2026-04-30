<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("users.title") }}</h1>
      <UButton to="/users/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />
        {{ t("users.addUser") }}
      </UButton>
    </div>

    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          :placeholder="t('users.searchPlaceholder')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <UButton variant="outline" @click="clearFilters">{{ t("actions.clear") }}</UButton>
      </div>
    </UCard>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("users.list") }}</span>
          <span class="text-sm text-gray-500">{{ t("users.count", { count: users.length }) }}</span>
        </div>
      </template>

      <UTable :data="users" :columns="columns" :loading="loading" striped>
        <template #id-cell="{ row }">
          <NuxtLink
            :to="`/users/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            #{{ row.original.id }}
          </NuxtLink>
        </template>

        <template #name-cell="{ row }">
          <div>
            <p class="font-medium">{{ row.original.firstName }} {{ row.original.lastName }}</p>
            <p v-if="row.original.email" class="text-sm text-gray-500">{{ row.original.email }}</p>
          </div>
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/users/${row.original.id}`">
              <UIcon name="i-lucide-eye" class="h-4 w-4" />
            </UButton>
            <UButton variant="ghost" size="sm" :to="`/users/${row.original.id}/edit`">
              <UIcon name="i-lucide-pencil" class="h-4 w-4" />
            </UButton>
          </div>
        </template>
      </UTable>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import type { PaginatedResponseModelsSanitizedUser } from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.users.users.read,
});

type User = NonNullable<PaginatedResponseModelsSanitizedUser["data"]>[0];

const { t } = useI18n();
const columns = computed<TableColumn<User>[]>(() => [
  { accessorKey: "id", header: t("users.columns.id") },
  { accessorKey: "name", header: t("users.columns.name") },
  { accessorKey: "actions", header: t("users.columns.actions") },
]);

const sourceUsers = ref<User[]>([]);
const users = ref<User[]>([]);
const loading = ref(false);

const filters = reactive({
  search: "",
});

let searchTimeout: ReturnType<typeof setTimeout> | null = null;

const applyFilters = () => {
  const query = filters.search.trim().toLowerCase();
  if (!query) {
    users.value = [...sourceUsers.value];
    return;
  }

  users.value = sourceUsers.value.filter((u) => {
    const fullName = `${u.firstName ?? ""} ${u.lastName ?? ""}`.toLowerCase();
    return fullName.includes(query) || (u.email ?? "").toLowerCase().includes(query);
  });
};

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(applyFilters, 300);
};

const fetchUsers = async () => {
  loading.value = true;
  try {
    const response = await getApiUsers({});
    sourceUsers.value = response.data ?? [];
    applyFilters();
  } catch (error) {
    console.error("Failed to fetch users:", error);
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  filters.search = "";
  applyFilters();
};

onMounted(fetchUsers);
</script>
