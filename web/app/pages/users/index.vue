<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("users.title") }}</h1>
      <UButton to="/users/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />
        {{ t("users.addUser") }}
      </UButton>
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          :placeholder="t('users.searchPlaceholder')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.role"
          :items="roleOptions"
          :placeholder="t('users.allRoles')"
          class="w-full sm:w-40"
          @change="fetchUsers"
        />
        <USelect
          v-model="filters.status"
          :items="statusOptions"
          :placeholder="t('users.allStatus')"
          class="w-full sm:w-40"
          @change="fetchUsers"
        />
        <UButton variant="outline" @click="clearFilters"> {{ t("actions.clear") }} </UButton>
      </div>
    </UCard>

    <!-- Users Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("users.list") }}</span>
          <span class="text-sm text-gray-500">{{ t("users.count", { count: pagination.total }) }}</span>
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
            <p v-if="row.original.email" class="text-sm text-gray-500">
              {{ row.original.email }}
            </p>
          </div>
        </template>

        <template #roles-cell="{ row }">
          <div class="flex flex-wrap gap-1">
            <UBadge
              v-for="role in row.original.roles"
              :key="role"
              :color="getRoleColor(role)"
              variant="soft"
              size="sm"
            >
              {{ formatRole(role) }}
            </UBadge>
          </div>
        </template>

        <template #status-cell="{ row }">
          <UBadge :color="row.original.isActive ? 'success' : 'error'" variant="soft">
            {{ row.original.isActive ? t("statuses.active") : t("statuses.inactive") }}
          </UBadge>
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

      <template #footer>
        <div class="flex items-center justify-between">
          <span class="text-sm text-gray-500">
            {{ t("pagination.pageOf", { page: pagination.page, totalPages: pagination.totalPages }) }}
          </span>
          <UPagination
            v-model="page"
            :page-count="pagination.limit"
            :total="pagination.total"
            @change="fetchUsers"
          />
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";

definePageMeta({
  requiresRole: ["admin"],
});

interface UserRow {
  id: number;
  email: string;
  firstName: string;
  lastName: string;
  isActive: boolean;
  roles: string[];
}

const { t } = useI18n();
const columns = computed<TableColumn<UserRow>[]>(() => [
  { accessorKey: "id", header: t("users.columns.id") },
  { accessorKey: "name", header: t("users.columns.name") },
  { accessorKey: "roles", header: t("users.columns.roles") },
  { accessorKey: "status", header: t("users.columns.status") },
  { accessorKey: "actions", header: t("users.columns.actions") },
]);

const roleOptions = computed(() => [
  { value: "all", label: t("users.allRoles") },
  { value: "admin", label: t("roles.admin") },
  { value: "manager", label: t("roles.manager") },
  { value: "receptionist", label: t("roles.receptionist") },
  { value: "staff", label: t("roles.staff") },
]);

const statusOptions = computed(() => [
  { value: "all", label: t("users.allStatus") },
  { value: "active", label: t("statuses.active") },
  { value: "inactive", label: t("statuses.inactive") },
]);

const users = ref<UserRow[]>([]);
const loading = ref(false);
const page = ref(1);

const filters = reactive({
  search: "",
  role: "",
  status: "",
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
    fetchUsers();
  }, 300);
};

const fetchUsers = async () => {
  loading.value = true;
  try {
    const params = new URLSearchParams();
    params.append("page", pagination.page.toString());
    params.append("limit", pagination.limit.toString());

    if (filters.search) params.append("search", filters.search);

    const response = await $fetch(`/api/users?${params.toString()}`);

    let filteredData = response.data;

    if (filters.role && filters.role !== "all") {
      filteredData = filteredData.filter((u: UserRow) => u.roles.includes(filters.role));
    }

    if (filters.status && filters.status !== "all") {
      const isActive = filters.status === "active";
      filteredData = filteredData.filter((u: UserRow) => u.isActive === isActive);
    }

    users.value = filteredData;
    pagination.total = response.pagination.total;
    pagination.totalPages = response.pagination.totalPages;
  } catch (error) {
    console.error("Failed to fetch users:", error);
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  filters.search = "";
  filters.role = "";
  filters.status = "";
  pagination.page = 1;
  fetchUsers();
};

const formatRole = (role: string): string => {
  const roles: Record<string, string> = {
    admin: t("roles.admin"),
    manager: t("roles.manager"),
    receptionist: t("roles.receptionist"),
    staff: t("roles.staff"),
  };
  return roles[role] || role;
};

const getRoleColor = (role: string): "success" | "info" | "warning" | "error" | "neutral" => {
  const colors: Record<string, "success" | "info" | "warning" | "error" | "neutral"> = {
    admin: "error",
    manager: "info",
    receptionist: "success",
    staff: "neutral",
  };
  return colors[role] || "neutral";
};

onMounted(fetchUsers);
</script>
