<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Users</h1>
      <UButton to="/users/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />
        Add User
      </UButton>
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          placeholder="Search by name or email..."
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.role"
          :items="roleOptions"
          placeholder="All Roles"
          class="w-full sm:w-40"
          @change="fetchUsers"
        />
        <USelect
          v-model="filters.status"
          :items="statusOptions"
          placeholder="All Status"
          class="w-full sm:w-40"
          @change="fetchUsers"
        />
        <UButton variant="outline" @click="clearFilters"> Clear </UButton>
      </div>
    </UCard>

    <!-- Users Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">User List</span>
          <span class="text-sm text-gray-500">{{ pagination.total }} users</span>
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
            {{ row.original.isActive ? "Active" : "Inactive" }}
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
            Page {{ pagination.page }} of {{ pagination.totalPages }}
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

const columns: TableColumn<UserRow>[] = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "name", header: "Name" },
  { accessorKey: "roles", header: "Roles" },
  { accessorKey: "status", header: "Status" },
  { accessorKey: "actions", header: "Actions" },
];

const roleOptions = [
  { value: "all", label: "All Roles" },
  { value: "admin", label: "Admin" },
  { value: "manager", label: "Manager" },
  { value: "receptionist", label: "Receptionist" },
  { value: "staff", label: "Staff" },
];

const statusOptions = [
  { value: "all", label: "All Status" },
  { value: "active", label: "Active" },
  { value: "inactive", label: "Inactive" },
];

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
    admin: "Admin",
    manager: "Manager",
    receptionist: "Receptionist",
    staff: "Staff",
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
