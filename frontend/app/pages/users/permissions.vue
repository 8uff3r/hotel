<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("permissions.title") }}</h1>
    </div>

    <UCard class="mb-6" v-if="!selectedUser">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="searchQuery"
          :placeholder="t('permissions.searchUser')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
      </div>
    </UCard>

    <!-- User Selection View -->
    <UCard v-if="!selectedUser">
      <UTable :data="filteredUsers" :columns="userColumns" :loading="loading" striped>
        <template #name-cell="{ row }">
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 items-center justify-center rounded-full bg-primary-100 text-primary-600 dark:bg-primary-900 dark:text-primary-300"
            >
              {{ row.original.firstName?.charAt(0) }}{{ row.original.lastName?.charAt(0) }}
            </div>
            <div>
              <div class="font-medium">
                {{ row.original.firstName }} {{ row.original.lastName }}
              </div>
              <div class="text-sm text-gray-500">{{ row.original.email }}</div>
            </div>
          </div>
        </template>
        <template #actions-cell="{ row }">
          <UButton size="sm" variant="outline" @click="selectUser(row.original)">
            {{ t("permissions.manage") }}
          </UButton>
        </template>
      </UTable>
    </UCard>

    <!-- Permission Management View -->
    <UCard v-if="selectedUser">
      <template #header>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <UButton variant="ghost" @click="selectedUser = undefined">
              <UIcon name="i-lucide-arrow-left" class="mr-2" />
              {{ t("actions.back") }}
            </UButton>
            <div>
              <h2 class="text-xl font-semibold">
                {{ selectedUser.firstName }} {{ selectedUser.lastName }}
              </h2>
              <p class="text-sm text-gray-500">{{ selectedUser.email }}</p>
            </div>
          </div>
          <div class="flex gap-2">
            <USelect
              v-model="selectedTemplate"
              :items="templates"
              :placeholder="t('permissions.selectTemplate')"
              class="w-48"
              @change="applyTemplate"
            />
          </div>
        </div>
      </template>

      <div class="space-y-4">
        <div v-for="category in permissionCategories" :key="category?.id ?? category?.label">
          <h3 class="mb-2 text-lg font-semibold">{{ getCategoryLabel(category) }}</h3>
          <div class="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
            <div
              v-for="perm in getPermissionsForCategory(category)"
              :key="perm.permissionId"
              class="flex items-center justify-between rounded-lg border p-3"
              :class="
                perm.granted
                  ? 'border-green-200 bg-green-50 dark:border-green-800 dark:bg-green-900/20'
                  : 'border-gray-200 bg-gray-50 dark:border-gray-700 dark:bg-gray-800'
              "
            >
              <div>
                <div class="font-medium">{{ perm.label }}</div>
                <div class="text-sm text-gray-500">{{ perm.page }} - {{ perm.action }}</div>
              </div>
              <div class="flex items-center gap-2">
                <span
                  v-if="perm.granted"
                  class="rounded-full bg-green-100 px-2 py-1 text-xs font-medium text-green-700 dark:bg-green-900 dark:text-green-300"
                >
                  {{ t("permissions.granted") }}
                </span>
                <span
                  v-else
                  class="rounded-full bg-gray-100 px-2 py-1 text-xs font-medium text-gray-700 dark:bg-gray-700 dark:text-gray-300"
                >
                  {{ t("permissions.denied") }}
                </span>
                <UButton size="sm" variant="ghost" @click="togglePermission(perm)">
                  <UIcon :name="perm.granted ? 'i-lucide-toggle-left' : 'i-lucide-toggle-right'" />
                </UButton>
              </div>
            </div>
          </div>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import type { PaginatedResponseModelsSanitizedUser, UserPermissionsResponse } from "~/utils/client";

type User = NonNullable<PaginatedResponseModelsSanitizedUser["data"]>[0];
type Permission = NonNullable<UserPermissionsResponse["permissions"]>[0];
type PermissionCategory = Permission["category"];

definePageMeta({
  requiresPermission: PERMISSIONS.users.users.update,
});

const { t } = useI18n();

const searchQuery = ref("");
const selectedUser = ref<User>();
const selectedTemplate = ref<number>();
const loading = ref(false);

const userColumns: TableColumn<User>[] = [
  { accessorKey: "name", header: t("users.name") },
  { accessorKey: "actions", header: "" },
];

const { data: users } = await useAsyncData("users", () => getApiUsers({}), {
  transform: (d) => d.data,
  default: () => [],
});

const { data: allPermissions } = await useAsyncData("permissions", () => getApiPermissions({}), {
  transform: (d) => d.data,
  default: () => [],
});

const { data: templates } = await useAsyncData("templates", () => getApiPermissionsTemplates({}), {
  transform: (d) => d.data,
  default: () => [],
});

const userPermissions = ref<UserPermissionsResponse["permissions"]>([]);

const filteredUsers = computed(() => {
  if (!users.value) return [];
  const query = searchQuery.value.toLowerCase();
  if (!query) return users.value;
  return users.value.filter(
    (u) =>
      u.firstName?.toLowerCase().includes(query) ||
      u.lastName?.toLowerCase().includes(query) ||
      u.email?.toLowerCase().includes(query)
  );
});

const permissionCategories = computed(() => {
  const groups = new Map<string, PermissionCategory>();
  (allPermissions.value ?? []).forEach((p: any) => {
    const category = p.category;
    if (!category) return;
    const key = `${category.id ?? 0}-${category.label ?? "General"}`;
    groups.set(key, category);
  });
  return Array.from(groups.values());
});

function getCategoryLabel(category: PermissionCategory) {
  return category?.label || "General";
}

function getPermissionsForCategory(category: PermissionCategory) {
  if (!category) return [];
  return (userPermissions.value ?? []).filter((p) => p.category?.id === category.id);
}

async function selectUser(user: User) {
  if (!user.id) return;
  selectedUser.value = user;
  selectedTemplate.value = undefined;
  loading.value = true;
  try {
    const response = await getApiPermissionsUserUserId({
      path: {
        userId: user.id.toString(),
      },
    });
    userPermissions.value = response.permissions ?? [];
  } catch (e) {
    userPermissions.value = [];
  } finally {
    loading.value = false;
  }
}

async function applyTemplate() {
  if (!selectedUser.value?.id || !selectedTemplate.value) return;
  loading.value = true;
  try {
    await postApiPermissionsUserUserIdTemplateTemplateId({
      path: {
        userId: selectedUser.value.id.toString(),
        templateId: selectedTemplate.value.toString(),
      },
    });
    await selectUser(selectedUser.value);
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

async function togglePermission(perm: Permission) {
  if (!selectedUser.value?.id || !perm.permissionId) return;
  loading.value = true;
  try {
    await postApiPermissionsUserUserIdPermissionId({
      path: {
        userId: selectedUser.value.id.toString(),
        permissionId: perm.permissionId.toString(),
      },
      body: {
        granted: !perm.granted,
      },
    });
    await selectUser(selectedUser.value);
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

let debounceTimer: ReturnType<typeof setTimeout>;
function debouncedSearch() {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {}, 300);
}
</script>
